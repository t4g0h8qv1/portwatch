// Package portmatch provides pattern-based port matching against named rules.
package portmatch

import (
	"fmt"
	"sort"
)

// Rule associates a name and description with a set of ports.
type Rule struct {
	Name  string
	Ports map[int]struct{}
	Desc  string
}

// Matcher holds a collection of named port rules.
type Matcher struct {
	rules map[string]*Rule
}

// New returns an empty Matcher.
func New() *Matcher {
	return &Matcher{rules: make(map[string]*Rule)}
}

// Add registers a named rule with the given ports.
func (m *Matcher) Add(name, desc string, ports []int) error {
	if name == "" {
		return fmt.Errorf("portmatch: rule name must not be empty")
	}
	if len(ports) == 0 {
		return fmt.Errorf("portmatch: rule %q has no ports", name)
	}
	set := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("portmatch: invalid port %d in rule %q", p, name)
		}
		set[p] = struct{}{}
	}
	m.rules[name] = &Rule{Name: name, Desc: desc, Ports: set}
	return nil
}

// Match returns the names of all rules that contain the given port.
func (m *Matcher) Match(port int) []string {
	var matched []string
	for name, rule := range m.rules {
		if _, ok := rule.Ports[port]; ok {
			matched = append(matched, name)
		}
	}
	sort.Strings(matched)
	return matched
}

// MatchAll returns a map of port -> rule names for every port in the input.
func (m *Matcher) MatchAll(ports []int) map[int][]string {
	result := make(map[int][]string, len(ports))
	for _, p := range ports {
		if names := m.Match(p); len(names) > 0 {
			result[p] = names
		}
	}
	return result
}

// Rules returns all registered rule names in sorted order.
func (m *Matcher) Rules() []string {
	names := make([]string, 0, len(m.rules))
	for n := range m.rules {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
