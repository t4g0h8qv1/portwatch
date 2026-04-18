package portbudget_test

import (
	"testing"

	"github.com/user/portwatch/internal/portbudget"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := portbudget.New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNew_ValidMax(t *testing.T) {
	b, err := portbudget.New(5)
	if err != nil || b == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_WithinBudget(t *testing.T) {
	b, _ := portbudget.New(3)
	v := b.Check("localhost", []int{80, 443})
	if v != nil {
		t.Fatalf("expected no violation, got %v", v)
	}
}

func TestCheck_AtLimit(t *testing.T) {
	b, _ := portbudget.New(3)
	v := b.Check("localhost", []int{80, 443, 8080})
	if v != nil {
		t.Fatalf("expected no violation at exact limit, got %v", v)
	}
}

func TestCheck_ExceedsLimit(t *testing.T) {
	b, _ := portbudget.New(2)
	v := b.Check("localhost", []int{80, 443, 8080})
	if v == nil {
		t.Fatal("expected violation")
	}
	if v.Max != 2 || v.Actual != 3 {
		t.Fatalf("unexpected violation values: %+v", v)
	}
	if v.Host != "localhost" {
		t.Fatalf("unexpected host: %s", v.Host)
	}
}

func TestSetHost_OverridesDefault(t *testing.T) {
	b, _ := portbudget.New(2)
	_ = b.SetHost("special", 10)
	v := b.Check("special", []int{80, 443, 8080, 9090})
	if v != nil {
		t.Fatalf("expected no violation for host with higher budget, got %v", v)
	}
	// default still applies to other hosts
	v2 := b.Check("other", []int{80, 443, 8080})
	if v2 == nil {
		t.Fatal("expected violation for host using default budget")
	}
}

func TestSetHost_InvalidMax(t *testing.T) {
	b, _ := portbudget.New(5)
	err := b.SetHost("host", 0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestViolation_Error(t *testing.T) {
	b, _ := portbudget.New(1)
	v := b.Check("myhost", []int{80, 443})
	if v == nil {
		t.Fatal("expected violation")
	}
	msg := v.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}
