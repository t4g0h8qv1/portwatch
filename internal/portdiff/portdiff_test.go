package portdiff_test

import (
	"testing"

	"github.com/example/portwatch/internal/portdiff"
)

func TestCompute_NoChanges(t *testing.T) {
	r := portdiff.Compute([]int{80, 443}, []int{80, 443})
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestCompute_Opened(t *testing.T) {
	r := portdiff.Compute([]int{80}, []int{80, 8080})
	if len(r.Opened) != 1 || r.Opened[0].Port != 8080 {
		t.Fatalf("expected opened 8080, got %v", r.Opened)
	}
	if len(r.Closed) != 0 {
		t.Fatalf("unexpected closed ports: %v", r.Closed)
	}
}

func TestCompute_Closed(t *testing.T) {
	r := portdiff.Compute([]int{80, 443}, []int{80})
	if len(r.Closed) != 1 || r.Closed[0].Port != 443 {
		t.Fatalf("expected closed 443, got %v", r.Closed)
	}
}

func TestCompute_BothChanges(t *testing.T) {
	r := portdiff.Compute([]int{22, 80}, []int{80, 9000})
	if len(r.Opened) != 1 || r.Opened[0].Port != 9000 {
		t.Fatalf("unexpected opened: %v", r.Opened)
	}
	if len(r.Closed) != 1 || r.Closed[0].Port != 22 {
		t.Fatalf("unexpected closed: %v", r.Closed)
	}
}

func TestCompute_FromEmpty(t *testing.T) {
	r := portdiff.Compute([]int{}, []int{22, 80})
	if len(r.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(r.Opened))
	}
}

func TestCompute_ToEmpty(t *testing.T) {
	r := portdiff.Compute([]int{22, 80}, []int{})
	if len(r.Closed) != 2 {
		t.Fatalf("expected 2 closed, got %d", len(r.Closed))
	}
}

func TestEntry_StatusValues(t *testing.T) {
	r := portdiff.Compute([]int{80}, []int{443})
	if r.Opened[0].Status != "opened" {
		t.Errorf("expected 'opened', got %q", r.Opened[0].Status)
	}
	if r.Closed[0].Status != "closed" {
		t.Errorf("expected 'closed', got %q", r.Closed[0].Status)
	}
}
