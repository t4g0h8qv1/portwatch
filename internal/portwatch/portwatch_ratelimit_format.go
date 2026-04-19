package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// RateLimitStatus holds a snapshot of rate limit state for a target.
type RateLimitStatus struct {
	Target    string
	Remaining int
	Max       int
}

// WriteRateLimitTable writes a formatted table of rate limit statuses to w.
func WriteRateLimitTable(w io.Writer, statuses []RateLimitStatus) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tREMAINING\tMAX")
	for _, s := range statuses {
		fmt.Fprintf(tw, "%s\t%d\t%d\n", s.Target, s.Remaining, s.Max)
	}
	tw.Flush()
}

// RateLimitSummary returns a one-line summary of rate limit statuses.
func RateLimitSummary(statuses []RateLimitStatus) string {
	if len(statuses) == 0 {
		return "no targets tracked"
	}
	throttled := 0
	for _, s := range statuses {
		if s.Remaining == 0 {
			throttled++
		}
	}
	if throttled == 0 {
		return fmt.Sprintf("%d target(s) within limit", len(statuses))
	}
	return fmt.Sprintf("%d/%d target(s) throttled", throttled, len(statuses))
}
