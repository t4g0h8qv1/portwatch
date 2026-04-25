package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DefaultSequenceConfig returns a SequenceConfig with sensible defaults.
func DefaultSequenceConfig() SequenceConfig {
	return SequenceConfig{
		MaxGap: 10 * time.Minute,
	}
}

// SequenceConfig controls how scan sequences are tracked.
type SequenceConfig struct {
	// MaxGap is the maximum time between scans before a sequence is reset.
	MaxGap time.Duration
}

func (c SequenceConfig) validate() error {
	if c.MaxGap <= 0 {
		return errors.New("portwatch: sequence MaxGap must be positive")
	}
	return nil
}

// SequenceEntry records the current scan sequence for a target.
type SequenceEntry struct {
	Target    string
	Count     int
	StartedAt time.Time
	LastSeen  time.Time
}

// ScanSequenceManager tracks consecutive scan counts per target,
// resetting the counter when the gap between scans exceeds MaxGap.
type ScanSequenceManager struct {
	mu      sync.Mutex
	cfg     SequenceConfig
	entries map[string]*SequenceEntry
	now     func() time.Time
}

// NewScanSequenceManager creates a new ScanSequenceManager with the given config.
func NewScanSequenceManager(cfg SequenceConfig) (*ScanSequenceManager, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &ScanSequenceManager{
		cfg:     cfg,
		entries: make(map[string]*SequenceEntry),
		now:     time.Now,
	}, nil
}

// Record marks a scan for the given target, incrementing or resetting
// the sequence counter depending on how long since the last scan.
func (m *ScanSequenceManager) Record(target string) SequenceEntry {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()
	e, ok := m.entries[target]
	if !ok || now.Sub(e.LastSeen) > m.cfg.MaxGap {
		e = &SequenceEntry{Target: target, Count: 0, StartedAt: now}
		m.entries[target] = e
	}
	e.Count++
	e.LastSeen = now
	return *e
}

// Get returns the current sequence entry for a target, and false if not found.
func (m *ScanSequenceManager) Get(target string) (SequenceEntry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[target]
	if !ok {
		return SequenceEntry{}, false
	}
	return *e, true
}

// Reset clears the sequence for the given target.
func (m *ScanSequenceManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}

// Targets returns all targets currently tracked.
func (m *ScanSequenceManager) Targets() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, 0, len(m.entries))
	for t := range m.entries {
		out = append(out, t)
	}
	return out
}
