package portwatch

import (
	"fmt"
	"io"
	"strings"
)

// ChangeEvent describes a set of port changes detected during a scan.
type ChangeEvent struct {
	Host    string
	Opened  []int
	Closed  []int
}

// HasChanges returns true if there are any opened or closed ports.
func (e ChangeEvent) HasChanges() bool {
	return len(e.Opened) > 0 || len(e.Closed) > 0
}

// Summary returns a one-line human-readable description of the event.
func (e ChangeEvent) Summary() string {
	parts := []string{}
	if len(e.Opened) > 0 {
		parts = append(parts, fmt.Sprintf("%d opened", len(e.Opened)))
	}
	if len(e.Closed) > 0 {
		parts = append(parts, fmt.Sprintf("%d closed", len(e.Closed)))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s: no changes", e.Host)
	}
	return fmt.Sprintf("%s: %s", e.Host, strings.Join(parts, ", "))
}

// WriteEvent writes a formatted change event to w.
func WriteEvent(w io.Writer, e ChangeEvent) {
	fmt.Fprintln(w, e.Summary())
	for _, p := range e.Opened {
		fmt.Fprintf(w, "  + %d\n", p)
	}
	for _, p := range e.Closed {
		fmt.Fprintf(w, "  - %d\n", p)
	}
}
