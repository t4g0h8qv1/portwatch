package portwatch

import (
	"errors"
	"sync"
	"time"
)

// FenceConfig holds configuration for the scan fence manager.
type FenceConfig struct {
	// MaxAge is the maximum duration a fence remains active.
	MaxAge time.Duration
}

// DefaultFenceConfig returns a FenceConfig with sensible defaults.
func DefaultFenceConfig() FenceConfig {
	return FenceConfig{
		MaxAge: 10 * time.Minute,
	}
}

type fenceEntry struct {
	reason  string
	expiresAt time.Time
}

// ScanFenceManager prevents scans from running against fenced targets.
type ScanFenceManager struct {
	mu     sync.RWMutex
	cfg    FenceConfig
	fences map[string]fenceEntry
	now    func() time.Time
}

// NewScanFenceManager creates a ScanFenceManager with the given config.
func NewScanFenceManager(cfg FenceConfig) (*ScanFenceManager, error) {
	if cfg.MaxAge <= 0 {
		return nil, errors.New("fence: MaxAge must be positive")
	}
	return &ScanFenceManager{
		cfg:    cfg,
		fences: make(map[string]fenceEntry),
		now:    time.Now,
	}, nil
}

// Fence marks target as fenced for the configured MaxAge with an optional reason.
func (m *ScanFenceManager) Fence(target, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fences[target] = fenceEntry{
		reason:    reason,
		expiresAt: m.now().Add(m.cfg.MaxAge),
	}
}

// Unfence removes the fence for target immediately.
func (m *ScanFenceManager) Unfence(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.fences, target)
}

// IsFenced returns true when target has an active (non-expired) fence.
func (m *ScanFenceManager) IsFenced(target string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.fences[target]
	if !ok {
		return false
	}
	if m.now().After(e.expiresAt) {
		return false
	}
	return true
}

// Reason returns the reason string for a fenced target, or empty string.
func (m *ScanFenceManager) Reason(target string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.fences[target]
	if !ok || m.now().After(e.expiresAt) {
		return ""
	}
	return e.reason
}

// Prune removes all expired fence entries.
func (m *ScanFenceManager) Prune() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	for k, e := range m.fences {
		if now.After(e.expiresAt) {
			delete(m.fences, k)
		}
	}
}
