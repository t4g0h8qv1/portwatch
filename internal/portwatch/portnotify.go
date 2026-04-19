package portwatch

import (
	"fmt"
	"io"
	"sort"
)

// Event describes port changes detected during a scan cycle.
type Event struct {
	Host   string
	Opened []int
	Closed []int
}

func buildEvent(host string, before, after []int) Event {
	bSet := toSet(before)
	aSet := toSet(after)

	var opened, closed []int
	for p := range aSet {
		if !bSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range bSet {
		if !aSet[p] {
			closed = append(closed, p)
		}
	}
	sort.Ints(opened)
	sort.Ints(closed)
	return Event{Host: host, Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}

// HasChanges reports whether the event contains any port changes.
func HasChanges(e Event) bool {
	return len(e.Opened) > 0 || len(e.Closed) > 0
}

// Summary returns a one-line human-readable description of the event.
func Summary(e Event) string {
	if !HasChanges(e) {
		return fmt.Sprintf("%s: no changes", e.Host)
	}
	return fmt.Sprintf("%s: +%d opened, -%d closed", e.Host, len(e.Opened), len(e.Closed))
}

// WriteEvent writes a detailed event description to w.
func WriteEvent(w io.Writer, e Event) {
	fmt.Fprintf(w, "host: %s\n", e.Host)
	for _, p := range e.Opened {
		fmt.Fprintf(w, "  + %d\n", p)
	}
	for _, p := range e.Closed {
		fmt.Fprintf(w, "  - %d\n", p)
	}
}
