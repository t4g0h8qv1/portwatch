package portrange_test

import (
	"testing"

	"github.com/example/portwatch/internal/portrange"
)

func TestParse_Single(t *testing.T) {
	ports, err := portrange.Parse("80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 || ports[0] != 80 {
		t.Errorf("expected [80], got %v", ports)
	}
}

func TestParse_Range(t *testing.T) {
	ports, err := portrange.Parse("8080-8082")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{8080, 8081, 8082}
	for i, p := range expected {
		if ports[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, ports[i])
		}
	}
}

func TestParse_Mixed(t *testing.T) {
	ports, err := portrange.Parse("22,80,443,8000-8002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{22, 80, 443, 8000, 8001, 8002}
	if len(ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(ports))
	}
}

func TestParse_Dedup(t *testing.T) {
	ports, err := portrange.Parse("80,80,80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 {
		t.Errorf("expected 1 unique port, got %d", len(ports))
	}
}

func TestParse_InvalidPort(t *testing.T) {
	_, err := portrange.Parse("99999")
	if err == nil {
		t.Error("expected error for out-of-range port")
	}
}

func TestParse_InvalidRange(t *testing.T) {
	_, err := portrange.Parse("1000-500")
	if err == nil {
		t.Error("expected error for reversed range")
	}
}

func TestParse_NonNumeric(t *testing.T) {
	_, err := portrange.Parse("http")
	if err == nil {
		t.Error("expected error for non-numeric input")
	}
}
