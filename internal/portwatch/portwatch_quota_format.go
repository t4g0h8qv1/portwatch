package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// QuotaStatus holds a snapshot of quota state for a single target.
type QuotaStatus struct {
	Target    string
	Remaining int
	Max       int
	Throttled bool
}

// WriteQuotaTable writes a formatted table of quota statuses to w.
func WriteQuotaTable(w io.Writer, statuses []QuotaStatus) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tREMAINING\tMAX\tSTATUS")
	for _, s := range statuses {
		status := "ok"
		if s.Throttled {
			status = "throttled"
		}
		fmt.Fprintf(tw, "%s\t%d\t%d\t%s\n", s.Target, s.Remaining, s.Max, status)
	}
	tw.Flush()
}

// QuotaSummary returns a human-readable summary of quota statuses.
func QuotaSummary(statuses []QuotaStatus) string {
	if len(statuses) == 0 {
		return "no quota entries"
	}
	throttled := 0
	for _, s := range statuses {
		if s.Throttled {
			throttled++
		}
	}
	if throttled == 0 {
		return fmt.Sprintf("%d target(s) within quota", len(statuses))
	}
	return fmt.Sprintf("%d/%d target(s) throttled", throttled, len(statuses))
}
