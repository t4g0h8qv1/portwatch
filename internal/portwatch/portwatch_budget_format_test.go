package portwatch

import (
	"strings"
	"testing"
)

func makeBudgetStatuses() []BudgetStatus {
	return []BudgetStatus{
		{Target: "host-b", Used: 10, Remaining: 90, Max: 100},
		{Target: "host-a", Used: 100, Remaining: 0, Max: 100},
	}
}

func TestWriteBudgetTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteBudgetTable(&sb, makeBudgetStatuses())
	out := sb.String()
	for _, hdr := range []string{"TARGET", "USED", "REMAINING", "MAX"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWriteBudgetTable_ShowsTarget(t *testing.T) {
	var sb strings.Builder
	WriteBudgetTable(&sb, makeBudgetStatuses())
	out := sb.String()
	if !strings.Contains(out, "host-a") {
		t.Error("expected host-a in output")
	}
	if !strings.Contains(out, "host-b") {
		t.Error("expected host-b in output")
	}
}

func TestWriteBudgetTable_SortedByTarget(t *testing.T) {
	var sb strings.Builder
	WriteBudgetTable(&sb, makeBudgetStatuses())
	out := sb.String()
	idxA := strings.Index(out, "host-a")
	idxB := strings.Index(out, "host-b")
	if idxA > idxB {
		t.Error("expected host-a before host-b in sorted output")
	}
}

func TestBudgetSummary_NoTargets(t *testing.T) {
	s := BudgetSummary(nil)
	if s != "no budget entries" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestBudgetSummary_AllWithinBudget(t *testing.T) {
	statuses := []BudgetStatus{
		{Target: "host-a", Used: 5, Remaining: 95, Max: 100},
	}
	s := BudgetSummary(statuses)
	if !strings.Contains(s, "within budget") {
		t.Errorf("expected 'within budget' in summary, got %q", s)
	}
}

func TestBudgetSummary_SomeExhausted(t *testing.T) {
	statuses := []BudgetStatus{
		{Target: "host-a", Used: 100, Remaining: 0, Max: 100},
		{Target: "host-b", Used: 50, Remaining: 50, Max: 100},
	}
	s := BudgetSummary(statuses)
	if !strings.Contains(s, "exhausted") {
		t.Errorf("expected 'exhausted' in summary, got %q", s)
	}
	if !strings.Contains(s, "1/2") {
		t.Errorf("expected '1/2' in summary, got %q", s)
	}
}
