package suppress_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/suppress"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "suppress.json")
}

func TestLoad_MissingFile(t *testing.T) {
	l, err := suppress.Load(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(l.Entries))
	}
}

func TestAdd_And_IsSuppressed(t *testing.T) {
	path := tempPath(t)
	l, _ := suppress.Load(path)

	if err := l.Add(8080, "known dev port", time.Hour); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if !l.IsSuppressed(8080) {
		t.Error("expected port 8080 to be suppressed")
	}
	if l.IsSuppressed(9090) {
		t.Error("port 9090 should not be suppressed")
	}
}

func TestAdd_Persists(t *testing.T) {
	path := tempPath(t)
	l, _ := suppress.Load(path)
	_ = l.Add(443, "https", time.Hour)

	// reload from disk
	l2, err := suppress.Load(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !l2.IsSuppressed(443) {
		t.Error("expected port 443 to be suppressed after reload")
	}
}

func TestIsSuppressed_Expired(t *testing.T) {
	path := tempPath(t)
	l, _ := suppress.Load(path)
	// add with a TTL that has already passed
	_ = l.Add(22, "ssh", -time.Second)

	if l.IsSuppressed(22) {
		t.Error("expired suppression should not suppress port")
	}
}

func TestFilter_RemovesSuppressed(t *testing.T) {
	path := tempPath(t)
	l, _ := suppress.Load(path)
	_ = l.Add(8080, "dev", time.Hour)

	result := l.Filter([]int{22, 80, 8080, 443})
	for _, p := range result {
		if p == 8080 {
			t.Error("suppressed port 8080 should have been filtered out")
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 ports, got %d", len(result))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := suppress.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
