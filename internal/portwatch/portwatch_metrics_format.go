package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// WriteMetrics writes a human-readable metrics table to w.
func WriteMetrics(w io.Writer, m Metrics) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tVALUE")
	fmt.Fprintf(tw, "scans_total\t%d\n", m.ScansTotal)
	fmt.Fprintf(tw, "alerts_total\t%d\n", m.AlertsTotal)
	fmt.Fprintf(tw, "errors_total\t%d\n", m.ErrorsTotal)
	fmt.Fprintf(tw, "open_port_count\t%d\n", m.OpenPortCount)
	fmt.Fprintf(tw, "last_scan_at\t%s\n", formatTime(m.LastScanAt))
	fmt.Fprintf(tw, "last_alert_at\t%s\n", formatTime(m.LastAlertAt))
	tw.Flush()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}
