package portwatch

import (
	"testing"
	"time"
)

func TestNewCircuitBreaker_InvalidMaxFailures(t *testing.T) {
	_, err := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 0, OpenDuration: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxFailures=0")
	}
}

func TestNewCircuitBreaker_InvalidOpenDuration(t *testing.T) {
	_, err := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 1, OpenDuration: 0})
	if err == nil {
		t.Fatal("expected error for OpenDuration=0")
	}
}

func TestNewCircuitBreaker_Valid(t *testing.T) {
	cb, err := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb == nil {
		t.Fatal("expected non-nil CircuitBreaker")
	}
}

func TestAllow_InitiallyClosed(t *testing.T) {
	cb, _ := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	if !cb.Allow("host1") {
		t.Fatal("expected Allow=true for fresh target")
	}
}

func TestRecordFailure_OpensCircuit(t *testing.T) {
	cfg := CircuitBreakerConfig{MaxFailures: 2, OpenDuration: time.Minute}
	cb, _ := NewCircuitBreaker(cfg)
	cb.RecordFailure("host1")
	if cb.State("host1") != CircuitClosed {
		t.Fatal("circuit should still be closed after 1 failure")
	}
	cb.RecordFailure("host1")
	if cb.State("host1") != CircuitOpen {
		t.Fatal("circuit should be open after reaching MaxFailures")
	}
	if cb.Allow("host1") {
		t.Fatal("Allow should return false when circuit is open")
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	cfg := CircuitBreakerConfig{MaxFailures: 1, OpenDuration: time.Minute}
	cb, _ := NewCircuitBreaker(cfg)
	cb.RecordFailure("host1")
	if cb.State("host1") != CircuitOpen {
		t.Fatal("expected open circuit")
	}
	cb.RecordSuccess("host1")
	if cb.State("host1") != CircuitClosed {
		t.Fatal("expected closed circuit after success")
	}
	if !cb.Allow("host1") {
		t.Fatal("Allow should return true after reset")
	}
}

func TestAllow_HalfOpenAfterExpiry(t *testing.T) {
	cfg := CircuitBreakerConfig{MaxFailures: 1, OpenDuration: 10 * time.Millisecond}
	cb, _ := NewCircuitBreaker(cfg)
	cb.RecordFailure("host1")
	time.Sleep(20 * time.Millisecond)
	if !cb.Allow("host1") {
		t.Fatal("expected Allow=true after open duration elapsed")
	}
	if cb.State("host1") != CircuitHalfOpen {
		t.Fatal("expected half-open state")
	}
}

func TestIndependentTargets(t *testing.T) {
	cfg := CircuitBreakerConfig{MaxFailures: 1, OpenDuration: time.Minute}
	cb, _ := NewCircuitBreaker(cfg)
	cb.RecordFailure("host1")
	if !cb.Allow("host2") {
		t.Fatal("host2 should be unaffected by host1 failures")
	}
}

func TestCircuitState_String(t *testing.T) {
	if CircuitClosed.String() != "closed" {
		t.Errorf("unexpected: %s", CircuitClosed.String())
	}
	if CircuitOpen.String() != "open" {
		t.Errorf("unexpected: %s", CircuitOpen.String())
	}
	if CircuitHalfOpen.String() != "half-open" {
		t.Errorf("unexpected: %s", CircuitHalfOpen.String())
	}
}
