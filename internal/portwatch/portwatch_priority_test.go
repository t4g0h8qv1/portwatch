package portwatch

import (
	"testing"
)

func TestDefaultPriorityConfig_Defaults(t *testing.T) {
	cfg := DefaultPriorityConfig()
	if cfg.DefaultLevel != PriorityNormal {
		t.Fatalf("expected default level %d, got %d", PriorityNormal, cfg.DefaultLevel)
	}
}

func TestNewScanPriorityManager_InvalidLevel(t *testing.T) {
	_, err := NewScanPriorityManager(PriorityConfig{DefaultLevel: 0})
	if err == nil {
		t.Fatal("expected error for zero default level")
	}
}

func TestNewScanPriorityManager_Valid(t *testing.T) {
	m, err := NewScanPriorityManager(DefaultPriorityConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestSet_And_Get(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	if err := m.Set("host1", PriorityHigh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Get("host1"); got != PriorityHigh {
		t.Fatalf("expected %d, got %d", PriorityHigh, got)
	}
}

func TestGet_ReturnsDefaultForUnknownTarget(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	if got := m.Get("unknown"); got != PriorityNormal {
		t.Fatalf("expected default %d, got %d", PriorityNormal, got)
	}
}

func TestSet_EmptyTarget(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	if err := m.Set("", PriorityHigh); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestSet_InvalidLevel(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	if err := m.Set("host1", 0); err == nil {
		t.Fatal("expected error for zero priority level")
	}
}

func TestReset_RevertsToDefault(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	_ = m.Set("host1", PriorityHigh)
	m.Reset("host1")
	if got := m.Get("host1"); got != PriorityNormal {
		t.Fatalf("expected default after reset, got %d", got)
	}
}

func TestTargets_ReturnsExplicitlySet(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	_ = m.Set("alpha", PriorityLow)
	_ = m.Set("beta", PriorityHigh)
	targets := m.Targets()
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
}

func TestTargets_EmptyAfterReset(t *testing.T) {
	m, _ := NewScanPriorityManager(DefaultPriorityConfig())
	_ = m.Set("host1", PriorityHigh)
	m.Reset("host1")
	if len(m.Targets()) != 0 {
		t.Fatal("expected no targets after reset")
	}
}
