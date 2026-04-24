package portwatch

import (
	"sync"
	"time"
)

// DebounceConfig holds configuration for the debounce manager.
type DebounceConfig struct {
	// Window is the minimum duration between alerts for the same target.
	Window time.Duration
}

// DefaultDebounceConfig returns a DebounceConfig with sensible defaults.
func DefaultDebounceConfig() DebounceConfig {
	return DebounceConfig{
		Window: 30 * time.Second,
	}
}

// debounceEntry tracks the last alert time for a target.
type debounceEntry struct {
	lastAlert time.Time
}

// DebounceManager suppresses repeated alerts for the same target within a
// configurable time window, preventing alert storms during transient changes.
type DebounceManager struct {
	mu      sync.Mutex
	cfg     DebounceConfig
	entries map[string]*debounceEntry
	now     func() time.Time
}

// NewDebounceManager creates a DebounceManager with the given config.
// Returns an error if the window is non-positive.
func NewDebounceManager(cfg DebounceConfig) (*DebounceManager, error) {
	if cfg.Window <= 0 {
		return nil, errInvalidDebounceWindow
	}
	return &DebounceManager{
		cfg:     cfg,
		entries: make(map[string]*debounceEntry),
		now:     time.Now,
	}, nil
}

// Ready reports whether target is ready to fire an alert.
// A target is ready if it has never alerted or its last alert
// was more than Window ago.
func (d *DebounceManager) Ready(target string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	e, ok := d.entries[target]
	if !ok {
		return true
	}
	return d.now().Sub(e.lastAlert) >= d.cfg.Window
}

// Observe records that an alert was fired for target at the current time.
func (d *DebounceManager) Observe(target string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if e, ok := d.entries[target]; ok {
		e.lastAlert = d.now()
		return
	}
	d.entries[target] = &debounceEntry{lastAlert: d.now()}
}

// Reset clears the debounce state for target.
func (d *DebounceManager) Reset(target string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, target)
}
