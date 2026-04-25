package portwatch

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

// ScorecardEntry holds cumulative scan statistics for a single target.
type ScorecardEntry struct {
	Target       string
	TotalScans   int
	TotalAlerts  int
	TotalErrors  int
	LastScan     time.Time
	LastAlert    time.Time
}

// ScorecardManager tracks per-target scan outcomes over the lifetime of the
// process. It is safe for concurrent use.
type ScorecardManager struct {
	mu      sync.RWMutex
	entries map[string]*ScorecardEntry
}

// NewScorecardManager returns an initialised ScorecardManager.
func NewScorecardManager() *ScorecardManager {
	return &ScorecardManager{
		entries: make(map[string]*ScorecardEntry),
	}
}

// RecordScan increments the scan counter for target and updates LastScan.
func (s *ScorecardManager) RecordScan(target string, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.ensure(target)
	e.TotalScans++
	e.LastScan = at
}

// RecordAlert increments the alert counter for target and updates LastAlert.
func (s *ScorecardManager) RecordAlert(target string, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.ensure(target)
	e.TotalAlerts++
	e.LastAlert = at
}

// RecordError increments the error counter for target.
func (s *ScorecardManager) RecordError(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensure(target).TotalErrors++
}

// Get returns a copy of the ScorecardEntry for target, or false if unknown.
func (s *ScorecardManager) Get(target string) (ScorecardEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[target]
	if !ok {
		return ScorecardEntry{}, false
	}
	return *e, true
}

// All returns a snapshot of all entries sorted by target name.
func (s *ScorecardManager) All() []ScorecardEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ScorecardEntry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target < out[j].Target })
	return out
}

// WriteScorecardTable writes a human-readable scorecard table to w.
func WriteScorecardTable(w io.Writer, entries []ScorecardEntry) {
	fmt.Fprintf(w, "%-30s %8s %8s %8s  %-20s\n", "TARGET", "SCANS", "ALERTS", "ERRORS", "LAST SCAN")
	for _, e := range entries {
		ls := "never"
		if !e.LastScan.IsZero() {
			ls = e.LastScan.Format(time.RFC3339)
		}
		fmt.Fprintf(w, "%-30s %8d %8d %8d  %-20s\n", e.Target, e.TotalScans, e.TotalAlerts, e.TotalErrors, ls)
	}
}

func (s *ScorecardManager) ensure(target string) *ScorecardEntry {
	if e, ok := s.entries[target]; ok {
		return e
	}
	e := &ScorecardEntry{Target: target}
	s.entries[target] = e
	return e
}
