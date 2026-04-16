package timeout_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/timeout"
)

func TestNew_InvalidTimeout(t *testing.T) {
	_, err := timeout.New(0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
	_, err = timeout.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestDefault_Returned(t *testing.T) {
	m, err := timeout.New(2 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Get("any-host"); got != 2*time.Second {
		t.Fatalf("expected 2s, got %v", got)
	}
}

func TestSet_OverridesDefault(t *testing.T) {
	m, _ := timeout.New(2 * time.Second)
	if err := m.Set("slow", 10*time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if got := m.Get("slow"); got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
	if got := m.Get("other"); got != 2*time.Second {
		t.Fatalf("expected default 2s, got %v", got)
	}
}

func TestSet_InvalidDuration(t *testing.T) {
	m, _ := timeout.New(2 * time.Second)
	if err := m.Set("host", 0); err == nil {
		t.Fatal("expected error for zero override")
	}
}

func TestRemove_RevertsToDefault(t *testing.T) {
	m, _ := timeout.New(3 * time.Second)
	_ = m.Set("h", 9*time.Second)
	m.Remove("h")
	if got := m.Get("h"); got != 3*time.Second {
		t.Fatalf("expected 3s after remove, got %v", got)
	}
}

func TestDefault_Method(t *testing.T) {
	m, _ := timeout.New(4 * time.Second)
	if m.Default() != 4*time.Second {
		t.Fatalf("Default() mismatch")
	}
}
