package portwatch

import (
	"testing"
	"time"
)

func TestDefaultRetryConfig_Defaults(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.BaseDelay != 2*time.Second {
		t.Errorf("expected BaseDelay=2s, got %v", cfg.BaseDelay)
	}
	if cfg.MaxDelay != 30*time.Second {
		t.Errorf("expected MaxDelay=30s, got %v", cfg.MaxDelay)
	}
}

func TestNewScanRetryManager_InvalidMaxAttempts(t *testing.T) {
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 0
	_, err := NewScanRetryManager(cfg)
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestNewScanRetryManager_InvalidBaseDelay(t *testing.T) {
	cfg := DefaultRetryConfig()
	cfg.BaseDelay = 0
	_, err := NewScanRetryManager(cfg)
	if err == nil {
		t.Fatal("expected error for BaseDelay=0")
	}
}

func TestNewScanRetryManager_Valid(t *testing.T) {
	_, err := NewScanRetryManager(DefaultRetryConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShouldRetry_InitiallyTrue(t *testing.T) {
	m, _ := NewScanRetryManager(DefaultRetryConfig())
	if !m.ShouldRetry("host1") {
		t.Error("expected ShouldRetry=true before any failures")
	}
}

func TestNextDelay_IncreasesExponentially(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 4, BaseDelay: time.Second, MaxDelay: time.Minute}
	m, _ := NewScanRetryManager(cfg)

	d1 := m.NextDelay("host1")
	d2 := m.NextDelay("host1")
	d3 := m.NextDelay("host1")

	if d1 != time.Second {
		t.Errorf("expected 1s, got %v", d1)
	}
	if d2 != 2*time.Second {
		t.Errorf("expected 2s, got %v", d2)
	}
	if d3 != 4*time.Second {
		t.Errorf("expected 4s, got %v", d3)
	}
}

func TestNextDelay_CapsAtMax(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 10, BaseDelay: time.Second, MaxDelay: 3 * time.Second}
	m, _ := NewScanRetryManager(cfg)
	for i := 0; i < 5; i++ {
		d := m.NextDelay("host1")
		if d > 3*time.Second {
			t.Errorf("delay %v exceeded MaxDelay", d)
		}
	}
}

func TestNextDelay_ExhaustedReturnsZero(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, BaseDelay: time.Second, MaxDelay: time.Minute}
	m, _ := NewScanRetryManager(cfg)
	m.NextDelay("host1")
	m.NextDelay("host1")
	if d := m.NextDelay("host1"); d != 0 {
		t.Errorf("expected 0 after exhaustion, got %v", d)
	}
}

func TestShouldRetry_FalseAfterExhaustion(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, BaseDelay: time.Second, MaxDelay: time.Minute}
	m, _ := NewScanRetryManager(cfg)
	m.NextDelay("host1")
	m.NextDelay("host1")
	if m.ShouldRetry("host1") {
		t.Error("expected ShouldRetry=false after exhaustion")
	}
}

func TestReset_ClearsAttempts(t *testing.T) {
	m, _ := NewScanRetryManager(DefaultRetryConfig())
	m.NextDelay("host1")
	m.NextDelay("host1")
	m.Reset("host1")
	if m.Attempts("host1") != 0 {
		t.Errorf("expected 0 attempts after reset, got %d", m.Attempts("host1"))
	}
	if !m.ShouldRetry("host1") {
		t.Error("expected ShouldRetry=true after reset")
	}
}

func TestAttempts_IndependentTargets(t *testing.T) {
	m, _ := NewScanRetryManager(DefaultRetryConfig())
	m.NextDelay("hostA")
	m.NextDelay("hostA")
	m.NextDelay("hostB")
	if m.Attempts("hostA") != 2 {
		t.Errorf("expected hostA=2, got %d", m.Attempts("hostA"))
	}
	if m.Attempts("hostB") != 1 {
		t.Errorf("expected hostB=1, got %d", m.Attempts("hostB"))
	}
}
