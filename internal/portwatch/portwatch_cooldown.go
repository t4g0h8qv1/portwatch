package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// ScanCooldownConfig holds configuration for the scan cooldown manager.
type ScanCooldownConfig struct {
	// MinGap is the minimum duration required between successive scans of the same target.
	MinGap time.Duration
}

// DefaultScanCooldownConfig returns a ScanCooldownConfig with sensible defaults.
func DefaultScanCooldownConfig() ScanCooldownConfig {
	return ScanCooldownConfig{
		MinGap: 30 * time.Second,
	}
}

// ScanCooldownManager enforces a minimum gap between scans of the same target.
type ScanCooldownManager struct {
	cfg    ScanCooldownConfig
	mu     sync.Mutex
	lastAt map[string]time.Time
	now    func() time.Time
}

// NewScanCooldownManager creates a ScanCooldownManager with the given config.
// Returns an error if MinGap is not positive.
func NewScanCooldownManager(cfg ScanCooldownConfig) (*ScanCooldownManager, error) {
	if cfg.MinGap <= 0 {
		return nil, fmt.Errorf("portwatch: cooldown MinGap must be positive, got %v", cfg.MinGap)
	}
	return &ScanCooldownManager{
		cfg:    cfg,
		lastAt: make(map[string]time.Time),
		now:    time.Now,
	}, nil
}

// Ready reports whether target may be scanned now (i.e. enough time has elapsed
// since the last scan). It does not record the attempt; call Observe after scanning.
func (m *ScanCooldownManager) Ready(target string) bool {
	if target == "" {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	last, ok := m.lastAt[target]
	if !ok {
		return true
	}
	return m.now().Sub(last) >= m.cfg.MinGap
}

// Observe records that target was scanned at the current time.
func (m *ScanCooldownManager) Observe(target string) {
	if target == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastAt[target] = m.now()
}

// Reset clears the cooldown state for target, allowing an immediate scan.
func (m *ScanCooldownManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.lastAt, target)
}

// NextReady returns the time at which target will next be eligible for scanning.
// If target has no recorded scan, the zero time is returned.
func (m *ScanCooldownManager) NextReady(target string) time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	last, ok := m.lastAt[target]
	if !ok {
		return time.Time{}
	}
	return last.Add(m.cfg.MinGap)
}
