package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// WindowConfig holds configuration for the scan window manager.
type WindowConfig struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// ScanWindowManager restricts scans to a configured daily time window.
type ScanWindowManager struct {
	mu      sync.RWMutex
	windows map[string]WindowConfig
	now     func() time.Time
}

// NewScanWindowManager returns a new ScanWindowManager.
func NewScanWindowManager() *ScanWindowManager {
	return &ScanWindowManager{
		windows: make(map[string]WindowConfig),
		now:     time.Now,
	}
}

// Set registers a time window for a target.
func (m *ScanWindowManager) Set(target string, cfg WindowConfig) error {
	if target == "" {
		return fmt.Errorf("portwatch: window target must not be empty")
	}
	if cfg.End <= cfg.Start {
		return fmt.Errorf("portwatch: window end must be after start")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.windows[target] = cfg
	return nil
}

// Remove deletes the window for a target.
func (m *ScanWindowManager) Remove(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.windows, target)
}

// Allowed reports whether a scan is permitted for target at the current time.
// If no window is registered, scans are always allowed.
func (m *ScanWindowManager) Allowed(target string) bool {
	m.mu.RLock()
	cfg, ok := m.windows[target]
	m.mu.RUnlock()
	if !ok {
		return true
	}
	now := m.now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	offset := now.Sub(midnight)
	return offset >= cfg.Start && offset < cfg.End
}

// Targets returns all registered target names.
func (m *ScanWindowManager) Targets() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, 0, len(m.windows))
	for t := range m.windows {
		out = append(out, t)
	}
	return out
}
