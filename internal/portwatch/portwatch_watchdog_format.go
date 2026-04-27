package portwatch

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// WatchdogStatus holds display data for a single target.
type WatchdogStatus struct {
	Target   string
	LastSeen time.Time
	Expired  bool
}

// WriteWatchdogTable writes a formatted table of watchdog statuses to w.
func WriteWatchdogTable(w io.Writer, statuses []WatchdogStatus) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tLAST SEEN\tSTATUS")
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Target < statuses[j].Target
	})
	for _, s := range statuses {
		var lastSeen string
		if s.LastSeen.IsZero() {
			lastSeen = "never"
		} else {
			lastSeen = s.LastSeen.Format(time.RFC3339)
		}
		status := "ok"
		if s.Expired {
			status = "expired"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", s.Target, lastSeen, status)
	}
	tw.Flush()
}

// WatchdogSummary returns a human-readable summary of watchdog statuses.
func WatchdogSummary(statuses []WatchdogStatus) string {
	if len(statuses) == 0 {
		return "watchdog: no targets registered"
	}
	expired := 0
	for _, s := range statuses {
		if s.Expired {
			expired++
		}
	}
	if expired == 0 {
		return fmt.Sprintf("watchdog: all %d target(s) active", len(statuses))
	}
	return fmt.Sprintf("watchdog: %d/%d target(s) expired", expired, len(statuses))
}
