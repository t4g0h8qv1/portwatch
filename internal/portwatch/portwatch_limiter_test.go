package portwatch

import (
	"testing"
	"time"
)

func TestNewScanLimiter_InvalidGap(t *testing.T) {
	_, err := NewScanLimiter(0)
	if err == nil {
		t.Fatal("expected error for zero gap")
	}
}

func TestNewScanLimiter_ValidGap(t *testing.T) {
	l, err := NewScanLimiter(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	l, _ := NewScanLimiter(time.Second)
	if err := l.Allow("host1"); err != nil {
		t.Fatalf("expected first call to be allowed, got: %v", err)
	}
}

func TestAllow_SecondCallWithinGap(t *testing.T) {
	now := time.Now()
	l, _ := NewScanLimiter(time.Minute)
	l.nowFunc = func() time.Time { return now }
	_ = l.Allow("host1")
	if err := l.Allow("host1"); err != ErrTooSoon {
		t.Fatalf("expected ErrTooSoon, got: %v", err)
	}
}

func TestAllow_PermittedAfterGap(t *testing.T) {
	now := time.Now()
	l, _ := NewScanLimiter(time.Minute)
	l.nowFunc = func() time.Time { return now }
	_ = l.Allow("host1")
	l.nowFunc = func() time.Time { return now.Add(2 * time.Minute) }
	if err := l.Allow("host1"); err != nil {
		t.Fatalf("expected allow after gap, got: %v", err)
	}
}

func TestReset_ClearsTarget(t *testing.T) {
	now := time.Now()
	l, _ := NewScanLimiter(time.Minute)
	l.nowFunc = func() time.Time { return now }
	_ = l.Allow("host1")
	l.Reset("host1")
	if err := l.Allow("host1"); err != nil {
		t.Fatalf("expected allow after reset, got: %v", err)
	}
}

func TestLastScan_Missing(t *testing.T) {
	l, _ := NewScanLimiter(time.Second)
	_, ok := l.LastScan("ghost")
	if ok {
		t.Fatal("expected no record for unknown target")
	}
}

func TestLastScan_AfterAllow(t *testing.T) {
	now := time.Now()
	l, _ := NewScanLimiter(time.Second)
	l.nowFunc = func() time.Time { return now }
	_ = l.Allow("host1")
	got, ok := l.LastScan("host1")
	if !ok {
		t.Fatal("expected record after allow")
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, got)
	}
}
