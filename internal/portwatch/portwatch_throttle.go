package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DefaultThrottleConfig returns a ThrottleConfig with sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		MinGap: 5 * time.Second,
	}
}

// ThrottleConfig holds configuration for the scan throttle manager.
type ThrottleConfig struct {
	// MinGap is the minimum duration required between successive scans
	// for the same target.
	MinGap time.Duration
}

// ScanThrottleManager enforces a minimum gap between scans per target.
type ScanThrottleManager struct {
	mu      sync.Mutex
	cfg     ThrottleConfig
	lastRun map[string]time.Time
	now     func() time.Time
}

// NewScanThrottleManager creates a new ScanThrottleManager with the given config.
func NewScanThrottleManager(cfg ThrottleConfig) (*ScanThrottleManager, error) {
	if cfg.MinGap <= 0 {
		return nil, errors.New("portwatch: throttle MinGap must be positive")
	}
	return &ScanThrottleManager{
		cfg:     cfg,
		lastRun: make(map[string]time.Time),
		now:     time.Now,
	}, nil
}

// Allow reports whether a scan for target is permitted given the configured
// minimum gap. It does not record the attempt; call Observe after the scan.
func (m *ScanThrottleManager) Allow(target string) bool {
	if target == "" {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	last, ok := m.lastRun[target]
	if !ok {
		return true
	}
	return m.now().Sub(last) >= m.cfg.MinGap
}

// Observe records that a scan for target occurred at the current time.
func (m *ScanThrottleManager) Observe(target string) {
	if target == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastRun[target] = m.now()
}

// Reset clears the recorded scan time for target, allowing the next scan
// to proceed immediately regardless of the configured gap.
func (m *ScanThrottleManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.lastRun, target)
}

// LastScan returns the time of the most recent observed scan for target
// and whether a record exists.
func (m *ScanThrottleManager) LastScan(target string) (time.Time, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.lastRun[target]
	return t, ok
}
