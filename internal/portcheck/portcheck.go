// Package portcheck provides a high-level check that compares a fresh scan
// against the saved baseline and returns a structured result.
package portcheck

import (
	"fmt"
	"time"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/portdiff"
	"github.com/example/portwatch/internal/scanner"
)

// Result holds the outcome of a single port check.
type Result struct {
	Host      string
	ScannedAt time.Time
	Diff      portdiff.Diff
	Baseline  []int
	Current   []int
}

// Changed reports whether the check found any difference.
func (r Result) Changed() bool {
	return len(r.Diff.Opened) > 0 || len(r.Diff.Closed) > 0
}

// Checker runs a port scan and compares the result to a stored baseline.
type Checker struct {
	baselinePath string
	timeout      time.Duration
}

// New returns a Checker that persists baselines at baselinePath.
func New(baselinePath string, timeout time.Duration) *Checker {
	if timeout <= 0 {
		timeout = 500 * time.Millisecond
	}
	return &Checker{baselinePath: baselinePath, timeout: timeout}
}

// Run scans host:ports, loads (or creates) the baseline, and returns a Result.
func (c *Checker) Run(host string, ports []int) (Result, error) {
	open, err := scanner.OpenPorts(host, ports, c.timeout)
	if err != nil {
		return Result{}, fmt.Errorf("scan: %w", err)
	}

	b, err := baseline.Load(c.baselinePath)
	if err != nil {
		// No baseline yet — save current scan as the new baseline.
		b = baseline.New(host, open)
		if saveErr := b.Save(c.baselinePath); saveErr != nil {
			return Result{}, fmt.Errorf("save baseline: %w", saveErr)
		}
		return Result{
			Host:      host,
			ScannedAt: time.Now(),
			Diff:      portdiff.Diff{},
			Baseline:  open,
			Current:   open,
		}, nil
	}

	diff := portdiff.Compute(b.Ports, open)
	return Result{
		Host:      host,
		ScannedAt: time.Now(),
		Diff:      diff,
		Baseline:  b.Ports,
		Current:   open,
	}, nil
}
