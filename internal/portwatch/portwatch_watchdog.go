package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DefaultWatchdogConfig returns a WatchdogConfig with sensible defaults.
func DefaultWatchdogConfig() WatchdogConfig {
	return WatchdogConfig{
		MaxSilence: 5 * time.Minute,
		CheckInterval: 30 * time.Second,
	}
}

// WatchdogConfig holds configuration for the watchdog manager.
type WatchdogConfig struct {
	MaxSilence    time.Duration
	CheckInterval time.Duration
}

// watchdogEntry records the last heartbeat time for a target.
type watchdogEntry struct {
	lastSeen time.Time
	triggered bool
}

// WatchdogManager tracks whether targets have been heard from recently.
type WatchdogManager struct {
	mu      sync.Mutex
	cfg     WatchdogConfig
	entries map[string]*watchdogEntry
	now     func() time.Time
}

// NewWatchdogManager creates a WatchdogManager with the given config.
func NewWatchdogManager(cfg WatchdogConfig) (*WatchdogManager, error) {
	if cfg.MaxSilence <= 0 {
		return nil, errors.New("watchdog: MaxSilence must be positive")
	}
	if cfg.CheckInterval <= 0 {
		return nil, errors.New("watchdog: CheckInterval must be positive")
	}
	return &WatchdogManager{
		cfg:     cfg,
		entries: make(map[string]*watchdogEntry),
		now:     time.Now,
	}, nil
}

// Ping records that target was seen at the current time.
func (w *WatchdogManager) Ping(target string) error {
	if target == "" {
		return errors.New("watchdog: target must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[target] = &watchdogEntry{lastSeen: w.now()}
	return nil
}

// IsExpired reports whether target has exceeded MaxSilence since its last ping.
func (w *WatchdogManager) IsExpired(target string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	e, ok := w.entries[target]
	if !ok {
		return true
	}
	return w.now().Sub(e.lastSeen) > w.cfg.MaxSilence
}

// Expired returns all targets that have exceeded MaxSilence.
func (w *WatchdogManager) Expired() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	var out []string
	for target, e := range w.entries {
		if w.now().Sub(e.lastSeen) > w.cfg.MaxSilence {
			out = append(out, target)
		}
	}
	return out
}

// Reset removes a target from the watchdog registry.
func (w *WatchdogManager) Reset(target string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, target)
}
