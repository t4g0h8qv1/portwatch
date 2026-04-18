package portquota_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/portquota"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := portquota.New(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
	_, err = portquota.New(-5)
	if err == nil {
		t.Fatal("expected error for max=-5")
	}
}

func TestNew_ValidMax(t *testing.T) {
	q, err := portquota.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Max() != 10 {
		t.Fatalf("expected max=10, got %d", q.Max())
	}
}

func TestCheck_WithinLimit(t *testing.T) {
	q, _ := portquota.New(5)
	ports := []int{80, 443, 8080}
	if err := q.Check(ports); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheck_AtLimit(t *testing.T) {
	q, _ := portquota.New(3)
	ports := []int{80, 443, 8080}
	if err := q.Check(ports); err != nil {
		t.Fatalf("expected no error at exact limit, got %v", err)
	}
}

func TestCheck_ExceedsLimit(t *testing.T) {
	q, _ := portquota.New(2)
	ports := []int{80, 443, 8080}
	err := q.Check(ports)
	if err == nil {
		t.Fatal("expected ErrQuotaExceeded")
	}
	if !errors.Is(err, portquota.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestRemaining_Positive(t *testing.T) {
	q, _ := portquota.New(10)
	ports := []int{80, 443}
	if r := q.Remaining(ports); r != 8 {
		t.Fatalf("expected 8 remaining, got %d", r)
	}
}

func TestRemaining_Negative(t *testing.T) {
	q, _ := portquota.New(2)
	ports := []int{80, 443, 8080, 9090}
	if r := q.Remaining(ports); r != -2 {
		t.Fatalf("expected -2 remaining, got %d", r)
	}
}
