package portwatch

import (
	"testing"
	"time"
)

func TestDefaultJitterConfig_Defaults(t *testing.T) {
	cfg := DefaultJitterConfig()
	if cfg.MaxJitter != 5*time.Second {
		t.Errorf("expected MaxJitter 5s, got %v", cfg.MaxJitter)
	}
}

func TestNewJitterManager_InvalidMaxJitter(t *testing.T) {
	_, err := NewJitterManager(JitterConfig{MaxJitter: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxJitter")
	}
}

func TestNewJitterManager_NegativeMaxJitter(t *testing.T) {
	_, err := NewJitterManager(JitterConfig{MaxJitter: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error for negative MaxJitter")
	}
}

func TestNewJitterManager_Valid(t *testing.T) {
	jm, err := NewJitterManager(JitterConfig{MaxJitter: 2 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jm == nil {
		t.Fatal("expected non-nil JitterManager")
	}
}

func TestDelay_WithinBounds(t *testing.T) {
	max := 100 * time.Millisecond
	jm, err := NewJitterManager(JitterConfig{MaxJitter: max})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 50; i++ {
		d := jm.Delay()
		if d < 0 || d >= max {
			t.Errorf("delay %v out of bounds [0, %v)", d, max)
		}
	}
}

func TestDelay_NonDeterministic(t *testing.T) {
	jm, err := NewJitterManager(JitterConfig{MaxJitter: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seen := map[time.Duration]bool{}
	for i := 0; i < 20; i++ {
		seen[jm.Delay()] = true
	}
	if len(seen) < 2 {
		t.Error("expected multiple distinct delay values")
	}
}

func TestMaxJitter_ReturnsConfigured(t *testing.T) {
	expected := 3 * time.Second
	jm, err := NewJitterManager(JitterConfig{MaxJitter: expected})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jm.MaxJitter() != expected {
		t.Errorf("expected %v, got %v", expected, jm.MaxJitter())
	}
}
