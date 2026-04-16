// Package timeout provides configurable per-host scan timeout management.
package timeout

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidTimeout is returned when a non-positive timeout is provided.
var ErrInvalidTimeout = errors.New("timeout: duration must be greater than zero")

// Manager holds per-host timeout overrides with a global default.
type Manager struct {
	mu       sync.RWMutex
	default_ time.Duration
	overrides map[string]time.Duration
}

// New creates a Manager with the given default timeout.
func New(defaultTimeout time.Duration) (*Manager, error) {
	if defaultTimeout <= 0 {
		return nil, ErrInvalidTimeout
	}
	return &Manager{
		default_:  defaultTimeout,
		overrides: make(map[string]time.Duration),
	}, nil
}

// Set registers a timeout override for a specific host.
func (m *Manager) Set(host string, d time.Duration) error {
	if d <= 0 {
		return ErrInvalidTimeout
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.overrides[host] = d
	return nil
}

// Get returns the effective timeout for a host.
func (m *Manager) Get(host string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if d, ok := m.overrides[host]; ok {
		return d
	}
	return m.default_
}

// Remove deletes a host override, reverting to the default.
func (m *Manager) Remove(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.overrides, host)
}

// Default returns the global default timeout.
func (m *Manager) Default() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.default_
}
