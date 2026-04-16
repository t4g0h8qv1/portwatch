package portping

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteTable writes a human-readable latency table to w.
func WriteTable(w io.Writer, results []Result) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tPORT\tLATENCY\tSTATUS")
	for _, r := range results {
		status := "ok"
		latency := r.Latency.Round(1000).String()
		if r.Err != nil {
			status = "error: " + r.Err.Error()
			latency = "-"
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n", r.Host, r.Port, latency, status)
	}
	return tw.Flush()
}

// Summary returns a one-line summary of ping results.
func Summary(results []Result) string {
	total := len(results)
	if total == 0 {
		return "no ports probed"
	}
	ok := 0
	for _, r := range results {
		if r.Err == nil {
			ok++
		}
	}
	return fmt.Sprintf("%d/%d ports reachable", ok, total)
}
