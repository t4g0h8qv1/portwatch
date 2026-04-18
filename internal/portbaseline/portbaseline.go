// Package portbaseline manages per-host port baselines with versioning.
package portbaseline

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"
)

// Entry holds a versioned baseline for a single host.
type Entry struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// Store manages baselines keyed by host.
type Store struct {
	path    string
	entries map[string]*Entry
}

// Load reads the store from disk. Missing file returns an empty store.
func Load(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[string]*Entry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Set updates the baseline for host, incrementing the version.
func (s *Store) Set(host string, ports []int) {
	deduped := dedup(ports)
	sort.Ints(deduped)
	e, ok := s.entries[host]
	if !ok {
		e = &Entry{Host: host}
		s.entries[host] = e
	}
	e.Ports = deduped
	e.UpdatedAt = time.Now().UTC()
	e.Version++
}

// Get returns the baseline entry for host, or false if not found.
func (s *Store) Get(host string) (Entry, bool) {
	e, ok := s.entries[host]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Hosts returns all tracked hosts in sorted order.
func (s *Store) Hosts() []string {
	out := make([]string, 0, len(s.entries))
	for h := range s.entries {
		out = append(out, h)
	}
	sort.Strings(out)
	return out
}

// Save persists the store to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func dedup(ports []int) []int {
	seen := make(map[int]struct{}, len(ports))
	out := ports[:0:0]
	for _, p := range ports {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}
