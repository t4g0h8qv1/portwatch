package portcooldown

import (
	"errors"
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if !errors.Is(err, ErrInvalidWindow) {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestNew_ValidWindow(t *testing.T) {
	tr, err := New(5 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestReady_BeforeWindow(t *testing.T) {
	tr, _ := New(10 * time.Second)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Observe("host1", 80)
	if tr.Ready("host1", 80) {
		t.Fatal("expected not ready before window elapses")
	}
}

func TestReady_AfterWindow(t *testing.T) {
	tr, _ := New(10 * time.Second)
	base := time.Now()
	tr.now = func() time.Time { return base }
	tr.Observe("host1", 443)
	tr.now = func() time.Time { return base.Add(15 * time.Second) }
	if !tr.Ready("host1", 443) {
		t.Fatal("expected ready after window elapses")
	}
}

func TestObserve_IdempotentClock(t *testing.T) {
	tr, _ := New(10 * time.Second)
	base := time.Now()
	tr.now = func() time.Time { return base }
	tr.Observe("host1", 22)
	tr.now = func() time.Time { return base.Add(5 * time.Second) }
	tr.Observe("host1", 22) // should not reset clock
	tr.now = func() time.Time { return base.Add(12 * time.Second) }
	if !tr.Ready("host1", 22) {
		t.Fatal("expected ready; second Observe must not reset first-seen time")
	}
}

func TestForget_ResetsPort(t *testing.T) {
	tr, _ := New(5 * time.Second)
	base := time.Now()
	tr.now = func() time.Time { return base }
	tr.Observe("host1", 8080)
	tr.Forget("host1", 8080)
	tr.now = func() time.Time { return base.Add(10 * time.Second) }
	if tr.Ready("host1", 8080) {
		t.Fatal("expected not ready after Forget")
	}
}

func TestReady_UnknownHost(t *testing.T) {
	tr, _ := New(5 * time.Second)
	if tr.Ready("ghost", 80) {
		t.Fatal("expected false for unknown host")
	}
}
