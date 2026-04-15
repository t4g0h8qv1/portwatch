package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestLoad_MissingFile(t *testing.T) {
	h, err := history.Load("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestRecord_And_Load(t *testing.T) {
	path := tempPath(t)
	ports := []int{22, 80, 443}

	if err := history.Record(path, "localhost", ports); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	h, err := history.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	e := h.Entries[0]
	if e.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", e.Host)
	}
	if len(e.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(e.Ports))
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_Accumulates(t *testing.T) {
	path := tempPath(t)

	for i := 0; i < 3; i++ {
		if err := history.Record(path, "host", []int{8080}); err != nil {
			t.Fatalf("Record iteration %d failed: %v", i, err)
		}
	}

	h, err := history.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(h.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(h.Entries))
	}
}

func TestLast_Empty(t *testing.T) {
	h := &history.History{}
	_, ok := h.Last()
	if ok {
		t.Error("expected false for empty history")
	}
}

func TestLast_ReturnsNewest(t *testing.T) {
	path := tempPath(t)
	_ = history.Record(path, "host", []int{22})
	time.Sleep(time.Millisecond)
	_ = history.Record(path, "host", []int{443})

	h, _ := history.Load(path)
	last, ok := h.Last()
	if !ok {
		t.Fatal("expected entry")
	}
	if len(last.Ports) != 1 || last.Ports[0] != 443 {
		t.Errorf("expected last entry ports [443], got %v", last.Ports)
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := history.Load(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
