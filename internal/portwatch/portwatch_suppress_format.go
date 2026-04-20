package portwatch

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// SuppressStatus holds display info for a suppressed target.
type SuppressStatus struct {
	Target string
	Expiry time.Time
}

// WriteSuppressTable writes a formatted table of suppressed targets to w.
func WriteSuppressTable(w io.Writer, statuses []SuppressStatus) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tSUPPRESSED UNTIL")
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Target < statuses[j].Target
	})
	for _, s := range statuses {
		fmt.Fprintf(tw, "%s\t%s\n", s.Target, s.Expiry.Format(time.RFC3339))
	}
	tw.Flush()
}

// SuppressSummary returns a human-readable summary of active suppressions.
func SuppressSummary(statuses []SuppressStatus) string {
	if len(statuses) == 0 {
		return "no targets suppressed"
	}
	return fmt.Sprintf("%d target(s) suppressed", len(statuses))
}
