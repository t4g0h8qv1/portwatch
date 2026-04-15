// Package filter provides port filtering utilities for portwatch.
// It allows callers to include or exclude specific ports from scan results
// based on allowlists and denylists.
package filter

// Options holds the configuration for filtering ports.
type Options struct {
	// Allow is an explicit set of ports that are permitted (allowlist).
	// If non-empty, only ports in this set pass the filter.
	Allow []int

	// Deny is a set of ports that are always rejected (denylist).
	Deny []int
}

// Filter applies the given Options to a slice of ports and returns
// the ports that pass the filter rules.
//
// Rules applied in order:
//  1. If a port appears in Deny, it is excluded.
//  2. If Allow is non-empty, only ports in Allow are included.
//  3. Otherwise the port is included.
func Filter(ports []int, opts Options) []int {
	denySet := toSet(opts.Deny)
	allowSet := toSet(opts.Allow)

	result := make([]int, 0, len(ports))
	for _, p := range ports {
		if denySet[p] {
			continue
		}
		if len(allowSet) > 0 && !allowSet[p] {
			continue
		}
		result = append(result, p)
	}
	return result
}

// toSet converts a slice of ints into a map for O(1) lookups.
func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
