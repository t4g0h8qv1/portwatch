package portwatch

import (
	"testing"
	"time"
)

func TestDefaultEvictConfig_Defaults(t *testing.T) {
	cfg := DefaultEvictConfig()
	if cfg.MaxAge != 24*time.Hour {
		t.Fatalf("expected 24h, got %v", cfg.MaxAge)
	}
}

func TestNewEvictManager_InvalidMaxAge(t *testing.T) {
	_, err := NewEvictManager(EvictConfig{MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxAge")
	}
}

func TestNewEvictManager_Valid(t *testing.T) {
	_, err := NewEvictManager(DefaultEvictConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTouch_EmptyTarget(t *testing.T) {
	m, _ := NewEvictManager(DefaultEvictConfig())
	if err := m.Touch(""); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestShouldEvict_NotTracked(t *testing.T) {
	m, _ := NewEvictManager(DefaultEvictConfig())
	if m.ShouldEvict("host:22") {
		t.Fatal("untracked target should not be evicted")
	}
}

func TestShouldEvict_RecentlyTouched(t *testing.T) {
	m, _ := NewEvictManager(EvictConfig{MaxAge: time.Hour})
	_ = m.Touch("host:22")
	if m.ShouldEvict("host:22") {
		t.Fatal("recently touched target should not be evicted")
	}
}

func TestShouldEvict_Expired(t *testing.T) {
	base := time.Now()
	m, _ := NewEvictManager(EvictConfig{MaxAge: time.Minute})
	m.now = func() time.Time { return base }
	_ = m.Touch("host:22")
	m.now = func() time.Time { return base.Add(2 * time.Minute) }
	if !m.ShouldEvict("host:22") {
		t.Fatal("expired target should be evicted")
	}
}

func TestEvict_RemovesTarget(t *testing.T) {
	m, _ := NewEvictManager(DefaultEvictConfig())
	_ = m.Touch("host:22")
	m.Evict("host:22")
	if m.ShouldEvict("host:22") {
		t.Fatal("evicted target should no longer be tracked")
	}
	targets := m.Targets()
	if len(targets) != 0 {
		t.Fatalf("expected no targets, got %v", targets)
	}
}

func TestTargets_ReturnsAll(t *testing.T) {
	m, _ := NewEvictManager(DefaultEvictConfig())
	_ = m.Touch("host:22")
	_ = m.Touch("host:80")
	targets := m.Targets()
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
}
