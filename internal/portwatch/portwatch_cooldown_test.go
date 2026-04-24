package portwatch

import (
	"testing"
	"time"
)

func TestNewScanCooldownManager_InvalidGap(t *testing.T) {
	_, err := NewScanCooldownManager(ScanCooldownConfig{MinGap: 0})
	if err == nil {
		t.Fatal("expected error for zero MinGap")
	}
	_, err = NewScanCooldownManager(ScanCooldownConfig{MinGap: -time.Second})
	if err == nil {
		t.Fatal("expected error for negative MinGap")
	}
}

func TestNewScanCooldownManager_Valid(t *testing.T) {
	m, err := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestReady_BeforeObserve(t *testing.T) {
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Second})
	if !m.Ready("host1") {
		t.Error("expected Ready=true for unseen target")
	}
}

func TestReady_WithinCooldown(t *testing.T) {
	now := time.Now()
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Minute})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	if m.Ready("host1") {
		t.Error("expected Ready=false within cooldown window")
	}
}

func TestReady_AfterCooldown(t *testing.T) {
	now := time.Now()
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Second})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	m.now = func() time.Time { return now.Add(2 * time.Second) }
	if !m.Ready("host1") {
		t.Error("expected Ready=true after cooldown elapsed")
	}
}

func TestReset_AllowsImmediateScan(t *testing.T) {
	now := time.Now()
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Minute})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	m.Reset("host1")
	if !m.Ready("host1") {
		t.Error("expected Ready=true after Reset")
	}
}

func TestNextReady_NoObservation(t *testing.T) {
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Second})
	if !m.NextReady("host1").IsZero() {
		t.Error("expected zero time for unseen target")
	}
}

func TestNextReady_AfterObserve(t *testing.T) {
	now := time.Now()
	gap := 5 * time.Second
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: gap})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	want := now.Add(gap)
	got := m.NextReady("host1")
	if !got.Equal(want) {
		t.Errorf("NextReady: got %v, want %v", got, want)
	}
}

func TestReady_EmptyTarget(t *testing.T) {
	m, _ := NewScanCooldownManager(ScanCooldownConfig{MinGap: time.Second})
	if m.Ready("") {
		t.Error("expected Ready=false for empty target")
	}
}

func TestDefaultScanCooldownConfig_Defaults(t *testing.T) {
	cfg := DefaultScanCooldownConfig()
	if cfg.MinGap <= 0 {
		t.Errorf("expected positive default MinGap, got %v", cfg.MinGap)
	}
}
