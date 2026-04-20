package portwatch

import (
	"testing"
	"time"
)

func TestNewScanSuppressManager_InvalidTTL(t *testing.T) {
	_, err := NewScanSuppressManager(SuppressConfig{DefaultTTL: 0})
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNewScanSuppressManager_Valid(t *testing.T) {
	m, err := NewScanSuppressManager(DefaultSuppressConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestSuppress_IsSuppressed(t *testing.T) {
	m, _ := NewScanSuppressManager(DefaultSuppressConfig())
	m.Suppress("host1", time.Hour)
	if !m.IsSuppressed("host1") {
		t.Fatal("expected host1 to be suppressed")
	}
}

func TestIsSuppressed_NotSuppressed(t *testing.T) {
	m, _ := NewScanSuppressManager(DefaultSuppressConfig())
	if m.IsSuppressed("host1") {
		t.Fatal("expected host1 to not be suppressed")
	}
}

func TestSuppress_Expires(t *testing.T) {
	m, _ := NewScanSuppressManager(DefaultSuppressConfig())
	fixed := time.Now()
	m.now = func() time.Time { return fixed }
	m.Suppress("host1", time.Second)
	m.now = func() time.Time { return fixed.Add(2 * time.Second) }
	if m.IsSuppressed("host1") {
		t.Fatal("expected suppression to have expired")
	}
}

func TestLift_RemovesSuppression(t *testing.T) {
	m, _ := NewScanSuppressManager(DefaultSuppressConfig())
	m.Suppress("host1", time.Hour)
	m.Lift("host1")
	if m.IsSuppressed("host1") {
		t.Fatal("expected host1 suppression to be lifted")
	}
}

func TestTargets_ReturnsSuppressed(t *testing.T) {
	m, _ := NewScanSuppressManager(DefaultSuppressConfig())
	m.Suppress("a", time.Hour)
	m.Suppress("b", time.Hour)
	targets := m.Targets()
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
}

func TestTargets_ExcludesExpired(t *testing.T) {
	m, _ := NewScanSuppressManager(DefaultSuppressConfig())
	fixed := time.Now()
	m.now = func() time.Time { return fixed }
	m.Suppress("expired", time.Millisecond)
	m.Suppress("active", time.Hour)
	m.now = func() time.Time { return fixed.Add(time.Second) }
	targets := m.Targets()
	if len(targets) != 1 || targets[0] != "active" {
		t.Fatalf("expected only active target, got %v", targets)
	}
}

func TestSuppress_UsesDefaultTTL(t *testing.T) {
	cfg := SuppressConfig{DefaultTTL: time.Hour}
	m, _ := NewScanSuppressManager(cfg)
	m.Suppress("host1", 0)
	if !m.IsSuppressed("host1") {
		t.Fatal("expected host1 to be suppressed using default TTL")
	}
}
