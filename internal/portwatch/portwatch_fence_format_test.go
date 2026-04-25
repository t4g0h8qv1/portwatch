package portwatch

import (
	"strings"
	"testing"
	"time"
)

func makeFenceStatuses() []FenceStatus {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return []FenceStatus{
		{Target: "host-b", Reason: "maintenance", ExpiresAt: base.Add(10 * time.Minute)},
		{Target: "host-a", Reason: "deploy", ExpiresAt: base.Add(5 * time.Minute)},
	}
}

func TestWriteFenceTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteFenceTable(&sb, makeFenceStatuses())
	out := sb.String()
	for _, h := range []string{"TARGET", "REASON", "EXPIRES"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestWriteFenceTable_ShowsTarget(t *testing.T) {
	var sb strings.Builder
	WriteFenceTable(&sb, makeFenceStatuses())
	out := sb.String()
	if !strings.Contains(out, "host-a") {
		t.Error("expected host-a in output")
	}
	if !strings.Contains(out, "host-b") {
		t.Error("expected host-b in output")
	}
}

func TestWriteFenceTable_SortedByTarget(t *testing.T) {
	var sb strings.Builder
	WriteFenceTable(&sb, makeFenceStatuses())
	out := sb.String()
	idxA := strings.Index(out, "host-a")
	idxB := strings.Index(out, "host-b")
	if idxA > idxB {
		t.Error("expected host-a before host-b in sorted output")
	}
}

func TestFenceSummary_NoTargets(t *testing.T) {
	s := FenceSummary(nil)
	if s != "no active fences" {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestFenceSummary_WithTargets(t *testing.T) {
	s := FenceSummary(makeFenceStatuses())
	if !strings.Contains(s, "2") {
		t.Fatalf("expected count in summary, got %q", s)
	}
}
