package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteRunnerResult writes a summary table of a RunnerResult to w.
func WriteRunnerResult(w io.Writer, res RunnerResult) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FIELD\tVALUE")
	fmt.Fprintf(tw, "Scans Completed\t%d\n", res.ScansCompleted)
	fmt.Fprintf(tw, "Errors\t%d\n", res.Errors)
	if res.LastScan.IsZero() {
		fmt.Fprintf(tw, "Last Scan\tnever\n")
	} else {
		fmt.Fprintf(tw, "Last Scan\t%s\n", res.LastScan.Format(time.RFC3339))
	}
	tw.Flush()
}

// RunnerSummary returns a one-line summary string for a RunnerResult.
func RunnerSummary(res RunnerResult) string {
	if res.ScansCompleted == 0 {
		return "no scans completed"
	}
	errPart := ""
	if res.Errors > 0 {
		errPart = fmt.Sprintf(", %d error(s)", res.Errors)
	}
	return fmt.Sprintf("%d scan(s) completed%s, last at %s",
		res.ScansCompleted, errPart, res.LastScan.Format(time.RFC3339))
}
