package portwatch

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// WriteStalenessTable writes a human-readable table of staleness status for
// all observed targets to w.
func WriteStalenessTable(w io.Writer, m *StalenessManager) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tLAST SCAN\tSTALE")
	targets := m.Targets()
	sort.Strings(targets)
	for _, target := range targets {
		t, _ := m.LastScan(target)
		stale := m.IsStale(target)
		fmt.Fprintf(tw, "%s\t%s\t%v\n", target, formatStalenessTime(t), stale)
	}
	tw.Flush()
}

// StalenessSummary returns a one-line summary of stale vs. healthy targets.
func StalenessSummary(m *StalenessManager) string {
	targets := m.Targets()
	if len(targets) == 0 {
		return "no targets observed"
	}
	staleCount := 0
	for _, target := range targets {
		if m.IsStale(target) {
			staleCount++
		}
	}
	healthy := len(targets) - staleCount
	if staleCount == 0 {
		return fmt.Sprintf("all %d target(s) up-to-date", healthy)
	}
	return fmt.Sprintf("%d/%d target(s) stale", staleCount, len(targets))
}

func formatStalenessTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}
