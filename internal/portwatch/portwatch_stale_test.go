package portwatch

import (
	"testing"
	"time"
)

func TestNewStalenessManager_InvalidMaxAge(t *testing.T) {
	_, err := NewStalenessManager(StalenessConfig{MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxAge")
	}
	_, err = NewStalenessManager(StalenessConfig{MaxAge: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error for negative MaxAge")
	}
}

func TestNewStalenessManager_Valid(t *testing.T) {
	m, err := NewStalenessManager(DefaultStalenessConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestIsStale_NoObservation(t *testing.T) {
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	if !m.IsStale("host1") {
		t.Error("expected stale when no observation recorded")
	}
}

func TestIsStale_RecentObservation(t *testing.T) {
	now := time.Now()
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	if m.IsStale("host1") {
		t.Error("expected not stale after immediate observation")
	}
}

func TestIsStale_ExpiredObservation(t *testing.T) {
	base := time.Now()
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	m.now = func() time.Time { return base }
	m.Observe("host1")
	// advance clock beyond MaxAge
	m.now = func() time.Time { return base.Add(2 * time.Minute) }
	if !m.IsStale("host1") {
		t.Error("expected stale after MaxAge exceeded")
	}
}

func TestObserve_IndependentTargets(t *testing.T) {
	base := time.Now()
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	m.now = func() time.Time { return base }
	m.Observe("host1")
	// host2 never observed
	if m.IsStale("host1") {
		t.Error("host1 should not be stale")
	}
	if !m.IsStale("host2") {
		t.Error("host2 should be stale")
	}
}

func TestLastScan_Missing(t *testing.T) {
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	_, ok := m.LastScan("ghost")
	if ok {
		t.Error("expected ok=false for unknown target")
	}
}

func TestLastScan_Present(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	m.now = func() time.Time { return now }
	m.Observe("host1")
	got, ok := m.LastScan("host1")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if !got.Equal(now) {
		t.Errorf("expected %v, got %v", now, got)
	}
}

func TestTargets_ReturnsObserved(t *testing.T) {
	m, _ := NewStalenessManager(StalenessConfig{MaxAge: time.Minute})
	m.Observe("a")
	m.Observe("b")
	targets := m.Targets()
	if len(targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(targets))
	}
}
