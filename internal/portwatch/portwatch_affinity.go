package portwatch

import (
	"errors"
	"sync"
)

// AffinityConfig holds configuration for the scan affinity manager.
type AffinityConfig struct {
	// MaxTargets is the maximum number of targets that can be pinned to a worker.
	MaxTargets int
}

// DefaultAffinityConfig returns a sensible default AffinityConfig.
func DefaultAffinityConfig() AffinityConfig {
	return AffinityConfig{
		MaxTargets: 64,
	}
}

// affinityEntry records which worker a target is pinned to.
type affinityEntry struct {
	worker int
}

// ScanAffinityManager pins targets to consistent workers to reduce
// cache thrashing and improve scan locality.
type ScanAffinityManager struct {
	mu      sync.RWMutex
	cfg     AffinityConfig
	entries map[string]affinityEntry
	next    int
}

// NewScanAffinityManager constructs a ScanAffinityManager with the given config.
func NewScanAffinityManager(cfg AffinityConfig) (*ScanAffinityManager, error) {
	if cfg.MaxTargets <= 0 {
		return nil, errors.New("portwatch: affinity MaxTargets must be greater than zero")
	}
	return &ScanAffinityManager{
		cfg:     cfg,
		entries: make(map[string]affinityEntry),
	}, nil
}

// Assign returns the worker ID assigned to the given target, creating a new
// assignment if one does not yet exist. Worker IDs are in [0, MaxTargets).
func (m *ScanAffinityManager) Assign(target string) (int, error) {
	if target == "" {
		return 0, errors.New("portwatch: affinity target must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.entries[target]; ok {
		return e.worker, nil
	}
	worker := m.next % m.cfg.MaxTargets
	m.next++
	m.entries[target] = affinityEntry{worker: worker}
	return worker, nil
}

// Get returns the worker assignment for target, and false if none exists.
func (m *ScanAffinityManager) Get(target string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[target]
	return e.worker, ok
}

// Remove deletes the affinity assignment for the given target.
func (m *ScanAffinityManager) Remove(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}

// Len returns the number of targets currently assigned.
func (m *ScanAffinityManager) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}
