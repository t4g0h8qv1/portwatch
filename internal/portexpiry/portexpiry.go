// Package portexpiry tracks how long ports have been continuously open
// and emits warnings when a port exceeds a configured maximum age.
package portexpiry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry records when a port was first seen open.
type Entry struct {
	Port      int       `json:"port"`
	FirstSeen time.Time `json:"first_seen"`
}

// Registry maps port numbers to their first-seen timestamps.
type Registry struct {
	path    string
	entries map[int]Entry
}

// Load reads the registry from disk, or returns an empty one if the file is missing.
func Load(path string) (*Registry, error) {
	r := &Registry{path: path, entries: make(map[int]Entry)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return r, nil
	}
	if err != nil {
		return nil, fmt.Errorf("portexpiry: read %s: %w", path, err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("portexpiry: parse %s: %w", path, err)
	}
	for _, e := range entries {
		r.entries[e.Port] = e
	}
	return r, nil
}

// Track updates the registry with the current set of open ports.
// New ports are recorded with now as their first-seen time.
// Ports no longer open are removed.
func (r *Registry) Track(openPorts []int, now time.Time) {
	next := make(map[int]Entry, len(openPorts))
	for _, p := range openPorts {
		if e, ok := r.entries[p]; ok {
			next[p] = e
		} else {
			next[p] = Entry{Port: p, FirstSeen: now}
		}
	}
	r.entries = next
}

// Expired returns ports whose first-seen age exceeds maxAge.
func (r *Registry) Expired(maxAge time.Duration, now time.Time) []Entry {
	var out []Entry
	for _, e := range r.entries {
		if now.Sub(e.FirstSeen) > maxAge {
			out = append(out, e)
		}
	}
	return out
}

// Save persists the registry to disk.
func (r *Registry) Save() error {
	entries := make([]Entry, 0, len(r.entries))
	for _, e := range r.entries {
		entries = append(entries, e)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("portexpiry: marshal: %w", err)
	}
	if err := os.WriteFile(r.path, data, 0o600); err != nil {
		return fmt.Errorf("portexpiry: write %s: %w", r.path, err)
	}
	return nil
}
