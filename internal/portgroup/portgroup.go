// Package portgroup provides named groupings of ports for use in config and alerts.
package portgroup

import "fmt"

// Group is a named set of ports.
type Group struct {
	Name  string
	Ports []int
}

// Registry holds named port groups.
type Registry struct {
	groups map[string][]int
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string][]int)}
}

// Register adds or replaces a named group.
func (r *Registry) Register(name string, ports []int) error {
	if name == "" {
		return fmt.Errorf("portgroup: name must not be empty")
	}
	if len(ports) == 0 {
		return fmt.Errorf("portgroup: ports must not be empty for group %q", name)
	}
	copy := make([]int, len(ports))
	for i, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("portgroup: invalid port %d in group %q", p, name)
		}
		copy[i] = p
	}
	r.groups[name] = copy
	return nil
}

// Lookup returns the ports for a named group, or an error if not found.
func (r *Registry) Lookup(name string) ([]int, error) {
	ports, ok := r.groups[name]
	if !ok {
		return nil, fmt.Errorf("portgroup: group %q not found", name)
	}
	return ports, nil
}

// Names returns all registered group names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.groups))
	for n := range r.groups {
		names = append(names, n)
	}
	return names
}

// Resolve expands a slice of group names into a deduplicated port list.
func (r *Registry) Resolve(names []string) ([]int, error) {
	seen := make(map[int]struct{})
	var result []int
	for _, name := range names {
		ports, err := r.Lookup(name)
		if err != nil {
			return nil, err
		}
		for _, p := range ports {
			if _, ok := seen[p]; !ok {
				seen[p] = struct{}{}
				result = append(result, p)
			}
		}
	}
	return result, nil
}
