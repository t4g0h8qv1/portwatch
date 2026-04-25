package portwatch

import (
	"testing"
	"time"
)

func TestDefaultScanBudgetConfig_Defaults(t *testing.T) {
	cfg := DefaultScanBudgetConfig()
	if cfg.MaxScansPerDay != 1440 {
		t.Errorf("expected MaxScansPerDay=1440, got %d", cfg.MaxScansPerDay)
	}
	if cfg.Window != 24*time.Hour {
		t.Errorf("expected Window=24h, got %s", cfg.Window)
	}
}

func TestNewScanBudgetManager_InvalidMax(t *testing.T) {
	_, err := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 0, Window: time.Hour})
	if err == nil {
		t.Fatal("expected error for MaxScansPerDay=0")
	}
}

func TestNewScanBudgetManager_InvalidWindow(t *testing.T) {
	_, err := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 10, Window: 0})
	if err == nil {
		t.Fatal("expected error for Window=0")
	}
}

func TestNewScanBudgetManager_Valid(t *testing.T) {
	m, err := NewScanBudgetManager(DefaultScanBudgetConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestAllow_PermitsUpToMax(t *testing.T) {
	m, _ := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 3, Window: time.Hour})
	for i := 0; i < 3; i++ {
		if !m.Allow("host1") {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
	if m.Allow("host1") {
		t.Fatal("expected Allow=false after budget exhausted")
	}
}

func TestAllow_IndependentTargets(t *testing.T) {
	m, _ := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 1, Window: time.Hour})
	if !m.Allow("host1") {
		t.Fatal("expected Allow=true for host1")
	}
	if !m.Allow("host2") {
		t.Fatal("expected Allow=true for host2 (independent budget)")
	}
	if m.Allow("host1") {
		t.Fatal("expected Allow=false for host1 after budget exhausted")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now := time.Now()
	m, _ := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 1, Window: time.Millisecond})
	m.now = func() time.Time { return now }
	if !m.Allow("host1") {
		t.Fatal("expected Allow=true")
	}
	if m.Allow("host1") {
		t.Fatal("expected Allow=false within window")
	}
	m.now = func() time.Time { return now.Add(2 * time.Millisecond) }
	if !m.Allow("host1") {
		t.Fatal("expected Allow=true after window expired")
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	m, _ := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 5, Window: time.Hour})
	if r := m.Remaining("host1"); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	m.Allow("host1")
	m.Allow("host1")
	if r := m.Remaining("host1"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestReset_ClearsBudget(t *testing.T) {
	m, _ := NewScanBudgetManager(ScanBudgetConfig{MaxScansPerDay: 1, Window: time.Hour})
	m.Allow("host1")
	if m.Allow("host1") {
		t.Fatal("expected budget exhausted before reset")
	}
	m.Reset("host1")
	if !m.Allow("host1") {
		t.Fatal("expected Allow=true after reset")
	}
}
