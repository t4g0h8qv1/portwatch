package portlock_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/portwatch/internal/portlock"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portlock.json")
}

func TestLoad_MissingFile(t *testing.T) {
	s, err := portlock.Load(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsLocked(80) {
		t.Fatal("expected port 80 to not be locked")
	}
}

func TestLock_And_IsLocked(t *testing.T) {
	p := tempPath(t)
	s, _ := portlock.Load(p)
	if err := s.Lock(443, "https"); err != nil {
		t.Fatalf("Lock: %v", err)
	}
	if !s.IsLocked(443) {
		t.Fatal("expected 443 to be locked")
	}
	if s.IsLocked(80) {
		t.Fatal("expected 80 to not be locked")
	}
}

func TestLock_Persists(t *testing.T) {
	p := tempPath(t)
	s, _ := portlock.Load(p)
	_ = s.Lock(22, "ssh")

	s2, err := portlock.Load(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !s2.IsLocked(22) {
		t.Fatal("expected 22 to be locked after reload")
	}
}

func TestUnlock_RemovesPort(t *testing.T) {
	p := tempPath(t)
	s, _ := portlock.Load(p)
	_ = s.Lock(8080, "")
	if err := s.Unlock(8080); err != nil {
		t.Fatalf("Unlock: %v", err)
	}
	if s.IsLocked(8080) {
		t.Fatal("expected 8080 to be unlocked")
	}
	_, statErr := os.Stat(p)
	if statErr != nil {
		t.Fatalf("lock file missing after unlock: %v", statErr)
	}
}

func TestMissing_ReturnsAbsentLockedPorts(t *testing.T) {
	p := tempPath(t)
	s, _ := portlock.Load(p)
	_ = s.Lock(22, "")
	_ = s.Lock(80, "")
	_ = s.Lock(443, "")

	missing := s.Missing([]int{80})
	sort.Ints(missing)

	if len(missing) != 2 || missing[0] != 22 || missing[1] != 443 {
		t.Fatalf("unexpected missing ports: %v", missing)
	}
}

func TestMissing_NoneWhenAllOpen(t *testing.T) {
	s, _ := portlock.Load(tempPath(t))
	_ = s.Lock(22, "")
	missing := s.Missing([]int{22, 80})
	if len(missing) != 0 {
		t.Fatalf("expected no missing ports, got %v", missing)
	}
}
