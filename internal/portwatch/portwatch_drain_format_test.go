package portwatch

import (
	"strings"
	"testing"
)

func TestWriteDrainStatus_ContainsHeaders(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())
	var sb strings.Builder
	WriteDrainStatus(&sb, dm)
	out := sb.String()
	for _, want := range []string{"METRIC", "VALUE", "in_flight"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestWriteDrainStatus_ShowsCount(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())
	dm.Acquire("h1")
	dm.Acquire("h2")
	var sb strings.Builder
	WriteDrainStatus(&sb, dm)
	if !strings.Contains(sb.String(), "2") {
		t.Errorf("expected count 2 in output: %s", sb.String())
	}
}

func TestDrainSummary_Idle(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())
	s := DrainSummary(dm)
	if !strings.Contains(s, "idle") {
		t.Errorf("expected idle summary, got: %s", s)
	}
}

func TestDrainSummary_Active(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())
	dm.Acquire("host-x")
	s := DrainSummary(dm)
	if !strings.Contains(s, "1") || !strings.Contains(s, "in flight") {
		t.Errorf("unexpected summary: %s", s)
	}
}
