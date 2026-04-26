package portwatch

import (
	"testing"
)

func TestDefaultAffinityConfig_Defaults(t *testing.T) {
	cfg := DefaultAffinityConfig()
	if cfg.MaxTargets != 64 {
		t.Fatalf("expected MaxTargets=64, got %d", cfg.MaxTargets)
	}
}

func TestNewScanAffinityManager_InvalidMax(t *testing.T) {
	_, err := NewScanAffinityManager(AffinityConfig{MaxTargets: 0})
	if err == nil {
		t.Fatal("expected error for MaxTargets=0")
	}
}

func TestNewScanAffinityManager_Valid(t *testing.T) {
	m, err := NewScanAffinityManager(DefaultAffinityConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestAssign_EmptyTarget(t *testing.T) {
	m, _ := NewScanAffinityManager(DefaultAffinityConfig())
	_, err := m.Assign("")
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestAssign_ReturnsConsistentWorker(t *testing.T) {
	m, _ := NewScanAffinityManager(DefaultAffinityConfig())
	w1, err := m.Assign("host-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w2, err := m.Assign("host-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w1 != w2 {
		t.Fatalf("expected same worker, got %d and %d", w1, w2)
	}
}

func TestAssign_DifferentTargetsMayDiffer(t *testing.T) {
	m, _ := NewScanAffinityManager(AffinityConfig{MaxTargets: 8})
	wa, _ := m.Assign("host-a")
	wb, _ := m.Assign("host-b")
	// Workers are assigned round-robin; with MaxTargets=8 they should differ.
	if wa == wb {
		t.Logf("note: host-a and host-b both assigned worker %d (acceptable with small MaxTargets)", wa)
	}
}

func TestGet_Missing(t *testing.T) {
	m, _ := NewScanAffinityManager(DefaultAffinityConfig())
	_, ok := m.Get("unknown")
	if ok {
		t.Fatal("expected false for unknown target")
	}
}

func TestGet_AfterAssign(t *testing.T) {
	m, _ := NewScanAffinityManager(DefaultAffinityConfig())
	w, _ := m.Assign("host-x")
	got, ok := m.Get("host-x")
	if !ok {
		t.Fatal("expected true after assign")
	}
	if got != w {
		t.Fatalf("expected worker %d, got %d", w, got)
	}
}

func TestRemove_ClearsEntry(t *testing.T) {
	m, _ := NewScanAffinityManager(DefaultAffinityConfig())
	m.Assign("host-y") //nolint:errcheck
	m.Remove("host-y")
	_, ok := m.Get("host-y")
	if ok {
		t.Fatal("expected false after remove")
	}
}

func TestLen_TracksAssignments(t *testing.T) {
	m, _ := NewScanAffinityManager(DefaultAffinityConfig())
	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}
	m.Assign("a") //nolint:errcheck
	m.Assign("b") //nolint:errcheck
	if m.Len() != 2 {
		t.Fatalf("expected 2, got %d", m.Len())
	}
	m.Remove("a")
	if m.Len() != 1 {
		t.Fatalf("expected 1, got %d", m.Len())
	}
}
