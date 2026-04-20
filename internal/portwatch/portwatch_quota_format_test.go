package portwatch

import (
	"strings"
	"testing"
)

func makeQuotaStatuses() []QuotaStatus {
	return []QuotaStatus{
		{Target: "host1", Remaining: 5, Max: 10, Throttled: false},
		{Target: "host2", Remaining: 0, Max: 10, Throttled: true},
	}
}

func TestWriteQuotaTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteQuotaTable(&sb, makeQuotaStatuses())
	out := sb.String()
	for _, h := range []string{"TARGET", "REMAINING", "MAX", "STATUS"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestWriteQuotaTable_ShowsTarget(t *testing.T) {
	var sb strings.Builder
	WriteQuotaTable(&sb, makeQuotaStatuses())
	out := sb.String()
	if !strings.Contains(out, "host1") {
		t.Error("expected host1 in output")
	}
	if !strings.Contains(out, "throttled") {
		t.Error("expected throttled status in output")
	}
}

func TestQuotaSummary_NoTargets(t *testing.T) {
	s := QuotaSummary(nil)
	if !strings.Contains(s, "no quota") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestQuotaSummary_AllWithinLimit(t *testing.T) {
	statuses := []QuotaStatus{
		{Target: "host1", Remaining: 3, Max: 10, Throttled: false},
	}
	s := QuotaSummary(statuses)
	if !strings.Contains(s, "within quota") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestQuotaSummary_SomeThrottled(t *testing.T) {
	s := QuotaSummary(makeQuotaStatuses())
	if !strings.Contains(s, "throttled") {
		t.Errorf("unexpected summary: %q", s)
	}
}
