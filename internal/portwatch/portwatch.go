// Package portwatch orchestrates the core scan-diff-notify loop.
package portwatch

import (
	"context"
	"fmt"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/scanner"
)

// Notifier is called when port changes are detected.
type Notifier interface {
	Notify(ctx context.Context, event Event) error
}

// Config holds the runtime configuration for a single watch target.
type Config struct {
	Host     string
	Ports    []int
	Baseline string // path to baseline file
}

// Run performs one scan cycle: scan → diff → notify.
func Run(ctx context.Context, cfg Config, n Notifier) error {
	open, err := scanner.OpenPorts(cfg.Host, cfg.Ports)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	bl, err := baseline.Load(cfg.Baseline)
	if err != nil {
		// First run — save baseline and return.
		bl = baseline.New(cfg.Host, open)
		return bl.Save(cfg.Baseline)
	}

	ev := buildEvent(cfg.Host, bl.Ports, open)
	if !HasChanges(ev) {
		return nil
	}

	if err := n.Notify(ctx, ev); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

	bl.Ports = open
	return bl.Save(cfg.Baseline)
}
