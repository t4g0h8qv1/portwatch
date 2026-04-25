package portwatch

import (
	"testing"
	"time"
)

func TestDefaultFenceConfig_Defaults(t *testing.T) {
	cfg := DefaultFenceConfig()
	if cfg.MaxAge != 10*time.Minute {
		t.Fatalf("expected 10m, got %v", cfg.MaxAge)
	}
}

func TestNewScanFenceManager_InvalidMaxAge(t *testing.T) {
	_, err := NewScanFenceManager(FenceConfig{MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxAge")
	}
}

func TestNewScanFenceManager_Valid(t *testing.T) {
	m, err := NewScanFenceManager(DefaultFenceConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestFence_IsFenced(t *testing.T) {
	m, _ := NewScanFenceManager(DefaultFenceConfig())
	m.Fence("host-a", "maintenance")
	if !m.IsFenced("host-a") {
		t.Fatal("expected host-a to be fenced")
	}
}

func TestIsFenced_NotFenced(t *testing.T) {
	m, _ := NewScanFenceManager(DefaultFenceConfig())
	if m.IsFenced("host-b") {
		t.Fatal("expected host-b not to be fenced")
	}
}

func TestUnfence_RemovesFence(t *testing.T) {
	m, _ := NewScanFenceManager(DefaultFenceConfig())
	m.Fence("host-a", "test")
	m.Unfence("host-a")
	if m.IsFenced("host-a") {
		t.Fatal("expected host-a to be unfenced")
	}
}

func TestFence_Expires(t *testing.T) {
	base := time.Now()
	m, _ := NewScanFenceManager(FenceConfig{MaxAge: 5 * time.Minute})
	m.now = func() time.Time { return base }
	m.Fence("host-a", "expire-test")

	m.now = func() time.Time { return base.Add(6 * time.Minute) }
	if m.IsFenced("host-a") {
		t.Fatal("expected fence to have expired")
	}
}

func TestReason_ReturnedCorrectly(t *testing.T) {
	m, _ := NewScanFenceManager(DefaultFenceConfig())
	m.Fence("host-a", "planned downtime")
	if got := m.Reason("host-a"); got != "planned downtime" {
		t.Fatalf("expected 'planned downtime', got %q", got)
	}
}

func TestReason_EmptyWhenNotFenced(t *testing.T) {
	m, _ := NewScanFenceManager(DefaultFenceConfig())
	if r := m.Reason("unknown"); r != "" {
		t.Fatalf("expected empty reason, got %q", r)
	}
}

func TestPrune_RemovesExpired(t *testing.T) {
	base := time.Now()
	m, _ := NewScanFenceManager(FenceConfig{MaxAge: 1 * time.Minute})
	m.now = func() time.Time { return base }
	m.Fence("host-a", "prune-test")

	m.now = func() time.Time { return base.Add(2 * time.Minute) }
	m.Prune()

	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.fences["host-a"]; ok {
		t.Fatal("expected expired entry to be pruned")
	}
}
