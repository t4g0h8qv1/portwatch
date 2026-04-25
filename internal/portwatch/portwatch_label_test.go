package portwatch

import (
	"testing"
)

func TestNewLabelManager_Empty(t *testing.T) {
	m := NewLabelManager()
	if got := m.Get("host1"); got != "" {
		t.Fatalf("expected empty label, got %q", got)
	}
}

func TestSet_And_Get(t *testing.T) {
	m := NewLabelManager()
	if err := m.Set("192.168.1.1", "router"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Get("192.168.1.1"); got != "router" {
		t.Fatalf("expected %q, got %q", "router", got)
	}
}

func TestSet_EmptyTarget_ReturnsError(t *testing.T) {
	m := NewLabelManager()
	if err := m.Set("", "label"); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestSet_EmptyLabel_ReturnsError(t *testing.T) {
	m := NewLabelManager()
	if err := m.Set("host1", ""); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestLabel_FallsBackToTarget(t *testing.T) {
	m := NewLabelManager()
	if got := m.Label("10.0.0.1"); got != "10.0.0.1" {
		t.Fatalf("expected fallback to target, got %q", got)
	}
}

func TestLabel_ReturnsSetLabel(t *testing.T) {
	m := NewLabelManager()
	_ = m.Set("10.0.0.1", "gateway")
	if got := m.Label("10.0.0.1"); got != "gateway" {
		t.Fatalf("expected %q, got %q", "gateway", got)
	}
}

func TestRemove_ClearsLabel(t *testing.T) {
	m := NewLabelManager()
	_ = m.Set("host1", "web")
	m.Remove("host1")
	if got := m.Get("host1"); got != "" {
		t.Fatalf("expected empty after remove, got %q", got)
	}
}

func TestAll_SortedByTarget(t *testing.T) {
	m := NewLabelManager()
	_ = m.Set("z-host", "last")
	_ = m.Set("a-host", "first")
	_ = m.Set("m-host", "middle")
	pairs := m.All()
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	if pairs[0][0] != "a-host" || pairs[1][0] != "m-host" || pairs[2][0] != "z-host" {
		t.Fatalf("unexpected order: %v", pairs)
	}
}

func TestAll_EmptyManager(t *testing.T) {
	m := NewLabelManager()
	if got := m.All(); len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}
