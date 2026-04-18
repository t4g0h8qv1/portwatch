package portpolicy_test

import (
	"testing"

	"github.com/example/portwatch/internal/portpolicy"
)

func TestEvaluate_Allowed(t *testing.T) {
	p := portpolicy.New([]int{80, 443}, nil)
	r := p.Evaluate(80)
	if r.Status != portpolicy.Allowed {
		t.Fatalf("expected Allowed, got %s", r.Status)
	}
}

func TestEvaluate_Denied(t *testing.T) {
	p := portpolicy.New(nil, []int{22})
	r := p.Evaluate(22)
	if r.Status != portpolicy.Denied {
		t.Fatalf("expected Denied, got %s", r.Status)
	}
}

func TestEvaluate_Unreviewed(t *testing.T) {
	p := portpolicy.New([]int{80}, []int{22})
	r := p.Evaluate(9090)
	if r.Status != portpolicy.Unreviewed {
		t.Fatalf("expected Unreviewed, got %s", r.Status)
	}
}

func TestEvaluate_DeniedTakesPrecedence(t *testing.T) {
	p := portpolicy.New([]int{443}, []int{443})
	r := p.Evaluate(443)
	if r.Status != portpolicy.Denied {
		t.Fatalf("denied should take precedence, got %s", r.Status)
	}
}

func TestEvaluateAll_ReturnsAllResults(t *testing.T) {
	p := portpolicy.New([]int{80}, []int{22})
	results := p.EvaluateAll([]int{80, 22, 9090})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Status != portpolicy.Allowed {
		t.Errorf("port 80: expected Allowed")
	}
	if results[1].Status != portpolicy.Denied {
		t.Errorf("port 22: expected Denied")
	}
	if results[2].Status != portpolicy.Unreviewed {
		t.Errorf("port 9090: expected Unreviewed")
	}
}

func TestViolations_OnlyDenied(t *testing.T) {
	p := portpolicy.New([]int{80, 443}, []int{22, 23})
	v := p.Violations([]int{80, 22, 443, 23, 8080})
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
	for _, r := range v {
		if r.Status != portpolicy.Denied {
			t.Errorf("expected Denied, got %s", r.Status)
		}
	}
}

func TestViolations_None(t *testing.T) {
	p := portpolicy.New([]int{80}, nil)
	v := p.Violations([]int{80})
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}
