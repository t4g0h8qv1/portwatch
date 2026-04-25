package portwatch

import (
	"testing"
	"time"
)

func TestDefaultPressureConfig_Defaults(t *testing.T) {
	cfg := DefaultPressureConfig()
	if cfg.HighWatermark != 0.80 {
		t.Fatalf("expected HighWatermark 0.80, got %v", cfg.HighWatermark)
	}
	if cfg.LowWatermark != 0.50 {
		t.Fatalf("expected LowWatermark 0.50, got %v", cfg.LowWatermark)
	}
	if cfg.Window != 30*time.Second {
		t.Fatalf("expected Window 30s, got %v", cfg.Window)
	}
}

func TestNewScanPressureManager_InvalidMax(t *testing.T) {
	_, err := NewScanPressureManager(DefaultPressureConfig(), 0)
	if err == nil {
		t.Fatal("expected error for zero maxPerWindow")
	}
}

func TestNewScanPressureManager_InvalidHighWatermark(t *testing.T) {
	cfg := DefaultPressureConfig()
	cfg.HighWatermark = 1.5
	_, err := NewScanPressureManager(cfg, 10)
	if err == nil {
		t.Fatal("expected error for HighWatermark > 1")
	}
}

func TestNewScanPressureManager_InvalidLowWatermark(t *testing.T) {
	cfg := DefaultPressureConfig()
	cfg.LowWatermark = cfg.HighWatermark + 0.1
	_, err := NewScanPressureManager(cfg, 10)
	if err == nil {
		t.Fatal("expected error when LowWatermark >= HighWatermark")
	}
}

func TestNewScanPressureManager_Valid(t *testing.T) {
	m, err := NewScanPressureManager(DefaultPressureConfig(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestLevel_NormalWhenEmpty(t *testing.T) {
	m, _ := NewScanPressureManager(DefaultPressureConfig(), 10)
	if got := m.Level(); got != PressureNormal {
		t.Fatalf("expected Normal, got %v", got)
	}
}

func TestLevel_HighWhenAtWatermark(t *testing.T) {
	cfg := DefaultPressureConfig()
	cfg.Window = time.Minute
	m, _ := NewScanPressureManager(cfg, 10)
	// record 8 out of 10 — 0.80 == HighWatermark
	for i := 0; i < 8; i++ {
		m.Record()
	}
	if got := m.Level(); got != PressureHigh {
		t.Fatalf("expected High, got %v", got)
	}
}

func TestLoad_ReflectsRecordedScans(t *testing.T) {
	cfg := DefaultPressureConfig()
	cfg.Window = time.Minute
	m, _ := NewScanPressureManager(cfg, 10)
	m.Record()
	m.Record()
	if got := m.Load(); got != 0.2 {
		t.Fatalf("expected load 0.2, got %v", got)
	}
}

func TestRecord_PrunesExpiredObservations(t *testing.T) {
	cfg := DefaultPressureConfig()
	cfg.Window = 10 * time.Millisecond
	m, _ := NewScanPressureManager(cfg, 10)

	now := time.Now()
	m.nowFn = func() time.Time { return now }
	m.Record()
	m.Record()

	// advance time past the window
	m.nowFn = func() time.Time { return now.Add(20 * time.Millisecond) }
	m.Record() // triggers prune

	if got := m.Load(); got != 0.1 {
		t.Fatalf("expected load 0.1 after prune, got %v", got)
	}
}

func TestPressureLevel_String(t *testing.T) {
	if PressureNormal.String() != "normal" {
		t.Fatalf("unexpected string for Normal")
	}
	if PressureHigh.String() != "high" {
		t.Fatalf("unexpected string for High")
	}
}
