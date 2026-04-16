package portdiff

import (
	"fmt"
	"io"
	"strings"
)

// Format writes a human-readable diff to w.
func Format(w io.Writer, r Result) {
	if !r.HasChanges() {
		fmt.Fprintln(w, "no port changes detected")
		return
	}
	for _, e := range r.Opened {
		fmt.Fprintf(w, "+ port %d %s\n", e.Port, e.Status)
	}
	for _, e := range r.Closed {
		fmt.Fprintf(w, "- port %d %s\n", e.Port, e.Status)
	}
}

// Summary returns a one-line summary string.
func Summary(r Result) string {
	if !r.HasChanges() {
		return "no changes"
	}
	parts := []string{}
	if n := len(r.Opened); n > 0 {
		parts = append(parts, fmt.Sprintf("%d opened", n))
	}
	if n := len(r.Closed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d closed", n))
	}
	return strings.Join(parts, ", ")
}
