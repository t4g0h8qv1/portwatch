package portwatch

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// SnapshotEntry holds a point-in-time record of open ports for a target.
type SnapshotEntry struct {
	Target    string
	Ports     []int
	TakenAt   time.Time
}

// SnapshotStore retains the most recent scan snapshot per target.
type SnapshotStore struct {
	mu      sync.RWMutex
	entries map[string]SnapshotEntry
}

// NewSnapshotStore returns an initialised SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{entries: make(map[string]SnapshotEntry)}
}

// Record stores or replaces the snapshot for target.
func (s *SnapshotStore) Record(target string, ports []int) {
	if target == "" {
		return
	}
	copy := make([]int, len(ports))
	for i, p := range ports {
		copy[i] = p
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[target] = SnapshotEntry{
		Target:  target,
		Ports:   copy,
		TakenAt: time.Now(),
	}
}

// Get returns the latest snapshot for target, or false if none exists.
func (s *SnapshotStore) Get(target string) (SnapshotEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[target]
	return e, ok
}

// Targets returns all tracked target names.
func (s *SnapshotStore) Targets() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.entries))
	for t := range s.entries {
		out = append(out, t)
	}
	return out
}

// WriteSnapshotTable writes a human-readable snapshot summary to w.
func WriteSnapshotTable(w io.Writer, store *SnapshotStore) {
	targets := store.Targets()
	fmt.Fprintf(w, "%-30s %-8s %s\n", "TARGET", "PORTS", "TAKEN AT")
	fmt.Fprintf(w, "%-30s %-8s %s\n", "------", "-----", "--------")
	for _, t := range targets {
		e, _ := store.Get(t)
		fmt.Fprintf(w, "%-30s %-8d %s\n", e.Target, len(e.Ports), e.TakenAt.Format(time.RFC3339))
	}
}
