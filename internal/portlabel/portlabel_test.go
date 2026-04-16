package portlabel_test

import (
	"testing"

	"github.com/example/portwatch/internal/portlabel"
)

func TestResolve_BuiltIn(t *testing.T) {
	l := portlabel.New(nil)
	if got := l.Resolve(443); got != "https" {
		t.Fatalf("expected https, got %s", got)
	}
}

func TestResolve_Unknown(t *testing.T) {
	l := portlabel.New(nil)
	if got := l.Resolve(9999); got != "unknown" {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestResolve_CustomOverridesBuiltIn(t *testing.T) {
	l := portlabel.New(map[int]string{80: "my-http"})
	if got := l.Resolve(80); got != "my-http" {
		t.Fatalf("expected my-http, got %s", got)
	}
}

func TestResolve_CustomNewEntry(t *testing.T) {
	l := portlabel.New(map[int]string{9000: "myapp"})
	if got := l.Resolve(9000); got != "myapp" {
		t.Fatalf("expected myapp, got %s", got)
	}
}

func TestLabel_String(t *testing.T) {
	l := portlabel.New(nil)
	lb := l.Label(22)
	if lb.String() != "22/ssh" {
		t.Fatalf("unexpected label string: %s", lb.String())
	}
}

func TestLabelAll(t *testing.T) {
	l := portlabel.New(nil)
	labels := l.LabelAll([]int{22, 80, 9999})
	if len(labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(labels))
	}
	expected := []string{"22/ssh", "80/http", "9999/unknown"}
	for i, lb := range labels {
		if lb.String() != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], lb.String())
		}
	}
}

func TestNew_NilCustom(t *testing.T) {
	l := portlabel.New(nil)
	if l == nil {
		t.Fatal("expected non-nil labeler")
	}
}
