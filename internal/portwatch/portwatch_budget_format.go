package portwatch

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// BudgetStatus holds a snapshot of budget usage for a single target.
type BudgetStatus struct {
	Target    string
	Used      int
	Remaining int
	Max       int
}

// WriteBudgetTable writes a formatted table of budget statuses to w.
func WriteBudgetTable(w io.Writer, statuses []BudgetStatus) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tUSED\tREMAINING\tMAX")
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Target < statuses[j].Target
	})
	for _, s := range statuses {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\n", s.Target, s.Used, s.Remaining, s.Max)
	}
	tw.Flush()
}

// BudgetSummary returns a human-readable summary of budget statuses.
func BudgetSummary(statuses []BudgetStatus) string {
	if len(statuses) == 0 {
		return "no budget entries"
	}
	exhausted := 0
	for _, s := range statuses {
		if s.Remaining == 0 {
			exhausted++
		}
	}
	if exhausted == 0 {
		return fmt.Sprintf("%d target(s) within budget", len(statuses))
	}
	return fmt.Sprintf("%d/%d target(s) budget exhausted", exhausted, len(statuses))
}
