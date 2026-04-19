package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteDrainStatus writes a human-readable summary of the DrainManager state
// to w.
func WriteDrainStatus(w io.Writer, dm *DrainManager) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tVALUE")
	fmt.Fprintf(tw, "in_flight\t%d\n", dm.InFlight())
	tw.Flush()
}

// DrainSummary returns a one-line string describing the drain state.
func DrainSummary(dm *DrainManager) string {
	n := dm.InFlight()
	if n == 0 {
		return "drain: idle (no in-flight scans)"
	}
	return fmt.Sprintf("drain: %d scan(s) in flight", n)
}
