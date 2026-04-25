package portwatch

import (
	"testing"
	"time"
)

func TestDefaultSequenceConfig_Defaults(t *testing.T) {
	cfg := DefaultSequenceConfig()
	if cfg.MaxGap != 10*time.Minute {
		t.Fatalf("expected 10m, got %v", cfg.MaxGap)
	}
}

func TestNewScanSequenceManager_InvalidMaxGap(t *testing.T) {
	_, err := NewScanSequenceManager(SequenceConfig{MaxGap: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxGap")
	}
}

func TestNewScanSequenceManager_Valid(t *testing.T) {
	m, err := NewScanSequenceManager(DefaultSequenceConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	m, _ := NewScanSequenceManager(DefaultSequenceConfig())
	e1 := m.Record("host-a")
	if e1.Count != 1 {
		t.Fatalf("expected count 1, got %d", e1.Count)
	}
	e2 := m.Record("host-a")
	if e2.Count != 2 {
		t.Fatalf("expected count 2, got %d", e2.Count)
	}
}

func TestRecord_ResetsAfterMaxGap(t *testing.T) {
	cfg := SequenceConfig{MaxGap: 100 * time.Millisecond}
	m, _ := NewScanSequenceManager(cfg)

	base := time.Now()
	m.now = func() time.Time { return base }
	m.Record("host-a")
	m.Record("host-a")

	m.now = func() time.Time { return base.Add(200 * time.Millisecond) }
	e := m.Record("host-a")
	if e.Count != 1 {
		t.Fatalf("expected reset to 1, got %d", e.Count)
	}
}

func TestRecord_IndependentTargets(t *testing.T) {
	m, _ := NewScanSequenceManager(DefaultSequenceConfig())
	m.Record("host-a")
	m.Record("host-a")
	e := m.Record("host-b")
	if e.Count != 1 {
		t.Fatalf("expected host-b count 1, got %d", e.Count)
	}
}

func TestGet_Missing(t *testing.T) {
	m, _ := NewScanSequenceManager(DefaultSequenceConfig())
	_, ok := m.Get("ghost")
	if ok {
		t.Fatal("expected false for unknown target")
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	m, _ := NewScanSequenceManager(DefaultSequenceConfig())
	m.Record("host-a")
	e, ok := m.Get("host-a")
	if !ok {
		t.Fatal("expected true")
	}
	if e.Count != 1 {
		t.Fatalf("expected count 1, got %d", e.Count)
	}
}

func TestReset_ClearsSequence(t *testing.T) {
	m, _ := NewScanSequenceManager(DefaultSequenceConfig())
	m.Record("host-a")
	m.Record("host-a")
	m.Reset("host-a")
	_, ok := m.Get("host-a")
	if ok {
		t.Fatal("expected entry to be cleared after Reset")
	}
}

func TestTargets_ReturnsAll(t *testing.T) {
	m, _ := NewScanSequenceManager(DefaultSequenceConfig())
	m.Record("host-a")
	m.Record("host-b")
	targets := m.Targets()
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
}
