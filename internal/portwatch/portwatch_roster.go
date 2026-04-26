package portwatch

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// RosterEntry holds metadata about a tracked target.
type RosterEntry struct {
	Target    string
	AddedAt   time.Time
	LastSeen  time.Time
	Active    bool
}

// RosterManager tracks the set of known scan targets.
type RosterManager struct {
	mu      sync.RWMutex
	entries map[string]*RosterEntry
}

// NewRosterManager returns an empty RosterManager.
func NewRosterManager() *RosterManager {
	return &RosterManager{
		entries: make(map[string]*RosterEntry),
	}
}

// Register adds a target to the roster if not already present.
// Returns an error if target is empty.
func (r *RosterManager) Register(target string) error {
	if target == "" {
		return errors.New("portwatch: roster: target must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[target]; !ok {
		r.entries[target] = &RosterEntry{
			Target:  target,
			AddedAt: time.Now(),
			Active:  true,
		}
	}
	return nil
}

// Touch records that target was seen at the given time.
// Returns an error if target is not registered.
func (r *RosterManager) Touch(target string, at time.Time) error {
	if target == "" {
		return errors.New("portwatch: roster: target must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[target]
	if !ok {
		return errors.New("portwatch: roster: target not registered")
	}
	e.LastSeen = at
	return nil
}

// Deactivate marks a target as inactive without removing it.
func (r *RosterManager) Deactivate(target string) error {
	if target == "" {
		return errors.New("portwatch: roster: target must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[target]
	if !ok {
		return errors.New("portwatch: roster: target not registered")
	}
	e.Active = false
	return nil
}

// Get returns the RosterEntry for a target, or false if not found.
func (r *RosterManager) Get(target string) (RosterEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[target]
	if !ok {
		return RosterEntry{}, false
	}
	return *e, true
}

// All returns all roster entries sorted by target name.
func (r *RosterManager) All() []RosterEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]RosterEntry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target < out[j].Target })
	return out
}
