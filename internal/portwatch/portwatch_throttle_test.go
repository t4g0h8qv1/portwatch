package portwatch

import (
	"testing"
	"time"
)

func TestDefaultThrottleConfig_Defaults(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.MinGap != 5*time.Second {
		t.Fatalf("expected MinGap=5s, got %v", cfg.MinGap)
	}
}

func TestNewScanThrottleManager_InvalidGap(t *testing.T) {
	_, err := NewScanThrottleManager(ThrottleConfig{MinGap: 0})
	if err == nil {
		t.Fatal("expected error for zero MinGap")
	}
}

func TestNewScanThrottleManager_Valid(t *testing.T) {
	m, err := NewScanThrottleManager(ThrottleConfig{MinGap: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: time.Second})
	if !m.Allow("host1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinGap(t *testing.T) {
	now := time.Now()
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: 10 * time.Second})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	m.now = func() time.Time { return now.Add(5 * time.Second) }
	if m.Allow("host1") {
		t.Fatal("expected call within gap to be denied")
	}
}

func TestAllow_PermittedAfterGap(t *testing.T) {
	now := time.Now()
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: 5 * time.Second})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	m.now = func() time.Time { return now.Add(6 * time.Second) }
	if !m.Allow("host1") {
		t.Fatal("expected call after gap to be allowed")
	}
}

func TestAllow_EmptyTarget(t *testing.T) {
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: time.Second})
	if m.Allow("") {
		t.Fatal("expected empty target to be denied")
	}
}

func TestReset_AllowsImmediateScan(t *testing.T) {
	now := time.Now()
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: 10 * time.Second})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	m.now = func() time.Time { return now.Add(2 * time.Second) }
	m.Reset("host1")
	if !m.Allow("host1") {
		t.Fatal("expected allow after reset")
	}
}

func TestLastScan_Missing(t *testing.T) {
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: time.Second})
	_, ok := m.LastScan("ghost")
	if ok {
		t.Fatal("expected no record for unknown target")
	}
}

func TestLastScan_ReturnsObservedTime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	m, _ := NewScanThrottleManager(ThrottleConfig{MinGap: time.Second})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	got, ok := m.LastScan("host1")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, got)
	}
}
