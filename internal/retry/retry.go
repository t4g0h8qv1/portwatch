// Package retry provides a simple exponential-backoff retry helper
// used when scanning or notifying over unreliable connections.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Config holds retry behaviour parameters.
type Config struct {
	// MaxAttempts is the total number of tries (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// MaxDelay caps the exponential growth of the delay.
	MaxDelay time.Duration
	// Multiplier is the factor applied to the delay after each failure.
	Multiplier float64
}

// DefaultConfig returns sensible defaults for most use-cases.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  3,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Do calls fn up to cfg.MaxAttempts times, backing off between failures.
// It stops early if ctx is cancelled or fn returns nil.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts < 1 {
		return errors.New("retry: MaxAttempts must be >= 1")
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = 2.0
	}

	delay := cfg.InitialDelay
	var last error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: context cancelled after %d attempt(s): %w", attempt-1, err)
		}

		last = fn()
		if last == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry: context cancelled: %w", ctx.Err())
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	return fmt.Errorf("retry: all %d attempt(s) failed: %w", cfg.MaxAttempts, last)
}
