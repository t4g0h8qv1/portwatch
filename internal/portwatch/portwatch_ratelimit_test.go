package portwatch

import (
	"testing"
	"time"
)

func TestNewScanRateLimiter_InvalidWindow(t *testing.T) {
	_, err := NewScanRateLimiter(0, 5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewScanRateLimiter_InvalidMaxScans(t *testing.T) {
	_, err := NewScanRateLimiter(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero maxScans")
	}
}

func TestNewScanRateLimiter_Valid(t *testing.T) {
	rl, err := NewScanRateLimiter(time.Minute, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestAllow_PermitsUpToMax(t *testing.T) {
	rl, _ := NewScanRateLimiter(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !rl.Allow("host1") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}
	if rl.Allow("host1") {
		t.Fatal("expected deny after max scans")
	}
}

func TestAllow_IndependentTargets(t *testing.T) {
	rl, _ := NewScanRateLimiter(time.Minute, 1)
	if !rl.Allow("a") {
		t.Fatal("expected allow for a")
	}
	if !rl.Allow("b") {
		t.Fatal("expected allow for b")
	}
	if rl.Allow("a") {
		t.Fatal("expected deny for a after max")
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	rl, _ := NewScanRateLimiter(time.Minute, 3)
	if r := rl.Remaining("h"); r != 3 {
		t.Fatalf("want 3, got %d", r)
	}
	rl.Allow("h")
	if r := rl.Remaining("h"); r != 2 {
		t.Fatalf("want 2, got %d", r)
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	rl, _ := NewScanRateLimiter(time.Minute, 1)
	rl.Allow("h")
	if rl.Allow("h") {
		t.Fatal("expected deny before reset")
	}
	rl.Reset("h")
	if !rl.Allow("h") {
		t.Fatal("expected allow after reset")
	}
}

func TestAllow_RefillsAfterWindow(t *testing.T) {
	rl, _ := NewScanRateLimiter(50*time.Millisecond, 1)
	rl.Allow("h")
	if rl.Allow("h") {
		t.Fatal("expected deny within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow("h") {
		t.Fatal("expected allow after window expired")
	}
}
