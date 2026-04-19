package portwatch

import (
	"testing"
	"time"
)

func TestDefaultBackoffConfig_Defaults(t *testing.T) {
	cfg := DefaultBackoffConfig()
	if cfg.InitialInterval != 5*time.Second {
		t.Errorf("expected 5s initial interval, got %v", cfg.InitialInterval)
	}
	if cfg.MaxInterval != 5*time.Minute {
		t.Errorf("expected 5m max interval, got %v", cfg.MaxInterval)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected multiplier 2.0, got %v", cfg.Multiplier)
	}
}

func TestRecordFailure_ReturnsIncreasingIntervals(t *testing.T) {
	cfg := BackoffConfig{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
	}
	bm := NewBackoffManager(cfg)

	d1 := bm.RecordFailure("host1")
	d2 := bm.RecordFailure("host1")
	d3 := bm.RecordFailure("host1")

	if d1 != 1*time.Second {
		t.Errorf("expected 1s, got %v", d1)
	}
	if d2 != 2*time.Second {
		t.Errorf("expected 2s, got %v", d2)
	}
	if d3 != 4*time.Second {
		t.Errorf("expected 4s, got %v", d3)
	}
}

func TestRecordFailure_CapsAtMax(t *testing.T) {
	cfg := BackoffConfig{
		InitialInterval: 1 * time.Second,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
	}
	bm := NewBackoffManager(cfg)
	var last time.Duration
	for i := 0; i < 10; i++ {
		last = bm.RecordFailure("host1")
	}
	if last != 5*time.Second {
		t.Errorf("expected max 5s, got %v", last)
	}
}

func TestRecordSuccess_ResetFailures(t *testing.T) {
	cfg := DefaultBackoffConfig()
	bm := NewBackoffManager(cfg)
	bm.RecordFailure("host1")
	bm.RecordFailure("host1")
	if bm.Failures("host1") != 2 {
		t.Fatalf("expected 2 failures")
	}
	bm.RecordSuccess("host1")
	if bm.Failures("host1") != 0 {
		t.Errorf("expected 0 failures after success, got %d", bm.Failures("host1"))
	}
}

func TestRecordFailure_IndependentTargets(t *testing.T) {
	cfg := DefaultBackoffConfig()
	bm := NewBackoffManager(cfg)
	bm.RecordFailure("hostA")
	bm.RecordFailure("hostA")
	bm.RecordFailure("hostB")

	if bm.Failures("hostA") != 2 {
		t.Errorf("hostA: expected 2, got %d", bm.Failures("hostA"))
	}
	if bm.Failures("hostB") != 1 {
		t.Errorf("hostB: expected 1, got %d", bm.Failures("hostB"))
	}
}
