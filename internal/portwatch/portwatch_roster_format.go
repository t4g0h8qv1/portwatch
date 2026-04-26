package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteRosterTable writes a formatted table of roster entries to w.
func WriteRosterTable(w io.Writer, entries []RosterEntry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tACTIVE\tADDED\tLAST SEEN")
	for _, e := range entries {
		active := "yes"
		if !e.Active {
			active = "no"
		}
		added := e.AddedAt.Format(time.RFC3339)
		lastSeen := "never"
		if !e.LastSeen.IsZero() {
			lastSeen = e.LastSeen.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.Target, active, added, lastSeen)
	}
	_ = tw.Flush()
}

// RosterSummary returns a human-readable summary of roster state.
func RosterSummary(entries []RosterEntry) string {
	if len(entries) == 0 {
		return "roster: no targets registered"
	}
	active := 0
	for _, e := range entries {
		if e.Active {
			active++
		}
	}
	inactive := len(entries) - active
	if inactive == 0 {
		return fmt.Sprintf("roster: %d target(s), all active", len(entries))
	}
	return fmt.Sprintf("roster: %d target(s), %d active, %d inactive", len(entries), active, inactive)
}
