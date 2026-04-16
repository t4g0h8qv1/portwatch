// Package portstate tracks the open/closed state of ports across scans
// and provides a simple diff between two snapshots.
package portstate

import "sort"

// State represents a snapshot of open ports for a host.
type State struct {
	Host  string
	Ports []int
}

// Diff holds the ports that changed between two states.
type Diff struct {
	Opened []int
	Closed []int
}

// HasChanges returns true when at least one port opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare returns a Diff between a previous and current State.
// Both slices in the result are sorted ascending.
func Compare(prev, curr State) Diff {
	prevSet := toSet(prev.Ports)
	currSet := toSet(curr.Ports)

	var opened, closed []int

	for p := range currSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			closed = append(closed, p)
		}
	}

	sort.Ints(opened)
	sort.Ints(closed)

	return Diff{Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
