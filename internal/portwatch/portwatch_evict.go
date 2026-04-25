package portwatch

import (
	"errors"
	"sync"
	"time"
)

// EvictConfig holds configuration for the eviction manager.
type EvictConfig struct {
	// MaxAge is the maximum time a target may go unscanned before being evicted.
	MaxAge time.Duration
}

// DefaultEvictConfig returns a conservative default eviction configuration.
func DefaultEvictConfig() EvictConfig {
	return EvictConfig{
		MaxAge: 24 * time.Hour,
	}
}

// evictEntry tracks the last activity time for a target.
type evictEntry struct {
	lastSeen time.Time
}

// ScanEvictManager evicts targets that have not been seen within MaxAge.
type ScanEvictManager struct {
	mu      sync.Mutex
	cfg     EvictConfig
	entries map[string]evictEntry
	now     func() time.Time
}

// NewEvictManager creates a new ScanEvictManager with the given config.
func NewEvictManager(cfg EvictConfig) (*ScanEvictManager, error) {
	if cfg.MaxAge <= 0 {
		return nil, errors.New("evict: MaxAge must be positive")
	}
	return &ScanEvictManager{
		cfg:     cfg,
		entries: make(map[string]evictEntry),
		now:     time.Now,
	}, nil
}

// Touch records that a target was active at the current time.
func (m *ScanEvictManager) Touch(target string) error {
	if target == "" {
		return errors.New("evict: target must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[target] = evictEntry{lastSeen: m.now()}
	return nil
}

// ShouldEvict reports whether the target has exceeded MaxAge without activity.
func (m *ScanEvictManager) ShouldEvict(target string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[target]
	if !ok {
		return false
	}
	return m.now().Sub(e.lastSeen) > m.cfg.MaxAge
}

// Evict removes a target from the manager.
func (m *ScanEvictManager) Evict(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}

// Targets returns all currently tracked targets.
func (m *ScanEvictManager) Targets() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, 0, len(m.entries))
	for t := range m.entries {
		out = append(out, t)
	}
	return out
}
