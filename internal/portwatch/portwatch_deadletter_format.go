package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteDeadLetterTable writes a human-readable table of dead-letter entries to w.
func WriteDeadLetterTable(w io.Writer, entries []DeadLetterEntry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tATTEMPTS\tOCCURRED AT\tERROR")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n",
			e.Target,
			e.Attempts,
			e.OccurredAt.Format(time.RFC3339),
			e.Err.Error(),
		)
	}
	tw.Flush()
}

// DeadLetterSummary returns a one-line summary of the queue state.
func DeadLetterSummary(entries []DeadLetterEntry) string {
	if len(entries) == 0 {
		return "dead-letter queue: empty"
	}
	return fmt.Sprintf("dead-letter queue: %d unprocessed event(s)", len(entries))
}
