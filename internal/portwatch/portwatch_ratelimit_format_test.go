package portwatch

import (
	"strings"
	"testing"
)

func TestWriteRateLimitTable_ContainsHeaders(t *testing.T) {
	var b strings.Builder
	WriteRateLimitTable(&b, nil)
	out := b.String()
	for _, h := range []string{"TARGET", "REMAINING", "MAX"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestWriteRateLimitTable_ShowsTarget(t *testing.T) {
	var b strings.Builder
	WriteRateLimitTable(&b, []RateLimitStatus{
		{Target: "192.168.1.1", Remaining: 2, Max: 5},
	})
	out := b.String()
	if !strings.Contains(out, "192.168.1.1") {
		t.Error("expected target in output")
	}
	if !strings.Contains(out, "2") {
		t.Error("expected remaining count in output")
	}
}

func TestRateLimitSummary_NoTargets(t *testing.T) {
	s := RateLimitSummary(nil)
	if !strings.Contains(s, "no targets") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestRateLimitSummary_AllWithinLimit(t *testing.T) {
	statuses := []RateLimitStatus{
		{Target: "a", Remaining: 3, Max: 5},
		{Target: "b", Remaining: 1, Max: 5},
	}
	s := RateLimitSummary(statuses)
	if !strings.Contains(s, "within limit") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestRateLimitSummary_SomeThrottled(t *testing.T) {
	statuses := []RateLimitStatus{
		{Target: "a", Remaining: 0, Max: 5},
		{Target: "b", Remaining: 2, Max: 5},
	}
	s := RateLimitSummary(statuses)
	if !strings.Contains(s, "throttled") {
		t.Errorf("expected throttled in summary: %q", s)
	}
	if !strings.Contains(s, "1/2") {
		t.Errorf("expected count in summary: %q", s)
	}
}
