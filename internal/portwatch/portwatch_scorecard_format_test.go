package portwatch

import (
	"strings"
	"testing"
)

func TestScorecardSummary_NoTargets(t *testing.T) {
	out := ScorecardSummary(nil)
	if !strings.Contains(out, "no targets") {
		t.Errorf("unexpected summary: %q", out)
	}
}

func TestScorecardSummary_WithTargets(t *testing.T) {
	entries := []ScorecardEntry{
		{Target: "a", TotalScans: 10, TotalAlerts: 2, TotalErrors: 1},
		{Target: "b", TotalScans: 5, TotalAlerts: 0, TotalErrors: 0},
	}
	out := ScorecardSummary(entries)
	for _, want := range []string{"2 target", "15 scan", "2 alert", "1 error"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in summary %q", want, out)
		}
	}
}

func TestAlertRate_NoScans(t *testing.T) {
	e := ScorecardEntry{TotalScans: 0, TotalAlerts: 0}
	if r := AlertRate(e); r != 0 {
		t.Errorf("AlertRate: got %f, want 0", r)
	}
}

func TestAlertRate_Calculated(t *testing.T) {
	e := ScorecardEntry{TotalScans: 4, TotalAlerts: 1}
	if r := AlertRate(e); r != 0.25 {
		t.Errorf("AlertRate: got %f, want 0.25", r)
	}
}

func TestErrorRate_NoScans(t *testing.T) {
	e := ScorecardEntry{TotalScans: 0, TotalErrors: 0}
	if r := ErrorRate(e); r != 0 {
		t.Errorf("ErrorRate: got %f, want 0", r)
	}
}

func TestErrorRate_Calculated(t *testing.T) {
	e := ScorecardEntry{TotalScans: 10, TotalErrors: 3}
	if r := ErrorRate(e); r != 0.3 {
		t.Errorf("ErrorRate: got %f, want 0.3", r)
	}
}
