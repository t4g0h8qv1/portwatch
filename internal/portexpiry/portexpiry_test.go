package portexpiry_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/portexpiry"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "expiry.json")
}

func TestLoad_MissingFile(t *testing.T) {
	r, err := portexpiry.Load(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestTrack_AddsNewPorts(t *testing.T) {
	r, _ := portexpiry.Load(tempPath(t))
	now := time.Now()
	r.Track([]int{80, 443}, now)
	expired := r.Expired(0, now.Add(time.Second))
	if len(expired) != 2 {
		t.Fatalf("expected 2 expired entries, got %d", len(expired))
	}
}

func TestTrack_RemovesClosedPorts(t *testing.T) {
	r, _ := portexpiry.Load(tempPath(t))
	now := time.Now()
	r.Track([]int{80, 443}, now)
	r.Track([]int{80}, now.Add(time.Minute))
	expired := r.Expired(0, now.Add(2*time.Minute))
	for _, e := range expired {
		if e.Port == 443 {
			t.Fatal("port 443 should have been removed")
		}
	}
}

func TestExpired_RespectsMaxAge(t *testing.T) {
	r, _ := portexpiry.Load(tempPath(t))
	now := time.Now()
	r.Track([]int{22, 8080}, now)
	// only ports older than 1 hour should appear
	expired := r.Expired(time.Hour, now.Add(30*time.Minute))
	if len(expired) != 0 {
		t.Fatalf("expected 0 expired, got %d", len(expired))
	}
	expired = r.Expired(time.Hour, now.Add(2*time.Hour))
	if len(expired) != 2 {
		t.Fatalf("expected 2 expired, got %d", len(expired))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempPath(t)
	r, _ := portexpiry.Load(path)
	now := time.Now().Truncate(time.Second)
	r.Track([]int{80, 443}, now)
	if err := r.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	r2, err := portexpiry.Load(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	expired := r2.Expired(0, now.Add(time.Second))
	if len(expired) != 2 {
		t.Fatalf("expected 2 entries after reload, got %d", len(expired))
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json"), 0o600)
	_, err := portexpiry.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
