package portwatch

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// FenceStatus holds display information for a single fenced target.
type FenceStatus struct {
	Target    string
	Reason    string
	ExpiresAt time.Time
}

// WriteFenceTable writes a formatted table of active fences to w.
func WriteFenceTable(w io.Writer, statuses []FenceStatus) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tREASON\tEXPIRES")
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Target < statuses[j].Target
	})
	for _, s := range statuses {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", s.Target, s.Reason, s.ExpiresAt.Format(time.RFC3339))
	}
	tw.Flush()
}

// FenceSummary returns a one-line summary of the fence statuses.
func FenceSummary(statuses []FenceStatus) string {
	if len(statuses) == 0 {
		return "no active fences"
	}
	return fmt.Sprintf("%d target(s) fenced", len(statuses))
}
