package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidMax(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for max=0, got nil")
	}
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for interval=0, got nil")
	}
}

func TestAllow_ConsumesTokens(t *testing.T) {
	l, err := New(3, time.Minute)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Freeze time so no refill occurs.
	fixed := time.Now()
	l.clock = func() time.Time { return fixed }
	l.lastTick = fixed

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("Allow() returned false on call %d, expected true", i+1)
		}
	}
	// Fourth call should be denied.
	if l.Allow() {
		t.Fatal("Allow() returned true after exhausting tokens")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l, err := New(2, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	base := time.Now()
	l.clock = func() time.Time { return base }
	l.lastTick = base

	// Drain all tokens.
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected denial after draining tokens")
	}

	// Advance time by 1 second — should refill 2 tokens.
	base = base.Add(time.Second)
	if !l.Allow() {
		t.Fatal("expected Allow() after refill")
	}
}

func TestAllow_TokensCapAtMax(t *testing.T) {
	l, err := New(3, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	base := time.Now()
	l.clock = func() time.Time { return base }
	l.lastTick = base

	// Advance by 10 seconds — tokens should cap at max (3), not go to 30.
	base = base.Add(10 * time.Second)
	l.Allow() // trigger refill

	if l.Tokens() > l.max {
		t.Fatalf("tokens %f exceeded max %f", l.Tokens(), l.max)
	}
}

func TestTokens_InitialValue(t *testing.T) {
	l, _ := New(5, time.Minute)
	if got := l.Tokens(); got != 5.0 {
		t.Fatalf("expected initial tokens=5, got %f", got)
	}
}
