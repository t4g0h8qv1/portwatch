// Package portlock allows ports to be "locked" (expected to always be open)
// and raises an alert if a locked port is found closed.
package portlock

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Entry represents a locked port expectation.
type Entry struct {
	Port      int       `json:"port"`
	AddedAt   time.Time `json:"added_at"`
	Comment   string    `json:"comment,omitempty"`
}

// Store holds locked port entries persisted to disk.
type Store struct {
	mu      sync.RWMutex
	entries map[int]Entry
	path    string
}

// Load reads the lock file from path. If the file does not exist, an empty
// Store is returned.
func Load(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[int]Entry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		s.entries[e.Port] = e
	}
	return s, nil
}

// Lock adds a port to the locked set and persists the store.
func (s *Store) Lock(port int, comment string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = Entry{Port: port, AddedAt: time.Now(), Comment: comment}
	return s.save()
}

// Unlock removes a port from the locked set and persists the store.
func (s *Store) Unlock(port int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
	return s.save()
}

// IsLocked reports whether port is in the locked set.
func (s *Store) IsLocked(port int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.entries[port]
	return ok
}

// Missing returns locked ports that are absent from open.
func (s *Store) Missing(open []int) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	openSet := make(map[int]struct{}, len(open))
	for _, p := range open {
		openSet[p] = struct{}{}
	}
	var missing []int
	for port := range s.entries {
		if _, found := openSet[port]; !found {
			missing = append(missing, port)
		}
	}
	return missing
}

func (s *Store) save() error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
