package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DefaultWarmupConfig returns a WarmupConfig with sensible defaults.
func DefaultWarmupConfig() WarmupConfig {
	return WarmupConfig{
		MinScans: 3,
		MaxWait:  2 * time.Minute,
	}
}

// WarmupConfig controls how long a target is considered "warming up".
type WarmupConfig struct {
	// MinScans is the minimum number of successful scans before a target is warm.
	MinScans int
	// MaxWait is the maximum time to wait before declaring a target warm regardless of scan count.
	MaxWait time.Duration
}

func (c WarmupConfig) validate() error {
	if c.MinScans < 1 {
		return errors.New("portwatch/warmup: MinScans must be at least 1")
	}
	if c.MaxWait <= 0 {
		return errors.New("portwatch/warmup: MaxWait must be positive")
	}
	return nil
}

type warmupEntry struct {
	scans     int
	firstSeen time.Time
}

// WarmupManager tracks whether a target has completed its warmup period.
type WarmupManager struct {
	mu      sync.Mutex
	cfg     WarmupConfig
	entries map[string]*warmupEntry
	now     func() time.Time
}

// NewWarmupManager creates a WarmupManager with the given config.
func NewWarmupManager(cfg WarmupConfig) (*WarmupManager, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &WarmupManager{
		cfg:     cfg,
		entries: make(map[string]*warmupEntry),
		now:     time.Now,
	}, nil
}

// RecordScan records a successful scan for target, incrementing its scan count.
func (w *WarmupManager) RecordScan(target string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	e, ok := w.entries[target]
	if !ok {
		e = &warmupEntry{firstSeen: w.now()}
		w.entries[target] = e
	}
	e.scans++
}

// IsWarm reports whether target has completed warmup.
func (w *WarmupManager) IsWarm(target string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	e, ok := w.entries[target]
	if !ok {
		return false
	}
	if e.scans >= w.cfg.MinScans {
		return true
	}
	return w.now().Sub(e.firstSeen) >= w.cfg.MaxWait
}

// Reset removes warmup state for target.
func (w *WarmupManager) Reset(target string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, target)
}

// Scans returns the current scan count for target.
func (w *WarmupManager) Scans(target string) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if e, ok := w.entries[target]; ok {
		return e.scans
	}
	return 0
}
