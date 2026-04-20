package portwatch

import (
	"errors"
	"sync"
	"time"
)

// RetryConfig holds configuration for the scan retry manager.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   2 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// ScanRetryManager tracks per-target retry state for failed scans.
type ScanRetryManager struct {
	mu      sync.Mutex
	cfg     RetryConfig
	attempt map[string]int
}

// NewScanRetryManager creates a ScanRetryManager with the given config.
// Returns an error if MaxAttempts < 1 or BaseDelay <= 0.
func NewScanRetryManager(cfg RetryConfig) (*ScanRetryManager, error) {
	if cfg.MaxAttempts < 1 {
		return nil, errors.New("portwatch: MaxAttempts must be >= 1")
	}
	if cfg.BaseDelay <= 0 {
		return nil, errors.New("portwatch: BaseDelay must be positive")
	}
	if cfg.MaxDelay < cfg.BaseDelay {
		cfg.MaxDelay = cfg.BaseDelay
	}
	return &ScanRetryManager{
		cfg:     cfg,
		attempt: make(map[string]int),
	}, nil
}

// ShouldRetry returns true if the target has not exhausted its retry budget.
func (m *ScanRetryManager) ShouldRetry(target string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.attempt[target] < m.cfg.MaxAttempts
}

// NextDelay records a failure for target and returns the backoff delay to
// wait before the next attempt. Returns 0 once retries are exhausted.
func (m *ScanRetryManager) NextDelay(target string) time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	attempts := m.attempt[target]
	if attempts >= m.cfg.MaxAttempts {
		return 0
	}
	m.attempt[target]++
	delay := m.cfg.BaseDelay * (1 << uint(attempts))
	if delay > m.cfg.MaxDelay {
		delay = m.cfg.MaxDelay
	}
	return delay
}

// Reset clears the retry counter for target, e.g. after a successful scan.
func (m *ScanRetryManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.attempt, target)
}

// Attempts returns the current failure count for target.
func (m *ScanRetryManager) Attempts(target string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.attempt[target]
}
