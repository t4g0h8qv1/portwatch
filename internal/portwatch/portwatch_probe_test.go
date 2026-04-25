package portwatch

import (
	"testing"
	"time"
)

func TestDefaultProbeConfig_Defaults(t *testing.T) {
	cfg := DefaultProbeConfig()
	if cfg.MaxConsecutiveFailures != 3 {
		t.Errorf("expected MaxConsecutiveFailures=3, got %d", cfg.MaxConsecutiveFailures)
	}
	if cfg.ProbeInterval != 30*time.Second {
		t.Errorf("expected ProbeInterval=30s, got %v", cfg.ProbeInterval)
	}
	if cfg.RecoveryThreshold != 2 {
		t.Errorf("expected RecoveryThreshold=2, got %d", cfg.RecoveryThreshold)
	}
}

func TestNewProbeManager_InvalidMaxFailures(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.MaxConsecutiveFailures = 0
	_, err := NewProbeManager(cfg)
	if err == nil {
		t.Fatal("expected error for MaxConsecutiveFailures=0")
	}
}

func TestNewProbeManager_InvalidInterval(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.ProbeInterval = 0
	_, err := NewProbeManager(cfg)
	if err == nil {
		t.Fatal("expected error for ProbeInterval=0")
	}
}

func TestNewProbeManager_Valid(t *testing.T) {
	_, err := NewProbeManager(DefaultProbeConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsDead_NotDeadInitially(t *testing.T) {
	m, _ := NewProbeManager(DefaultProbeConfig())
	if m.IsDead("host-a") {
		t.Fatal("expected host-a to be alive initially")
	}
}

func TestRecordFailure_MarksDeadAtThreshold(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.MaxConsecutiveFailures = 2
	m, _ := NewProbeManager(cfg)

	m.RecordFailure("host-a")
	if m.IsDead("host-a") {
		t.Fatal("should not be dead after 1 failure")
	}
	m.RecordFailure("host-a")
	if !m.IsDead("host-a") {
		t.Fatal("expected host-a to be dead after 2 failures")
	}
}

func TestRecordSuccess_RecoversDead(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.MaxConsecutiveFailures = 1
	cfg.RecoveryThreshold = 2
	m, _ := NewProbeManager(cfg)

	m.RecordFailure("host-b")
	if !m.IsDead("host-b") {
		t.Fatal("expected dead after failure")
	}
	m.RecordSuccess("host-b")
	if !m.IsDead("host-b") {
		t.Fatal("should still be dead after 1 success (threshold=2)")
	}
	m.RecordSuccess("host-b")
	if m.IsDead("host-b") {
		t.Fatal("expected recovery after 2 successes")
	}
}

func TestRecordSuccess_ResetsFailureCount(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.MaxConsecutiveFailures = 3
	m, _ := NewProbeManager(cfg)

	m.RecordFailure("host-c")
	m.RecordFailure("host-c")
	m.RecordSuccess("host-c")
	m.RecordFailure("host-c")
	if m.IsDead("host-c") {
		t.Fatal("failure count should have reset after success")
	}
}

func TestReset_ClearsState(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.MaxConsecutiveFailures = 1
	m, _ := NewProbeManager(cfg)

	m.RecordFailure("host-d")
	if !m.IsDead("host-d") {
		t.Fatal("expected dead")
	}
	m.Reset("host-d")
	if m.IsDead("host-d") {
		t.Fatal("expected alive after reset")
	}
}

func TestIndependentTargets(t *testing.T) {
	cfg := DefaultProbeConfig()
	cfg.MaxConsecutiveFailures = 1
	m, _ := NewProbeManager(cfg)

	m.RecordFailure("host-x")
	if m.IsDead("host-y") {
		t.Fatal("host-y should not be affected by host-x failures")
	}
}
