package portsnap_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/portwatch/internal/portsnap"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func sorted(s []int) []int {
	sort.Ints(s)
	return s
}

func TestTake_SetsFields(t *testing.T) {
	s := portsnap.Take("localhost", []int{80, 443})
	if s.Host != "localhost" {
		t.Fatalf("expected host localhost, got %s", s.Host)
	}
	if len(s.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(s.Ports))
	}
	if s.CapturedAt.IsZero() {
		t.Fatal("expected non-zero CapturedAt")
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	orig := portsnap.Take("10.0.0.1", []int{22, 8080})
	if err := portsnap.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := portsnap.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Host != orig.Host {
		t.Errorf("host mismatch: got %s", loaded.Host)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("ports length mismatch")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := portsnap.Load(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json"), 0600)
	_, err := portsnap.Load(path)
	if err == nil {
		t.Fatal("expected decode error")
	}
}

func TestDiff_NoChanges(t *testing.T) {
	a := portsnap.Take("h", []int{80, 443})
	b := portsnap.Take("h", []int{80, 443})
	opened, closed := portsnap.Diff(a, b)
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no diff, got opened=%v closed=%v", opened, closed)
	}
}

func TestDiff_Changes(t *testing.T) {
	a := portsnap.Take("h", []int{80, 443})
	b := portsnap.Take("h", []int{80, 8080})
	opened, closed := portsnap.Diff(a, b)
	if got := sorted(opened); len(got) != 1 || got[0] != 8080 {
		t.Errorf("opened: want [8080], got %v", got)
	}
	if got := sorted(closed); len(got) != 1 || got[0] != 443 {
		t.Errorf("closed: want [443], got %v", got)
	}
}
