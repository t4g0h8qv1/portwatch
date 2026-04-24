package portwatch

import (
	"testing"
	"time"
)

func TestDefaultDebounceConfig_Defaults(t *testing.T) {
	cfg := DefaultDebounceConfig()
	if cfg.Window <= 0 {
		t.Fatalf("expected positive window, got %v", cfg.Window)
	}
}

func TestNewDebounceManager_InvalidWindow(t *testing.T) {
	_, err := NewDebounceManager(DebounceConfig{Window: 0})
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewDebounceManager_Valid(t *testing.T) {
	_, err := NewDebounceManager(DebounceConfig{Window: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReady_BeforeObserve(t *testing.T) {
	dm, _ := NewDebounceManager(DebounceConfig{Window: time.Second})
	if !dm.Ready("host1") {
		t.Fatal("expected ready before any observation")
	}
}

func TestReady_WithinWindow(t *testing.T) {
	now := time.Now()
	dm, _ := NewDebounceManager(DebounceConfig{Window: time.Minute})
	dm.now = func() time.Time { return now }
	dm.Observe("host1")
	if dm.Ready("host1") {
		t.Fatal("expected not ready within debounce window")
	}
}

func TestReady_AfterWindow(t *testing.T) {
	base := time.Now()
	dm, _ := NewDebounceManager(DebounceConfig{Window: time.Second})
	dm.now = func() time.Time { return base }
	dm.Observe("host1")
	dm.now = func() time.Time { return base.Add(2 * time.Second) }
	if !dm.Ready("host1") {
		t.Fatal("expected ready after window elapsed")
	}
}

func TestObserve_IndependentTargets(t *testing.T) {
	now := time.Now()
	dm, _ := NewDebounceManager(DebounceConfig{Window: time.Minute})
	dm.now = func() time.Time { return now }
	dm.Observe("host1")
	if !dm.Ready("host2") {
		t.Fatal("host2 should be unaffected by host1 observation")
	}
	if dm.Ready("host1") {
		t.Fatal("host1 should not be ready within window")
	}
}

func TestReset_ClearsDebounce(t *testing.T) {
	now := time.Now()
	dm, _ := NewDebounceManager(DebounceConfig{Window: time.Minute})
	dm.now = func() time.Time { return now }
	dm.Observe("host1")
	dm.Reset("host1")
	if !dm.Ready("host1") {
		t.Fatal("expected ready after reset")
	}
}
