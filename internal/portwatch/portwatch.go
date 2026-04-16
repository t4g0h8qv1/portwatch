// Package portwatch wires together scanning, baselining, alerting, and
// notification into a single reusable Run function.
package portwatch

import (
	"context"
	"fmt"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/notify"
	"github.com/example/portwatch/internal/scanner"
)

// Config holds the runtime parameters for a single watch cycle.
type Config struct {
	Target       string
	Ports        []int
	BaselinePath string
	HistoryPath  string
	Timeout      time.Duration
	Notifier     notify.Notifier
}

// Result is returned by Run after one scan cycle completes.
type Result struct {
	Open    []int
	Added   []int
	Removed []int
	ScannedAt time.Time
}

// Run performs one full scan cycle: scan → diff against baseline → notify.
func Run(ctx context.Context, cfg Config) (Result, error) {
	open, err := scanner.OpenPorts(ctx, cfg.Target, cfg.Ports, cfg.Timeout)
	if err != nil {
		return Result{}, fmt.Errorf("scan: %w", err)
	}

	bl, err := baseline.Load(cfg.BaselinePath)
	if err != nil {
		bl = baseline.New(cfg.Target, open)
		if saveErr := bl.Save(cfg.BaselinePath); saveErr != nil {
			return Result{}, fmt.Errorf("save baseline: %w", saveErr)
		}
	}

	evt := alert.Evaluate(bl, open)

	if len(evt.Added) > 0 || len(evt.Removed) > 0 {
		if cfg.Notifier != nil {
			if nErr := cfg.Notifier.Notify(ctx, evt); nErr != nil {
				return Result{}, fmt.Errorf("notify: %w", nErr)
			}
		}
		bl = baseline.New(cfg.Target, open)
		if saveErr := bl.Save(cfg.BaselinePath); saveErr != nil {
			return Result{}, fmt.Errorf("update baseline: %w", saveErr)
		}
	}

	if cfg.HistoryPath != "" {
		if hErr := history.Record(cfg.HistoryPath, cfg.Target, open); hErr != nil {
			return Result{}, fmt.Errorf("history: %w", hErr)
		}
	}

	return Result{
		Open:      open,
		Added:     evt.Added,
		Removed:   evt.Removed,
		ScannedAt: time.Now(),
	}, nil
}
