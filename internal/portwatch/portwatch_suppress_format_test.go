package portwatch

import (
	"strings"
	"testing"
	"time"
)

func makeSuppressions() []SuppressStatus {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	return []SuppressStatus{
		{Target: "host-b", Expiry: base.Add(time.Hour)},
		{Target: "host-a", Expiry: base.Add(2 * time.Hour)},
	}
}

func TestWriteSuppressTable_ContainsHeaders(t *testing.T) {
	var b strings.Builder
	WriteSuppressTable(&b, makeSuppressions())
	out := b.String()
	if !strings.Contains(out, "TARGET") {
		t.Error("expected TARGET header")
	}
	if !strings.Contains(out, "SUPPRESSED UNTIL") {
		t.Error("expected SUPPRESSED UNTIL header")
	}
}

func TestWriteSuppressTable_ShowsTarget(t *testing.T) {
	var b strings.Builder
	WriteSuppressTable(&b, makeSuppressions())
	out := b.String()
	if !strings.Contains(out, "host-a") {
		t.Error("expected host-a in output")
	}
	if !strings.Contains(out, "host-b") {
		t.Error("expected host-b in output")
	}
}

func TestWriteSuppressTable_SortedByTarget(t *testing.T) {
	var b strings.Builder
	WriteSuppressTable(&b, makeSuppressions())
	out := b.String()
	idxA := strings.Index(out, "host-a")
	idxB := strings.Index(out, "host-b")
	if idxA > idxB {
		t.Error("expected host-a before host-b")
	}
}

func TestSuppressSummary_NoTargets(t *testing.T) {
	s := SuppressSummary(nil)
	if s != "no targets suppressed" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSuppressSummary_WithTargets(t *testing.T) {
	s := SuppressSummary(makeSuppressions())
	if !strings.Contains(s, "2") {
		t.Errorf("expected count in summary: %s", s)
	}
}
