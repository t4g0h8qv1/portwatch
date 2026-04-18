// Package portinventory maintains a named inventory of hosts and their
// expected open ports, enabling policy-driven comparisons.
package portinventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

// Entry holds the expected ports for a single host.
type Entry struct {
	Host  string `json:"host"`
	Ports []int  `json:"ports"`
}

// Inventory maps host names to their expected port sets.
type Inventory struct {
	mu      sync.RWMutex
	entries map[string][]int
}

// New returns an empty Inventory.
func New() *Inventory {
	return &Inventory{entries: make(map[string][]int)}
}

// Set registers or replaces the expected ports for a host.
func (inv *Inventory) Set(host string, ports []int) error {
	if host == "" {
		return errors.New("portinventory: host must not be empty")
	}
	cp := make([]int, len(ports))
	copy(cp, ports)
	sort.Ints(cp)
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.entries[host] = cp
	return nil
}

// Get returns the expected ports for a host and whether the host exists.
func (inv *Inventory) Get(host string) ([]int, bool) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	p, ok := inv.entries[host]
	return p, ok
}

// Hosts returns a sorted list of all registered hosts.
func (inv *Inventory) Hosts() []string {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	out := make([]string, 0, len(inv.entries))
	for h := range inv.entries {
		out = append(out, h)
	}
	sort.Strings(out)
	return out
}

// Save persists the inventory to a JSON file at path.
func (inv *Inventory) Save(path string) error {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	entries := make([]Entry, 0, len(inv.entries))
	for h, p := range inv.entries {
		entries = append(entries, Entry{Host: h, Ports: p})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Host < entries[j].Host })
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("portinventory: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a JSON inventory file and returns a populated Inventory.
func Load(path string) (*Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return nil, fmt.Errorf("portinventory: read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("portinventory: unmarshal: %w", err)
	}
	inv := New()
	for _, e := range entries {
		if err := inv.Set(e.Host, e.Ports); err != nil {
			return nil, err
		}
	}
	return inv, nil
}
