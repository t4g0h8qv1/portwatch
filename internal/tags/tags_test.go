package tags_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/tags"
)

func TestNew_WellKnownPorts(t *testing.T) {
	r := tags.New()
	cases := map[int]string{
		22:  "ssh",
		80:  "http",
		443: "https",
	}
	for port, want := range cases {
		got, ok := r.Get(port)
		if !ok || got != want {
			t.Errorf("port %d: got %q ok=%v, want %q", port, got, ok, want)
		}
	}
}

func TestSet_And_Get(t *testing.T) {
	r := tags.New()
	r.Set(9200, "elasticsearch")
	label, ok := r.Get(9200)
	if !ok || label != "elasticsearch" {
		t.Fatalf("got %q ok=%v, want elasticsearch true", label, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	r := tags.New()
	_, ok := r.Get(65000)
	if ok {
		t.Fatal("expected no tag for unknown port")
	}
}

func TestLabel_Fallback(t *testing.T) {
	r := tags.New()
	got := r.Label(65000)
	if got != "port/65000" {
		t.Fatalf("got %q, want port/65000", got)
	}
}

func TestLabel_KnownPort(t *testing.T) {
	r := tags.New()
	got := r.Label(22)
	if got != "ssh" {
		t.Fatalf("got %q, want ssh", got)
	}
}

func TestLoadFile_MergesTags(t *testing.T) {
	extra := map[int]string{9200: "elasticsearch", 5601: "kibana"}
	data, _ := json.Marshal(extra)

	tmp := filepath.Join(t.TempDir(), "tags.json")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		t.Fatal(err)
	}

	r := tags.New()
	if err := r.LoadFile(tmp); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	for port, want := range extra {
		got, ok := r.Get(port)
		if !ok || got != want {
			t.Errorf("port %d: got %q ok=%v", port, got, ok)
		}
	}
	// built-ins still present
	if l := r.Label(22); l != "ssh" {
		t.Errorf("ssh gone after LoadFile: %q", l)
	}
}

func TestLoadFile_Missing(t *testing.T) {
	r := tags.New()
	if err := r.LoadFile("/no/such/file.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}
