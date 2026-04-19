package portwatch

import (
	"context"
	"fmt"
	"io"
	"time"
)

// RunnerConfig holds configuration for a continuous scan runner.
type RunnerConfig struct {
	ScanConfig Config
	Interval   time.Duration
	MaxScans   int // 0 means unlimited
	Out        io.Writer
}

// RunnerResult summarises a completed runner session.
type RunnerResult struct {
	ScansCompleted int
	Errors         int
	LastScan       time.Time
}

// Runner executes periodic scans until the context is cancelled or MaxScans is reached.
type Runner struct {
	cfg RunnerConfig
}

// NewRunner creates a Runner with the given config.
func NewRunner(cfg RunnerConfig) (*Runner, error) {
	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("portwatch runner: interval must be positive")
	}
	if cfg.Out == nil {
		return nil, fmt.Errorf("portwatch runner: output writer must not be nil")
	}
	return &Runner{cfg: cfg}, nil
}

// Start runs scans on the configured interval until ctx is done or MaxScans reached.
func (r *Runner) Start(ctx context.Context) RunnerResult {
	result := RunnerResult{}
	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()

	runScan := func() {
		if err := Run(ctx, r.cfg.ScanConfig); err != nil {
			fmt.Fprintf(r.cfg.Out, "[runner] scan error: %v\n", err)
			result.Errors++
		}
		result.ScansCompleted++
		result.LastScan = time.Now()
	}

	// Run immediately on start.
	runScan()
	if r.cfg.MaxScans > 0 && result.ScansCompleted >= r.cfg.MaxScans {
		return result
	}

	for {
		select {
		case <-ctx.Done():
			return result
		case <-ticker.C:
			runScan()
			if r.cfg.MaxScans > 0 && result.ScansCompleted >= r.cfg.MaxScans {
				return result
			}
		}
	}
}
