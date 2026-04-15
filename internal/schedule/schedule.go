// Package schedule provides utilities for running periodic port scans
// based on a configured interval.
package schedule

import (
	"context"
	"fmt"
	"time"
)

// Job holds the configuration for a scheduled scan task.
type Job struct {
	Interval time.Duration
	Task     func(ctx context.Context) error
	OnError  func(err error)
}

// Run starts the job and executes Task on every tick of the interval.
// It blocks until the provided context is cancelled.
func (j *Job) Run(ctx context.Context) {
	if j.Interval <= 0 {
		panic("schedule: interval must be positive")
	}
	if j.Task == nil {
		panic("schedule: task must not be nil")
	}

	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()

	// Run immediately on first invocation.
	j.runTask(ctx)

	for {
		select {
		case <-ticker.C:
			j.runTask(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (j *Job) runTask(ctx context.Context) {
	if err := j.Task(ctx); err != nil {
		if j.OnError != nil {
			j.OnError(err)
		}
	}
}

// ParseDuration wraps time.ParseDuration and returns a user-friendly error.
func ParseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("schedule: invalid duration %q: %w", s, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("schedule: duration must be positive, got %q", s)
	}
	return d, nil
}
