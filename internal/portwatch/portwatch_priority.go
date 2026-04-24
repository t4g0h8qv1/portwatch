package portwatch

import (
	"fmt"
	"sync"
)

// PriorityLevel represents the scan priority for a target.
type PriorityLevel int

const (
	PriorityLow    PriorityLevel = 1
	PriorityNormal PriorityLevel = 5
	PriorityHigh   PriorityLevel = 10
)

// DefaultPriorityConfig returns a PriorityConfig with sensible defaults.
func DefaultPriorityConfig() PriorityConfig {
	return PriorityConfig{
		DefaultLevel: PriorityNormal,
	}
}

// PriorityConfig holds configuration for the scan priority manager.
type PriorityConfig struct {
	DefaultLevel PriorityLevel
}

// ScanPriorityManager assigns and retrieves scan priorities per target.
type ScanPriorityManager struct {
	mu       sync.RWMutex
	levels   map[string]PriorityLevel
	defLevel PriorityLevel
}

// NewScanPriorityManager creates a new ScanPriorityManager.
// Returns an error if the default level is not positive.
func NewScanPriorityManager(cfg PriorityConfig) (*ScanPriorityManager, error) {
	if cfg.DefaultLevel <= 0 {
		return nil, fmt.Errorf("portwatch: default priority level must be positive, got %d", cfg.DefaultLevel)
	}
	return &ScanPriorityManager{
		levels:   make(map[string]PriorityLevel),
		defLevel: cfg.DefaultLevel,
	}, nil
}

// Set assigns a priority level to a target.
func (m *ScanPriorityManager) Set(target string, level PriorityLevel) error {
	if target == "" {
		return fmt.Errorf("portwatch: target must not be empty")
	}
	if level <= 0 {
		return fmt.Errorf("portwatch: priority level must be positive, got %d", level)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.levels[target] = level
	return nil
}

// Get returns the priority level for a target, or the default if not set.
func (m *ScanPriorityManager) Get(target string) PriorityLevel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if lvl, ok := m.levels[target]; ok {
		return lvl
	}
	return m.defLevel
}

// Reset removes any explicit priority for the target, reverting to the default.
func (m *ScanPriorityManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.levels, target)
}

// Targets returns all targets with an explicit priority set.
func (m *ScanPriorityManager) Targets() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, 0, len(m.levels))
	for t := range m.levels {
		out = append(out, t)
	}
	return out
}
