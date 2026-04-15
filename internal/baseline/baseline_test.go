package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/baseline"
)

func TestNew(t *testing.T) {
	b := baseline.New("localhost", []int{443, 80, 22})
	if b.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", b.Host)
	}
	if len(b.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(b.Ports))
	}
	// ports should be sorted
	expected := []int{22, 80, 443}
	for i, p := range b.Ports {
		if p != expected[i] {
			t.Errorf("port[%d]: expected %d, got %d", i, expected[i], p)
		}
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	orig := baseline.New("192.168.1.1", []int{22, 80, 443})
	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Host != orig.Host {
		t.Errorf("host mismatch: got %s", loaded.Host)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("ports length mismatch: got %d", len(loaded.Ports))
	}
	if loaded.CreatedAt.Round(time.Second) != orig.CreatedAt.Round(time.Second) {
		t.Errorf("CreatedAt mismatch")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestDiff(t *testing.T) {
	b := baseline.New("localhost", []int{22, 80, 443})

	newPorts, missing := b.Diff([]int{22, 80, 8080})

	if len(newPorts) != 1 || newPorts[0] != 8080 {
		t.Errorf("expected new port 8080, got %v", newPorts)
	}
	if len(missing) != 1 || missing[0] != 443 {
		t.Errorf("expected missing port 443, got %v", missing)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	b := baseline.New("localhost", []int{22, 80})
	newPorts, missing := b.Diff([]int{22, 80})
	if len(newPorts) != 0 {
		t.Errorf("expected no new ports, got %v", newPorts)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing ports, got %v", missing)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	b := baseline.New("localhost", []int{80})
	err := b.Save(filepath.Join(os.DevNull, "subdir", "baseline.json"))
	if err == nil {
		t.Error("expected error for invalid save path, got nil")
	}
}
