package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteSignalTable writes a human-readable table of scan signals to w.
func WriteSignalTable(w io.Writer, signals []ScanSignal) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tPORT\tKIND\tOBSERVED AT")
	fmt.Fprintln(tw, "------\t----\t----\t-----------")
	for _, s := range signals {
		ts := "never"
		if !s.ObservedAt.IsZero() {
			ts = s.ObservedAt.UTC().Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n", s.Target, s.Port, s.Kind, ts)
	}
	tw.Flush()
}

// SignalSummary returns a one-line summary of the provided signals.
func SignalSummary(signals []ScanSignal) string {
	if len(signals) == 0 {
		return "no signals recorded"
	}
	counts := map[SignalKind]int{}
	for _, s := range signals {
		counts[s.Kind]++
	}
	return fmt.Sprintf("%d signal(s): %d opened, %d closed, %d stable",
		len(signals),
		counts[SignalOpened],
		counts[SignalClosed],
		counts[SignalStable],
	)
}
