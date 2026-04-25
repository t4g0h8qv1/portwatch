package portwatch

import (
	"testing"
	"time"
)

func TestDefaultDecayConfig_Defaults(t *testing.T) {
	cfg := DefaultDecayConfig()
	if cfg.HalfLife <= 0 {
		t.Fatal("expected positive HalfLife")
	}
	if cfg.InitialScore <= 0 {
		t.Fatal("expected positive InitialScore")
	}
}

func TestNewDecayManager_InvalidHalfLife(t *testing.T) {
	_, err := NewDecayManager(DecayConfig{HalfLife: 0, InitialScore: 100})
	if err == nil {
		t.Fatal("expected error for zero HalfLife")
	}
}

func TestNewDecayManager_InvalidInitialScore(t *testing.T) {
	_, err := NewDecayManager(DecayConfig{HalfLife: time.Minute, InitialScore: 0})
	if err == nil {
		t.Fatal("expected error for zero InitialScore")
	}
}

func TestNewDecayManager_Valid(t *testing.T) {
	m, err := NewDecayManager(DefaultDecayConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestScore_UnobservedTargetIsZero(t *testing.T) {
	m, _ := NewDecayManager(DefaultDecayConfig())
	if got := m.Score("host:22"); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestObserve_SetsInitialScore(t *testing.T) {
	cfg := DefaultDecayConfig()
	m, _ := NewDecayManager(cfg)
	now := time.Now()
	m.now = func() time.Time { return now }
	m.Observe("host:80")
	got := m.Score("host:80")
	if got != cfg.InitialScore {
		t.Fatalf("expected %.2f, got %.2f", cfg.InitialScore, got)
	}
}

func TestScore_DecaysOverTime(t *testing.T) {
	cfg := DecayConfig{HalfLife: time.Hour, InitialScore: 100}
	m, _ := NewDecayManager(cfg)
	base := time.Now()
	m.now = func() time.Time { return base }
	m.Observe("host:443")

	// Advance time by one half-life; score should be ~50.
	m.now = func() time.Time { return base.Add(time.Hour) }
	got := m.Score("host:443")
	const want = 50.0
	const tolerance = 0.5
	if got < want-tolerance || got > want+tolerance {
		t.Fatalf("expected score near %.1f after one half-life, got %.4f", want, got)
	}
}

func TestReset_RemovesTarget(t *testing.T) {
	m, _ := NewDecayManager(DefaultDecayConfig())
	m.Observe("host:22")
	m.Reset("host:22")
	if got := m.Score("host:22"); got != 0 {
		t.Fatalf("expected 0 after reset, got %v", got)
	}
}

func TestTargets_ReturnsTracked(t *testing.T) {
	m, _ := NewDecayManager(DefaultDecayConfig())
	m.Observe("host:22")
	m.Observe("host:80")
	targets := m.Targets()
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
}

func TestObserve_ResetScoreOnReobserve(t *testing.T) {
	cfg := DecayConfig{HalfLife: time.Hour, InitialScore: 100}
	m, _ := NewDecayManager(cfg)
	base := time.Now()
	m.now = func() time.Time { return base }
	m.Observe("host:22")

	// Advance and let score decay.
	m.now = func() time.Time { return base.Add(2 * time.Hour) }
	_ = m.Score("host:22")

	// Re-observe should reset to initial.
	m.Observe("host:22")
	got := m.Score("host:22")
	if got != cfg.InitialScore {
		t.Fatalf("expected %.2f after re-observe, got %.4f", cfg.InitialScore, got)
	}
}
