package portwatch

import (
	"math"
	"sync"
	"time"
)

// BackoffConfig holds parameters for exponential backoff.
type BackoffConfig struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// DefaultBackoffConfig returns a BackoffConfig with sensible defaults.
func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		InitialInterval: 5 * time.Second,
		MaxInterval:     5 * time.Minute,
		Multiplier:      2.0,
	}
}

// BackoffManager tracks per-target consecutive failure counts and computes
// the next backoff interval using exponential backoff.
type BackoffManager struct {
	mu      sync.Mutex
	cfg     BackoffConfig
	failures map[string]int
}

// NewBackoffManager creates a BackoffManager with the given config.
func NewBackoffManager(cfg BackoffConfig) *BackoffManager {
	return &BackoffManager{
		cfg:      cfg,
		failures: make(map[string]int),
	}
}

// RecordFailure increments the failure count for target and returns the next
// wait interval.
func (b *BackoffManager) RecordFailure(target string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[target]++
	return b.interval(b.failures[target])
}

// RecordSuccess resets the failure count for target.
func (b *BackoffManager) RecordSuccess(target string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.failures, target)
}

// Failures returns the current consecutive failure count for target.
func (b *BackoffManager) Failures(target string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures[target]
}

func (b *BackoffManager) interval(failures int) time.Duration {
	if failures <= 0 {
		return 0
	}
	mult := math.Pow(b.cfg.Multiplier, float64(failures-1))
	d := time.Duration(float64(b.cfg.InitialInterval) * mult)
	if d > b.cfg.MaxInterval {
		d = b.cfg.MaxInterval
	}
	return d
}
