package portprofile_test

import (
	"testing"

	"github.com/example/portwatch/internal/portprofile"
)

func TestRegister_And_Get(t *testing.T) {
	r := portprofile.New()
	if err := r.Register("web", []int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := r.Get("web")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if p.Name != "web" {
		t.Errorf("name = %q, want web", p.Name)
	}
	if len(p.Ports) != 2 {
		t.Errorf("len(ports) = %d, want 2", len(p.Ports))
	}
}

func TestRegister_EmptyName(t *testing.T) {
	r := portprofile.New()
	if err := r.Register("", []int{80}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestRegister_EmptyPorts(t *testing.T) {
	r := portprofile.New()
	if err := r.Register("empty", []int{}); err == nil {
		t.Error("expected error for empty ports")
	}
}

func TestRegister_InvalidPort(t *testing.T) {
	r := portprofile.New()
	if err := r.Register("bad", []int{0}); err == nil {
		t.Error("expected error for port 0")
	}
	if err := r.Register("bad", []int{65536}); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestRegister_DeduplicatesPorts(t *testing.T) {
	r := portprofile.New()
	_ = r.Register("dup", []int{80, 80, 443, 443})
	p, _ := r.Get("dup")
	if len(p.Ports) != 2 {
		t.Errorf("expected 2 unique ports, got %d", len(p.Ports))
	}
}

func TestGet_Missing(t *testing.T) {
	r := portprofile.New()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestNames_Sorted(t *testing.T) {
	r := portprofile.New()
	_ = r.Register("ssh", []int{22})
	_ = r.Register("web", []int{80})
	_ = r.Register("db", []int{5432})
	names := r.Names()
	expected := []string{"db", "ssh", "web"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("names[%d] = %q, want %q", i, n, expected[i])
		}
	}
}

func TestDefault_HasWebProfile(t *testing.T) {
	r := portprofile.Default()
	p, err := r.Get("web")
	if err != nil {
		t.Fatalf("expected web profile: %v", err)
	}
	if len(p.Ports) == 0 {
		t.Error("web profile should have ports")
	}
}
