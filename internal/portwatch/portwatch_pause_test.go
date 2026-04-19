package portwatch

import (
	"testing"
	"time"
)

func TestPause_InvalidDuration(t *testing.T) {
	pm := NewPauseManager()
	if err := pm.Pause(-1*time.Second, "test"); err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestPause_IsPaused(t *testing.T) {
	pm := NewPauseManager()
	if err := pm.Pause(1*time.Hour, "maintenance"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pm.IsPaused() {
		t.Fatal("expected paused")
	}
}

func TestResume_ClearsPause(t *testing.T) {
	pm := NewPauseManager()
	_ = pm.Pause(1*time.Hour, "test")
	pm.Resume()
	if pm.IsPaused() {
		t.Fatal("expected not paused after resume")
	}
}

func TestPause_ExpiresAutomatically(t *testing.T) {
	now := time.Now()
	pm := &PauseManager{clock: func() time.Time { return now }}
	_ = pm.Pause(1*time.Second, "short")
	// advance clock past expiry
	pm.clock = func() time.Time { return now.Add(2 * time.Second) }
	if pm.IsPaused() {
		t.Fatal("expected pause to have expired")
	}
}

func TestState_ReflectsPause(t *testing.T) {
	pm := NewPauseManager()
	_ = pm.Pause(30*time.Minute, "deploy window")
	s := pm.State()
	if !s.Paused {
		t.Error("expected Paused=true")
	}
	if s.Reason != "deploy window" {
		t.Errorf("unexpected reason: %q", s.Reason)
	}
	if s.Until.IsZero() {
		t.Error("expected non-zero Until")
	}
}

func TestState_NotPaused(t *testing.T) {
	pm := NewPauseManager()
	s := pm.State()
	if s.Paused {
		t.Error("expected Paused=false")
	}
}

func TestResume_IdempotentWhenNotPaused(t *testing.T) {
	pm := NewPauseManager()
	pm.Resume() // should not panic
	if pm.IsPaused() {
		t.Fatal("expected not paused")
	}
}
