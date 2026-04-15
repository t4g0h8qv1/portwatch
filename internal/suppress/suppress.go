// Package suppress provides a mechanism to silence alerts for known or
// expected open ports, preventing repeated notifications for ports that
// have been explicitly acknowledged by the operator.
package suppress

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single suppression rule for a port.
type Entry struct {
	Port      int       `json:"port"`
	Reason    string    `json:"reason"`
	ExpiresAt time.Time `json:"expires_at"`
}

// List holds all active suppression entries.
type List struct {
	Entries []Entry `json:"entries"`
	path    string
}

// Load reads a suppression list from disk. If the file does not exist,
// an empty list is returned without error.
func Load(path string) (*List, error) {
	l := &List{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, l); err != nil {
		return nil, err
	}
	return l, nil
}

// Add inserts a new suppression entry and persists the list to disk.
func (l *List) Add(port int, reason string, ttl time.Duration) error {
	l.prune()
	l.Entries = append(l.Entries, Entry{
		Port:      port,
		Reason:    reason,
		ExpiresAt: time.Now().Add(ttl),
	})
	return l.save()
}

// IsSuppressed reports whether the given port has an active (non-expired)
// suppression entry.
func (l *List) IsSuppressed(port int) bool {
	now := time.Now()
	for _, e := range l.Entries {
		if e.Port == port && now.Before(e.ExpiresAt) {
			return true
		}
	}
	return false
}

// Filter removes suppressed ports from the provided slice.
func (l *List) Filter(ports []int) []int {
	var out []int
	for _, p := range ports {
		if !l.IsSuppressed(p) {
			out = append(out, p)
		}
	}
	return out
}

// prune removes expired entries from memory (does not persist).
func (l *List) prune() {
	now := time.Now()
	var active []Entry
	for _, e := range l.Entries {
		if now.Before(e.ExpiresAt) {
			active = append(active, e)
		}
	}
	l.Entries = active
}

func (l *List) save() error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
