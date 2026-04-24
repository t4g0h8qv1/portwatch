package portwatch

import (
	"errors"
	"sync"
	"time"
)

// ConcurrencyConfig holds configuration for the scan concurrency manager.
type ConcurrencyConfig struct {
	// MaxConcurrent is the maximum number of simultaneous scans allowed.
	MaxConcurrent int
	// AcquireTimeout is how long to wait for a slot before returning an error.
	AcquireTimeout time.Duration
}

// DefaultConcurrencyConfig returns a ConcurrencyConfig with sensible defaults.
func DefaultConcurrencyConfig() ConcurrencyConfig {
	return ConcurrencyConfig{
		MaxConcurrent:  4,
		AcquireTimeout: 5 * time.Second,
	}
}

// ScanConcurrencyManager limits the number of scans running at the same time.
type ScanConcurrencyManager struct {
	sem     chan struct{}
	timeout time.Duration
	mu      sync.Mutex
	active  map[string]bool
}

// NewScanConcurrencyManager creates a ScanConcurrencyManager from cfg.
// Returns an error if MaxConcurrent < 1 or AcquireTimeout <= 0.
func NewScanConcurrencyManager(cfg ConcurrencyConfig) (*ScanConcurrencyManager, error) {
	if cfg.MaxConcurrent < 1 {
		return nil, errors.New("portwatch: MaxConcurrent must be at least 1")
	}
	if cfg.AcquireTimeout <= 0 {
		return nil, errors.New("portwatch: AcquireTimeout must be positive")
	}
	sem := make(chan struct{}, cfg.MaxConcurrent)
	for i := 0; i < cfg.MaxConcurrent; i++ {
		sem <- struct{}{}
	}
	return &ScanConcurrencyManager{
		sem:     sem,
		timeout: cfg.AcquireTimeout,
		active:  make(map[string]bool),
	}, nil
}

// Acquire attempts to reserve a scan slot for the given target.
// Returns ErrConcurrencyTimeout if no slot becomes available within the timeout.
// Returns ErrTargetAlreadyScanning if the target is already being scanned.
func (m *ScanConcurrencyManager) Acquire(target string) error {
	m.mu.Lock()
	if m.active[target] {
		m.mu.Unlock()
		return ErrTargetAlreadyScanning
	}
	m.mu.Unlock()

	select {
	case <-m.sem:
		m.mu.Lock()
		m.active[target] = true
		m.mu.Unlock()
		return nil
	case <-time.After(m.timeout):
		return ErrConcurrencyTimeout
	}
}

// Release frees the scan slot held by target. Calling Release for a target
// that is not currently active is a no-op.
func (m *ScanConcurrencyManager) Release(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.active[target] {
		return
	}
	delete(m.active, target)
	m.sem <- struct{}{}
}

// Active returns the set of targets currently holding a scan slot.
func (m *ScanConcurrencyManager) Active() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, 0, len(m.active))
	for t := range m.active {
		out = append(out, t)
	}
	return out
}

// Slots returns the number of free concurrency slots remaining.
func (m *ScanConcurrencyManager) Slots() int {
	return len(m.sem)
}
