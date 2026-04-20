package portwatch

import (
	"sync"
	"time"
)

// ShadowEntry records a port that appeared unexpectedly but has not yet
// been promoted to an alert — it is being "shadowed" for confirmation.
type ShadowEntry struct {
	Port      int
	FirstSeen time.Time
	Count     int
}

// ShadowTracker holds ports that have been seen fewer than MinObservations
// times. Once a port crosses the threshold it is considered confirmed.
type ShadowTracker struct {
	mu              sync.Mutex
	entries         map[string]map[int]*ShadowEntry // target -> port -> entry
	MinObservations int
	MaxAge          time.Duration
	now             func() time.Time
}

// NewShadowTracker returns a ShadowTracker that requires minObs observations
// within maxAge before a port is considered confirmed.
func NewShadowTracker(minObs int, maxAge time.Duration) (*ShadowTracker, error) {
	if minObs < 1 {
		return nil, errShadowInvalidObs
	}
	if maxAge <= 0 {
		return nil, errShadowInvalidAge
	}
	return &ShadowTracker{
		entries:         make(map[string]map[int]*ShadowEntry),
		MinObservations: minObs,
		MaxAge:          maxAge,
		now:             time.Now,
	}, nil
}

// Observe records a sighting of port on target. Returns true when the port
// has been seen at least MinObservations times and should be alerted on.
func (s *ShadowTracker) Observe(target string, port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune(target)
	if s.entries[target] == nil {
		s.entries[target] = make(map[int]*ShadowEntry)
	}
	e, ok := s.entries[target][port]
	if !ok {
		e = &ShadowEntry{Port: port, FirstSeen: s.now()}
		s.entries[target][port] = e
	}
	e.Count++
	if e.Count >= s.MinObservations {
		delete(s.entries[target], port)
		return true
	}
	return false
}

// Pending returns all shadow entries for target that have not yet confirmed.
func (s *ShadowTracker) Pending(target string) []ShadowEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune(target)
	out := make([]ShadowEntry, 0, len(s.entries[target]))
	for _, e := range s.entries[target] {
		out = append(out, *e)
	}
	return out
}

// prune removes expired entries for target. Caller must hold s.mu.
func (s *ShadowTracker) prune(target string) {
	now := s.now()
	for port, e := range s.entries[target] {
		if now.Sub(e.FirstSeen) > s.MaxAge {
			delete(s.entries[target], port)
		}
	}
}
