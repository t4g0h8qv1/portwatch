package portwatch

import (
	"testing"
	"time"
)

func TestDefaultWarmupConfig_Defaults(t *testing.T) {
	cfg := DefaultWarmupConfig()
	if cfg.MinScans < 1 {
		t.Fatalf("expected MinScans >= 1, got %d", cfg.MinScans)
	}
	if cfg.MaxWait <= 0 {
		t.Fatalf("expected MaxWait > 0, got %v", cfg.MaxWait)
	}
}

func TestNewWarmupManager_InvalidMinScans(t *testing.T) {
	cfg := DefaultWarmupConfig()
	cfg.MinScans = 0
	_, err := NewWarmupManager(cfg)
	if err == nil {
		t.Fatal("expected error for MinScans=0")
	}
}

func TestNewWarmupManager_InvalidMaxWait(t *testing.T) {
	cfg := DefaultWarmupConfig()
	cfg.MaxWait = 0
	_, err := NewWarmupManager(cfg)
	if err == nil {
		t.Fatal("expected error for MaxWait=0")
	}
}

func TestNewWarmupManager_Valid(t *testing.T) {
	_, err := NewWarmupManager(DefaultWarmupConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsWarm_NotWarmBeforeMinScans(t *testing.T) {
	cfg := WarmupConfig{MinScans: 3, MaxWait: time.Hour}
	wm, _ := NewWarmupManager(cfg)
	wm.RecordScan("host1")
	wm.RecordScan("host1")
	if wm.IsWarm("host1") {
		t.Fatal("expected host1 to not be warm after 2 scans")
	}
}

func TestIsWarm_WarmAfterMinScans(t *testing.T) {
	cfg := WarmupConfig{MinScans: 3, MaxWait: time.Hour}
	wm, _ := NewWarmupManager(cfg)
	wm.RecordScan("host1")
	wm.RecordScan("host1")
	wm.RecordScan("host1")
	if !wm.IsWarm("host1") {
		t.Fatal("expected host1 to be warm after 3 scans")
	}
}

func TestIsWarm_WarmAfterMaxWait(t *testing.T) {
	cfg := WarmupConfig{MinScans: 10, MaxWait: time.Millisecond}
	wm, _ := NewWarmupManager(cfg)
	past := time.Now().Add(-time.Second)
	wm.mu.Lock()
	wm.entries["host2"] = &warmupEntry{scans: 1, firstSeen: past}
	wm.mu.Unlock()
	if !wm.IsWarm("host2") {
		t.Fatal("expected host2 to be warm after MaxWait elapsed")
	}
}

func TestIsWarm_UnknownTarget(t *testing.T) {
	wm, _ := NewWarmupManager(DefaultWarmupConfig())
	if wm.IsWarm("unknown") {
		t.Fatal("expected unknown target to not be warm")
	}
}

func TestReset_ClearsWarmup(t *testing.T) {
	cfg := WarmupConfig{MinScans: 1, MaxWait: time.Hour}
	wm, _ := NewWarmupManager(cfg)
	wm.RecordScan("host3")
	if !wm.IsWarm("host3") {
		t.Fatal("expected host3 to be warm")
	}
	wm.Reset("host3")
	if wm.IsWarm("host3") {
		t.Fatal("expected host3 to not be warm after reset")
	}
}

func TestScans_ReturnsCount(t *testing.T) {
	wm, _ := NewWarmupManager(DefaultWarmupConfig())
	wm.RecordScan("host4")
	wm.RecordScan("host4")
	if got := wm.Scans("host4"); got != 2 {
		t.Fatalf("expected 2 scans, got %d", got)
	}
}

func TestScans_MissingTargetIsZero(t *testing.T) {
	wm, _ := NewWarmupManager(DefaultWarmupConfig())
	if got := wm.Scans("nobody"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
