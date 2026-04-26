package portwatch

import (
	"testing"
	"time"
)

func TestNewRosterManager_Empty(t *testing.T) {
	rm := NewRosterManager()
	if got := rm.All(); len(got) != 0 {
		t.Fatalf("expected empty roster, got %d entries", len(got))
	}
}

func TestRegister_AddsTarget(t *testing.T) {
	rm := NewRosterManager()
	if err := rm.Register("host1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := rm.Get("host1")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if !e.Active {
		t.Error("expected entry to be active")
	}
}

func TestRegister_EmptyTarget_ReturnsError(t *testing.T) {
	rm := NewRosterManager()
	if err := rm.Register(""); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestRegister_Idempotent(t *testing.T) {
	rm := NewRosterManager()
	_ = rm.Register("host1")
	_ = rm.Register("host1")
	if got := rm.All(); len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}

func TestTouch_UpdatesLastSeen(t *testing.T) {
	rm := NewRosterManager()
	_ = rm.Register("host1")
	now := time.Now()
	if err := rm.Touch("host1", now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := rm.Get("host1")
	if !e.LastSeen.Equal(now) {
		t.Errorf("expected LastSeen %v, got %v", now, e.LastSeen)
	}
}

func TestTouch_UnregisteredTarget_ReturnsError(t *testing.T) {
	rm := NewRosterManager()
	if err := rm.Touch("ghost", time.Now()); err == nil {
		t.Fatal("expected error for unregistered target")
	}
}

func TestDeactivate_MarksInactive(t *testing.T) {
	rm := NewRosterManager()
	_ = rm.Register("host1")
	if err := rm.Deactivate("host1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := rm.Get("host1")
	if e.Active {
		t.Error("expected entry to be inactive")
	}
}

func TestDeactivate_UnregisteredTarget_ReturnsError(t *testing.T) {
	rm := NewRosterManager()
	if err := rm.Deactivate("ghost"); err == nil {
		t.Fatal("expected error for unregistered target")
	}
}

func TestAll_SortedByTarget(t *testing.T) {
	rm := NewRosterManager()
	for _, h := range []string{"charlie", "alpha", "bravo"} {
		_ = rm.Register(h)
	}
	all := rm.All()
	expected := []string{"alpha", "bravo", "charlie"}
	for i, e := range all {
		if e.Target != expected[i] {
			t.Errorf("pos %d: want %q, got %q", i, expected[i], e.Target)
		}
	}
}
