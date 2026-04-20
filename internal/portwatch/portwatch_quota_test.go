package portwatch

import (
	"testing"
	"time"
)

func TestNewScanQuotaManager_InvalidMax(t *testing.T) {
	_, err := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 0, Window: time.Hour})
	if err == nil {
		t.Fatal("expected error for zero MaxScansPerHour")
	}
}

func TestNewScanQuotaManager_InvalidWindow(t *testing.T) {
	_, err := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 10, Window: 0})
	if err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestNewScanQuotaManager_Valid(t *testing.T) {
	q, err := NewScanQuotaManager(DefaultScanQuotaConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestAllow_PermitsUpToMax(t *testing.T) {
	q, _ := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 3, Window: time.Hour})
	for i := 0; i < 3; i++ {
		if !q.Allow("host1") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
	if q.Allow("host1") {
		t.Fatal("expected Allow to return false after quota exceeded")
	}
}

func TestAllow_IndependentTargets(t *testing.T) {
	q, _ := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 1, Window: time.Hour})
	q.Allow("host1")
	if !q.Allow("host2") {
		t.Fatal("expected host2 to be allowed independently")
	}
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	q, _ := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 5, Window: time.Hour})
	if r := q.Remaining("host1"); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	q.Allow("host1")
	q.Allow("host1")
	if r := q.Remaining("host1"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now := time.Now()
	q, _ := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 1, Window: time.Millisecond})
	q.now = func() time.Time { return now }
	q.Allow("host1")
	if q.Allow("host1") {
		t.Fatal("expected quota to be exhausted")
	}
	q.now = func() time.Time { return now.Add(2 * time.Millisecond) }
	if !q.Allow("host1") {
		t.Fatal("expected quota to reset after window")
	}
}

func TestReset_ClearsTarget(t *testing.T) {
	q, _ := NewScanQuotaManager(ScanQuotaConfig{MaxScansPerHour: 1, Window: time.Hour})
	q.Allow("host1")
	if q.Allow("host1") {
		t.Fatal("expected quota exhausted")
	}
	q.Reset("host1")
	if !q.Allow("host1") {
		t.Fatal("expected quota reset after Reset()")
	}
}
