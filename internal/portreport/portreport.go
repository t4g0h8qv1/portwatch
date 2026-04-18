// Package portreport aggregates scan results into a structured summary
// suitable for display or downstream processing.
package portreport

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Entry holds the details for a single port in a report.
type Entry struct {
	Port     int
	Label    string
	Status   string // "open", "new", "closed"
	Severity string
}

// Report is a point-in-time summary of port scan results.
type Report struct {
	Host      string
	ScannedAt time.Time
	Entries   []Entry
}

// New creates a Report for the given host and entries.
func New(host string, entries []Entry) Report {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Port < sorted[j].Port
	})
	return Report{
		Host:      host,
		ScannedAt: time.Now(),
		Entries:   sorted,
	}
}

// CountByStatus returns a map of status -> count.
func (r Report) CountByStatus() map[string]int {
	m := make(map[string]int)
	for _, e := range r.Entries {
		m[e.Status]++
	}
	return m
}

// WriteSummary writes a human-readable summary to w.
func (r Report) WriteSummary(w io.Writer) {
	counts := r.CountByStatus()
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "Scanned: %s\n", r.ScannedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Total ports: %d (new: %d, closed: %d, open: %d)\n",
		len(r.Entries), counts["new"], counts["closed"], counts["open"])
	for _, e := range r.Entries {
		fmt.Fprintf(w, "  [%s] %d %s (%s)\n", e.Status, e.Port, e.Label, e.Severity)
	}
}
