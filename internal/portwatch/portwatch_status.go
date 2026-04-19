package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Status holds a point-in-time summary of the portwatch run loop.
type Status struct {
	Target    string
	Ports     []int
	LastScan  time.Time
	LastAlert time.Time
	UpSince   time.Time
	ScanCount int
	AlertCount int
	ErrorCount int
	LastError  error
}

// IsHealthy returns true when the last scan completed without error.
func (s Status) IsHealthy() bool {
	return s.LastError == nil
}

// WriteStatus formats a Status as a human-readable table to w.
func WriteStatus(w io.Writer, s Status) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Target:\t%s\n", s.Target)
	fmt.Fprintf(tw, "Monitored ports:\t%d\n", len(s.Ports))
	fmt.Fprintf(tw, "Up since:\t%s\n", formatStatusTime(s.UpSince))
	fmt.Fprintf(tw, "Last scan:\t%s\n", formatStatusTime(s.LastScan))
	fmt.Fprintf(tw, "Last alert:\t%s\n", formatStatusTime(s.LastAlert))
	fmt.Fprintf(tw, "Scans:\t%d\n", s.ScanCount)
	fmt.Fprintf(tw, "Alerts:\t%d\n", s.AlertCount)
	fmt.Fprintf(tw, "Errors:\t%d\n", s.ErrorCount)
	healthy := "yes"
	if !s.IsHealthy() {
		healthy = fmt.Sprintf("no (%v)", s.LastError)
	}
	fmt.Fprintf(tw, "Healthy:\t%s\n", healthy)
	tw.Flush()
}

func formatStatusTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}
