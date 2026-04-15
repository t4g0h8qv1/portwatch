package throttle

import (
	"testing"
	"time"
)

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := New(0)
	if err != ErrInvalidCooldown {
		t.Fatalf("expected ErrInvalidCooldown, got %v", err)
	}

	_, err = New(-time.Second)
	if err != ErrInvalidCooldown {
		t.Fatalf("expected ErrInvalidCooldown for negative duration, got %v", err)
	}
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	th, err := New(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !th.Allow("port:22") {
		t.Error("expected first Allow call to return true")
	}
}

func TestAllow_SecondCallWithinCooldown(t *testing.T) {
	th, err := New(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	th.Allow("port:22")
	if th.Allow("port:22") {
		t.Error("expected second Allow within cooldown to return false")
	}
}

func TestAllow_AllowedAfterCooldown(t *testing.T) {
	now := time.Now()
	th, err := New(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	th.now = func() time.Time { return now }
	th.Allow("port:22")

	th.now = func() time.Time { return now.Add(2 * time.Minute) }
	if !th.Allow("port:22") {
		t.Error("expected Allow to return true after cooldown has elapsed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	th, err := New(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	th.Allow("port:80")
	th.Reset("port:80")
	if !th.Allow("port:80") {
		t.Error("expected Allow to return true after Reset")
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	th, err := New(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	th.now = func() time.Time { return now }
	th.Allow("port:443")
	th.Allow("port:8080")

	th.now = func() time.Time { return now.Add(90 * time.Second) }
	th.Prune()

	if len(th.last) != 0 {
		t.Errorf("expected all entries pruned, got %d remaining", len(th.last))
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th, err := New(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	th.Allow("port:22")
	if !th.Allow("port:80") {
		t.Error("expected different key to be allowed independently")
	}
}
