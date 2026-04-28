package portwatch

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// SignalKind classifies the type of scan signal emitted.
type SignalKind string

const (
	SignalOpened  SignalKind = "opened"
	SignalClosed  SignalKind = "closed"
	SignalStable  SignalKind = "stable"
	SignalUnknown SignalKind = "unknown"
)

// ScanSignal represents a discrete event observed during a port scan.
type ScanSignal struct {
	Target    string
	Port      int
	Kind      SignalKind
	ObservedAt time.Time
}

// SignalManager records and retrieves scan signals per target.
type SignalManager struct {
	mu      sync.RWMutex
	records map[string][]ScanSignal
	maxAge  time.Duration
}

// DefaultSignalConfig returns sensible defaults for the signal manager.
var DefaultSignalConfig = struct {
	MaxAge time.Duration
}{
	MaxAge: 24 * time.Hour,
}

// NewSignalManager constructs a SignalManager with the given max retention age.
func NewSignalManager(maxAge time.Duration) (*SignalManager, error) {
	if maxAge <= 0 {
		return nil, errors.New("portwatch: signal maxAge must be positive")
	}
	return &SignalManager{
		records: make(map[string][]ScanSignal),
		maxAge:  maxAge,
	}, nil
}

// Record appends a signal for the given target, pruning entries older than maxAge.
func (m *SignalManager) Record(target string, port int, kind SignalKind) error {
	if target == "" {
		return errors.New("portwatch: target must not be empty")
	}
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	cutoff := now.Add(-m.maxAge)
	existing := m.records[target]
	pruned := existing[:0]
	for _, s := range existing {
		if s.ObservedAt.After(cutoff) {
			pruned = append(pruned, s)
		}
	}
	pruned = append(pruned, ScanSignal{
		Target:     target,
		Port:       port,
		Kind:       kind,
		ObservedAt: now,
	})
	m.records[target] = pruned
	return nil
}

// All returns all signals for the given target, sorted by observation time.
func (m *SignalManager) All(target string) []ScanSignal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	src := m.records[target]
	out := make([]ScanSignal, len(src))
	copy(out, src)
	sort.Slice(out, func(i, j int) bool {
		return out[i].ObservedAt.Before(out[j].ObservedAt)
	})
	return out
}

// Targets returns all targets that have at least one recorded signal.
func (m *SignalManager) Targets() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, 0, len(m.records))
	for t := range m.records {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}
