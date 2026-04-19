package portwatch

import (
	"os"
	"path/filepath"
	"testing"
)

func checkpointPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestLoadCheckpoint_MissingFile(t *testing.T) {
	cp, err := LoadCheckpoint(checkpointPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cp.Get("localhost"); ok {
		t.Fatal("expected no entry for missing file")
	}
}

func TestRecord_And_Get(t *testing.T) {
	cp, _ := LoadCheckpoint(checkpointPath(t))
	if err := cp.Record("localhost", []int{80, 443}); err != nil {
		t.Fatalf("Record: %v", err)
	}
	e, ok := cp.Get("localhost")
	if !ok {
		t.Fatal("expected entry after Record")
	}
	if e.ScanCount != 1 {
		t.Errorf("ScanCount = %d, want 1", e.ScanCount)
	}
	if len(e.OpenPorts) != 2 {
		t.Errorf("OpenPorts len = %d, want 2", len(e.OpenPorts))
	}
}

func TestRecord_Accumulates(t *testing.T) {
	cp, _ := LoadCheckpoint(checkpointPath(t))
	cp.Record("host", []int{22})
	cp.Record("host", []int{22, 80})
	e, _ := cp.Get("host")
	if e.ScanCount != 2 {
		t.Errorf("ScanCount = %d, want 2", e.ScanCount)
	}
}

func TestRecord_Persists(t *testing.T) {
	path := checkpointPath(t)
	cp, _ := LoadCheckpoint(path)
	cp.Record("remote", []int{8080})

	cp2, err := LoadCheckpoint(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := cp2.Get("remote")
	if !ok {
		t.Fatal("entry not persisted")
	}
	if len(e.OpenPorts) != 1 || e.OpenPorts[0] != 8080 {
		t.Errorf("OpenPorts = %v, want [8080]", e.OpenPorts)
	}
}

func TestLoadCheckpoint_CorruptFile(t *testing.T) {
	path := checkpointPath(t)
	os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadCheckpoint(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
