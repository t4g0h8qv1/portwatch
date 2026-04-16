package portgroup_test

import (
	"sort"
	"testing"

	"github.com/example/portwatch/internal/portgroup"
)

func sorted(ports []int) []int {
	s := make([]int, len(ports))
	copy(s, ports)
	sort.Ints(s)
	return s
}

func TestRegister_And_Lookup(t *testing.T) {
	r := portgroup.New()
	if err := r.Register("web", []int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ports, err := r.Lookup("web")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if len(ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(ports))
	}
}

func TestLookup_Missing(t *testing.T) {
	r := portgroup.New()
	_, err := r.Lookup("missing")
	if err == nil {
		t.Error("expected error for missing group")
	}
}

func TestRegister_EmptyName(t *testing.T) {
	r := portgroup.New()
	if err := r.Register("", []int{80}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestRegister_InvalidPort(t *testing.T) {
	r := portgroup.New()
	if err := r.Register("bad", []int{0}); err == nil {
		t.Error("expected error for port 0")
	}
	if err := r.Register("bad", []int{65536}); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestRegister_EmptyPorts(t *testing.T) {
	r := portgroup.New()
	if err := r.Register("empty", []int{}); err == nil {
		t.Error("expected error for empty ports")
	}
}

func TestResolve_Dedup(t *testing.T) {
	r := portgroup.New()
	_ = r.Register("web", []int{80, 443})
	_ = r.Register("tls", []int{443, 8443})
	ports, err := r.Resolve([]string{"web", "tls"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := sorted(ports)
	want := []int{80, 443, 8443}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: expected %d, got %d", i, want[i], got[i])
		}
	}
}

func TestResolve_UnknownGroup(t *testing.T) {
	r := portgroup.New()
	_, err := r.Resolve([]string{"unknown"})
	if err == nil {
		t.Error("expected error for unknown group")
	}
}

func TestNames(t *testing.T) {
	r := portgroup.New()
	_ = r.Register("web", []int{80})
	_ = r.Register("db", []int{5432})
	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
