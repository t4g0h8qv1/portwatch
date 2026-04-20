package portwatch

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// WindowStatus holds display info for a single target window.
type WindowStatus struct {
	Target  string
	Start   time.Duration
	End     time.Duration
	Allowed bool
}

// WriteWindowTable writes a tabular summary of all registered windows.
func WriteWindowTable(w io.Writer, m *ScanWindowManager) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tWINDOW START\tWINDOW END\tALLOWED")
	targets := m.Targets()
	sort.Strings(targets)
	for _, t := range targets {
		m.mu.RLock()
		cfg := m.windows[t]
		m.mu.RUnlock()
		allowed := m.Allowed(t)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", t, fmtOffset(cfg.Start), fmtOffset(cfg.End), allowed)
	}
	tw.Flush()
}

// WindowSummary returns a one-line summary of window states.
func WindowSummary(m *ScanWindowManager) string {
	targets := m.Targets()
	if len(targets) == 0 {
		return "no scan windows configured"
	}
	allowed := 0
	for _, t := range targets {
		if m.Allowed(t) {
			allowed++
		}
	}
	return fmt.Sprintf("%d/%d targets within scan window", allowed, len(targets))
}

func fmtOffset(d time.Duration) string {
	h := int(d.Hours())
	min := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", h, min)
}
