package portbaseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/portbaseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestLoad_MissingFile(t *testing.T) {
	s, err := portbaseline.Load(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Hosts()) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestSet_And_Get(t *testing.T) {
	s, _ := portbaseline.Load(tempPath(t))
	s.Set("host1", []int{80, 443, 22})
	e, ok := s.Get("host1")
	if !ok {
		t.Fatal("expected entry")
	}
	if len(e.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(e.Ports))
	}
	if e.Version != 1 {
		t.Fatalf("expected version 1, got %d", e.Version)
	}
}

func TestSet_IncrementsVersion(t *testing.T) {
	s, _ := portbaseline.Load(tempPath(t))
	s.Set("host1", []int{80})
	s.Set("host1", []int{80, 443})
	e, _ := s.Get("host1")
	if e.Version != 2 {
		t.Fatalf("expected version 2, got %d", e.Version)
	}
}

func TestSet_DeduplicatesPorts(t *testing.T) {
	s, _ := portbaseline.Load(tempPath(t))
	s.Set("host1", []int{80, 80, 443})
	e, _ := s.Get("host1")
	if len(e.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(e.Ports))
	}
}

func TestSaveAndLoad(t *testing.T) {
	p := tempPath(t)
	s, _ := portbaseline.Load(p)
	s.Set("alpha", []int{22, 80})
	s.Set("beta", []int{443})
	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	s2, err := portbaseline.Load(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(s2.Hosts()) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(s2.Hosts()))
	}
	e, ok := s2.Get("alpha")
	if !ok || len(e.Ports) != 2 {
		t.Fatal("alpha not restored correctly")
	}
}

func TestGet_Missing(t *testing.T) {
	s, _ := portbaseline.Load(tempPath(t))
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := portbaseline.Load(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
