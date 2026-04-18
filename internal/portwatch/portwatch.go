// Package portwatch wires together scanning, baseline diffing, and notification.
package portwatch

import (
	"context"
	"fmt"
	"time"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/notify"
	"github.com/example/portwatch/internal/portdiff"
	"github.com/example/portwatch/internal/scanner"
)

// Config holds the runtime configuration for a single watch cycle.
type Config struct {
	Target      string
	Ports       []int
	BaselinePath string
	Timeout     time.Duration
	Notifier    notify.Notifier
}

// Run performs one scan cycle: scan ports, compare against baseline,
// notify on changes, and persist the updated baseline.
func Run(ctx context.Context, cfg Config) error {
	open, err := scanner.OpenPorts(cfg.Target, cfg.Ports, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	b, err := baseline.Load(cfg.BaselinePath)
	if err != nil {
		// No baseline yet — save current scan as baseline.
		b = baseline.New(cfg.Target, open)
		if saveErr := b.Save(cfg.BaselinePath); saveErr != nil {
			return fmt.Errorf("save baseline: %w", saveErr)
		}
		return nil
	}

	diff := portdiff.Compute(b.Ports, open)
	if diff.Empty() {
		return nil
	}

	msg := portdiff.Summary(diff)
	if notifyErr := cfg.Notifier.Notify(ctx, msg); notifyErr != nil {
		return fmt.Errorf("notify: %w", notifyErr)
	}

	b.Ports = open
	b.UpdatedAt = time.Now()
	if saveErr := b.Save(cfg.BaselinePath); saveErr != nil {
		return fmt.Errorf("update baseline: %w", saveErr)
	}
	return nil
}
