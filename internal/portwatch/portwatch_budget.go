package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// DefaultScanBudgetConfig returns a ScanBudgetConfig with sensible defaults.
func DefaultScanBudgetConfig() ScanBudgetConfig {
	return ScanBudgetConfig{
		MaxScansPerDay: 1440,
		Window:         24 * time.Hour,
	}
}

// ScanBudgetConfig holds configuration for the scan budget manager.
type ScanBudgetConfig struct {
	MaxScansPerDay int
	Window         time.Duration
}

// scanBudgetEntry tracks usage for a single target.
type scanBudgetEntry struct {
	count     int
	windowEnd time.Time
}

// ScanBudgetManager enforces a maximum number of scans per target within a
// rolling time window.
type ScanBudgetManager struct {
	mu      sync.Mutex
	cfg     ScanBudgetConfig
	entries map[string]*scanBudgetEntry
	now     func() time.Time
}

// NewScanBudgetManager creates a new ScanBudgetManager.
func NewScanBudgetManager(cfg ScanBudgetConfig) (*ScanBudgetManager, error) {
	if cfg.MaxScansPerDay <= 0 {
		return nil, fmt.Errorf("portwatch: MaxScansPerDay must be > 0, got %d", cfg.MaxScansPerDay)
	}
	if cfg.Window <= 0 {
		return nil, fmt.Errorf("portwatch: Window must be > 0, got %s", cfg.Window)
	}
	return &ScanBudgetManager{
		cfg:     cfg,
		entries: make(map[string]*scanBudgetEntry),
		now:     time.Now,
	}, nil
}

// Allow returns true and records the scan if the target is within budget.
// Returns false if the budget for the current window is exhausted.
func (m *ScanBudgetManager) Allow(target string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	e, ok := m.entries[target]
	if !ok || now.After(e.windowEnd) {
		m.entries[target] = &scanBudgetEntry{count: 1, windowEnd: now.Add(m.cfg.Window)}
		return true
	}
	if e.count >= m.cfg.MaxScansPerDay {
		return false
	}
	e.count++
	return true
}

// Remaining returns the number of scans remaining in the current window for
// the given target.
func (m *ScanBudgetManager) Remaining(target string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	e, ok := m.entries[target]
	if !ok || now.After(e.windowEnd) {
		return m.cfg.MaxScansPerDay
	}
	rem := m.cfg.MaxScansPerDay - e.count
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the budget entry for the given target.
func (m *ScanBudgetManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}
