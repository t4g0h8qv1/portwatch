// Package portprofile groups named port profiles that can be referenced
// in configuration to describe expected open ports for a host role.
package portprofile

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Profile is a named set of ports.
type Profile struct {
	Name  string
	Ports []int
}

// Registry holds named port profiles.
type Registry struct {
	mu       sync.RWMutex
	profiles map[string]Profile
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{profiles: make(map[string]Profile)}
}

// Register adds or replaces a profile. Name must be non-empty and ports non-nil.
func (r *Registry) Register(name string, ports []int) error {
	if name == "" {
		return errors.New("portprofile: name must not be empty")
	}
	if len(ports) == 0 {
		return errors.New("portprofile: ports must not be empty")
	}
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("portprofile: invalid port %d", p)
		}
	}
	deduped := dedup(ports)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profiles[name] = Profile{Name: name, Ports: deduped}
	return nil
}

// Get returns the profile for name, or an error if not found.
func (r *Registry) Get(name string) (Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("portprofile: profile %q not found", name)
	}
	return p, nil
}

// Names returns all registered profile names in sorted order.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.profiles))
	for n := range r.profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func dedup(ports []int) []int {
	seen := make(map[int]struct{}, len(ports))
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	sort.Ints(out)
	return out
}
