package portacl_test

import (
	"testing"

	"github.com/user/portwatch/internal/portacl"
)

func TestEvaluate_DefaultAllow(t *testing.T) {
	a := portacl.New()
	if got := a.Evaluate(80, portacl.Inbound); got != portacl.Allow {
		t.Fatalf("expected Allow, got %s", got)
	}
}

func TestEvaluate_DenyRule(t *testing.T) .New()
	_ = a.Add(portacl.Rule{Port: 22, Direction: portacl.Inbound, Action: portacl.Deny})
	if got := a.Evaluate(22, portacl.Inbound); got != portacl.Deny {
		t.Fatalf("expected Deny, got %s", got)
	}
}

func TestEvaluate_DirectionMismatch(t *testing.T) {
	a := portacl.New()
	_ = a.Add(portacl.Rule{Port: 22, Direction: portacl.Inbound, Action: portacl.// Outbound 22 should not match the inbound rule.
	if got := a.Evaluate(22, portacl.Outbound); got != portacl.Allow {
		t.Fatalf("expected Allow for outbound, got %s", got)
	}
}

func TestEvaluate_FirstRuleWins(t *testing.T) {
	a := portacl.New()
	_ = a.Add(portacl.Rule{Port: 443, Direction: portacl.Inbound, Action: portacl.Allow})
	_ = a.Add(portacl.Rule{Port: 443, Direction: portacl.Inbound, Action: portacl.Deny})
	if got := a.Evaluate(443, portacl.Inbound); got != portacl.Allow {
		t.Fatalf("expected Allow (first rule wins), got %s", got)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	a := portacl.New()
	if err := a.Add(portacl.Rule{Port: 0, Direction: portacl.Inbound, Action: portacl.Allow}); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := a.Add(portacl.Rule{Port: 99999, Direction: portacl.Inbound, Action: portacl.Allow}); err == nil {
		t.Fatal("expected error for port 99999")
	}
}

func TestAdd_InvalidDirection(t *testing.T) {
	a := portacl.New()
	err := a.Add(portacl.Rule{Port: 80, Direction: "sideways", Action: portacl.Allow})
	if err == nil {
		t.Fatal("expected error for invalid direction")
	}
}

func TestEvaluateAll(t *testing.T) {
	a := portacl.New()
	_ = a.Add(portacl.Rule{Port: 22, Direction: portacl.Inbound, Action: portacl.Deny})
	results := a.EvaluateAll([]int{22, 80, 443}, portacl.Inbound)
	if results[22] != portacl.Deny {
		t.Errorf("port 22: expected Deny")
	}
	if results[80] != portacl.Allow {
		t.Errorf("port 80: expected Allow")
	}
	if results[443] != portacl.Allow {
		t.Errorf("port 443: expected Allow")
	}
}
