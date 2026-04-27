package portwatch

import (
	"sort"
	"testing"
	"time"
)

func TestDefaultWatchdogConfig_Defaults(t *testing.T) {
	cfg := DefaultWatchdogConfig()
	if cfg.MaxSilence <= 0 {
		t.Fatal("expected positive MaxSilence")
	}
	if cfg.CheckInterval <= 0 {
		t.Fatal("expected positive CheckInterval")
	}
}

func TestNewWatchdogManager_InvalidMaxSilence(t *testing.T) {
	_, err := NewWatchdogManager(WatchdogConfig{MaxSilence: 0, CheckInterval: time.Second})
	if err == nil {
		t.Fatal("expected error for zero MaxSilence")
	}
}

func TestNewWatchdogManager_InvalidCheckInterval(t *testing.T) {
	_, err := NewWatchdogManager(WatchdogConfig{MaxSilence: time.Minute, CheckInterval: 0})
	if err == nil {
		t.Fatal("expected error for zero CheckInterval")
	}
}

func TestNewWatchdogManager_Valid(t *testing.T) {
	_, err := NewWatchdogManager(DefaultWatchdogConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPing_EmptyTarget(t *testing.T) {
	w, _ := NewWatchdogManager(DefaultWatchdogConfig())
	if err := w.Ping(""); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestIsExpired_NotTracked(t *testing.T) {
	w, _ := NewWatchdogManager(DefaultWatchdogConfig())
	if !w.IsExpired("unknown") {
		t.Fatal("unknown target should be considered expired")
	}
}

func TestIsExpired_RecentPing(t *testing.T) {
	w, _ := NewWatchdogManager(DefaultWatchdogConfig())
	now := time.Now()
	w.now = func() time.Time { return now }
	_ = w.Ping("host1")
	if w.IsExpired("host1") {
		t.Fatal("recently pinged target should not be expired")
	}
}

func TestIsExpired_AfterMaxSilence(t *testing.T) {
	cfg := WatchdogConfig{MaxSilence: time.Minute, CheckInterval: time.Second}
	w, _ := NewWatchdogManager(cfg)
	base := time.Now()
	w.now = func() time.Time { return base }
	_ = w.Ping("host1")
	w.now = func() time.Time { return base.Add(2 * time.Minute) }
	if !w.IsExpired("host1") {
		t.Fatal("target should be expired after MaxSilence")
	}
}

func TestExpired_ReturnsMultiple(t *testing.T) {
	cfg := WatchdogConfig{MaxSilence: time.Minute, CheckInterval: time.Second}
	w, _ := NewWatchdogManager(cfg)
	base := time.Now()
	w.now = func() time.Time { return base }
	_ = w.Ping("a")
	_ = w.Ping("b")
	_ = w.Ping("c")
	w.now = func() time.Time { return base.Add(2 * time.Minute) }
	expired := w.Expired()
	sort.Strings(expired)
	if len(expired) != 3 {
		t.Fatalf("expected 3 expired targets, got %d", len(expired))
	}
}

func TestReset_RemovesTarget(t *testing.T) {
	w, _ := NewWatchdogManager(DefaultWatchdogConfig())
	now := time.Now()
	w.now = func() time.Time { return now }
	_ = w.Ping("host1")
	w.Reset("host1")
	if !w.IsExpired("host1") {
		t.Fatal("reset target should be considered expired (not tracked)")
	}
}
