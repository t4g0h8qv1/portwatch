// Package portdiff computes human-readable diffs between two port scan results.
package portdiff

import "sort"

// Entry represents a single port change in a diff.
type Entry struct {
	Port   int
	Status string // "opened" or "closed"
}

// Result holds the full diff between two scans.
type Result struct {
	Opened []Entry
	Closed []Entry
}

// HasChanges reports whether any ports changed.
func (r Result) HasChanges() bool {
	return len(r.Opened) > 0 || len(r.Closed) > 0
}

// Compute returns the diff between previous and current port lists.
func Compute(previous, current []int) Result {
	prev := toSet(previous)
	curr := toSet(current)

	var opened, closed []Entry

	for p := range curr {
		if !prev[p] {
			opened = append(opened, Entry{Port: p, Status: "opened"})
		}
	}
	for p := range prev {
		if !curr[p] {
			closed = append(closed, Entry{Port: p, Status: "closed"})
		}
	}

	sort.Slice(opened, func(i, j int) bool { return opened[i].Port < opened[j].Port })
	sort.Slice(closed, func(i, j int) bool { return closed[i].Port < closed[j].Port })

	return Result{Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
