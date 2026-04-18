// Package portwatch wires together scanning, baselining, and alerting.
package portwatch

import (
	"context"
	"fmt"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/notify"
	"github.com/example/portwatch/internal/portdiff"
	"github.com/example/portwatch/internal/scanner"
)

// Config holds the runtime options for a single watch run.
type Config struct {
	Target      string
	Ports       []int
	BaselinePath string
	Notifier    notify.Notifier
}

// Run performs one scan cycle: scan ports, compare against baseline,
// notify on changes, and persist an updated baseline.
func Run(ctx context.Context, cfg Config) error {
	open, err := scanner.OpenPorts(ctx, cfg.Target, cfg.Ports)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	b, err := baseline.Load(cfg.BaselinePath)
	if err != nil {
		// First run — save current state as baseline.
		b = baseline.New(cfg.Target, open)
		return b.Save(cfg.BaselinePath)
	}

	diff := portdiff.Compute(b.Ports, open)
	if diff.HasChanges() {
		msg := portdiff.Summary(diff)
		if notifyErr := cfg.Notifier.Notify(ctx, msg); notifyErr != nil {
			return fmt.Errorf("notify: %w", notifyErr)
		}
	}

	b.Ports = open
	return b.Save(cfg.BaselinePath)
}
