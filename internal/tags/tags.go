// Package tags provides port tagging — associating human-readable labels
// with port numbers for richer output and alerting.
package tags

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Registry maps port numbers to tag labels.
type Registry struct {
	mu   sync.RWMutex
	tags map[int]string
}

// New returns a Registry pre-populated with common well-known port labels.
func New() *Registry {
	return &Registry{
		tags: map[int]string{
			21:   "ftp",
			22:   "ssh",
			23:   "telnet",
			25:   "smtp",
			53:   "dns",
			80:   "http",
			110:  "pop3",
			143:  "imap",
			443:  "https",
			3306: "mysql",
			5432: "postgres",
			6379: "redis",
			8080: "http-alt",
			8443: "https-alt",
			27017: "mongodb",
		},
	}
}

// Set adds or updates a tag for the given port.
func (r *Registry) Set(port int, label string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[port] = label
}

// Get returns the tag for a port and whether one exists.
func (r *Registry) Get(port int) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	label, ok := r.tags[port]
	return label, ok
}

// Label returns the tag label for a port, or a default "port/<n>" string.
func (r *Registry) Label(port int) string {
	if label, ok := r.Get(port); ok {
		return label
	}
	return fmt.Sprintf("port/%d", port)
}

// LoadFile merges tags from a JSON file (map of port string -> label).
func (r *Registry) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("tags: open %s: %w", path, err)
	}
	defer f.Close()

	var raw map[int]string
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return fmt.Errorf("tags: decode %s: %w", path, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for port, label := range raw {
		r.tags[port] = label
	}
	return nil
}
