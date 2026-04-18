package portpolicy

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteTable writes a formatted table of policy results to w.
func WriteTable(w io.Writer, results []Result) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tSTATUS\tREASON")
	fmt.Fprintln(tw, "----\t------\t------")
	for _, r := range results {
		fmt.Fprintf(tw, "%d\t%s\t%s\n", r.Port, r.Status, r.Reason)
	}
	tw.Flush()
}

// Summary returns a one-line summary of policy evaluation results.
func Summary(results []Result) string {
	var allowed, denied, unreviewed int
	for _, r := range results {
		switch r.Status {
		case Allowed:
			allowed++
		case Denied:
			denied++
		case Unreviewed:
			unreviewed++
		}
	}
	return fmt.Sprintf("%d allowed, %d denied, %d unreviewed", allowed, denied, unreviewed)
}
