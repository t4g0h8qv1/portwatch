package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// ScanQuotaConfig holds configuration for the scan quota manager.
type ScanQuotaConfig struct {
	MaxScansPerHour int
	Window          time.Duration
}

// DefaultScanQuotaConfig returns a ScanQuotaConfig with sensible defaults.
func DefaultScanQuotaConfig() ScanQuotaConfig {
	return ScanQuotaConfig{
		MaxScansPerHour: 60,
		Window:          time.Hour,
	}
}

type scanQuotaEntry struct {
	count     int
	windowEnd time.Time
}

// ScanQuotaManager tracks per-target scan counts within a rolling window.
type ScanQuotaManager struct {
	mu      sync.Mutex
	cfg     ScanQuotaConfig
	entries map[string]*scanQuotaEntry
	now     func() time.Time
}

// NewScanQuotaManager creates a ScanQuotaManager with the given config.
func NewScanQuotaManager(cfg ScanQuotaConfig) (*ScanQuotaManager, error) {
	if cfg.MaxScansPerHour <= 0 {
		return nil, fmt.Errorf("portwatch: MaxScansPerHour must be positive")
	}
	if cfg.Window <= 0 {
		return nil, fmt.Errorf("portwatch: Window must be positive")
	}
	return &ScanQuotaManager{
		cfg:     cfg,
		entries: make(map[string]*scanQuotaEntry),
		now:     time.Now,
	}, nil
}

// Allow returns true if the target has not exceeded its quota for the current window.
func (m *ScanQuotaManager) Allow(target string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	e, ok := m.entries[target]
	if !ok || now.After(e.windowEnd) {
		m.entries[target] = &scanQuotaEntry{count: 1, windowEnd: now.Add(m.cfg.Window)}
		return true
	}
	if e.count >= m.cfg.MaxScansPerHour {
		return false
	}
	e.count++
	return true
}

// Remaining returns the number of scans remaining for the target in the current window.
func (m *ScanQuotaManager) Remaining(target string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	e, ok := m.entries[target]
	if !ok || now.After(e.windowEnd) {
		return m.cfg.MaxScansPerHour
	}
	remaining := m.cfg.MaxScansPerHour - e.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears the quota state for a target.
func (m *ScanQuotaManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}
