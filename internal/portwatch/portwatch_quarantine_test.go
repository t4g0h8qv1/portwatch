package portwatch

import (
	"testing"
	"time"
)

func TestNewQuarantineManager_Empty(t *testing.T) {
	q := NewQuarantineManager()
	if q.Count() != 0 {
		t.Fatalf("expected 0 entries, got %d", q.Count())
	}
}

func TestQuarantine_InvalidTarget(t *testing.T) {
	q := NewQuarantineManager()
	if err := q.Quarantine("", time.Minute); err != ErrEmptyQuarantineTarget {
		t.Fatalf("expected ErrEmptyQuarantineTarget, got %v", err)
	}
}

func TestQuarantine_InvalidDuration(t *testing.T) {
	q := NewQuarantineManager()
	if err := q.Quarantine("host1", 0); err != ErrInvalidQuarantineDuration {
		t.Fatalf("expected ErrInvalidQuarantineDuration, got %v", err)
	}
	if err := q.Quarantine("host1", -time.Second); err != ErrInvalidQuarantineDuration {
		t.Fatalf("expected ErrInvalidQuarantineDuration, got %v", err)
	}
}

func TestQuarantine_IsQuarantined(t *testing.T) {
	q := NewQuarantineManager()
	if err := q.Quarantine("host1", time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !q.IsQuarantined("host1") {
		t.Fatal("expected host1 to be quarantined")
	}
}

func TestQuarantine_NotQuarantined(t *testing.T) {
	q := NewQuarantineManager()
	if q.IsQuarantined("host1") {
		t.Fatal("expected host1 not to be quarantined")
	}
}

func TestQuarantine_Expires(t *testing.T) {
	now := time.Now()
	q := NewQuarantineManager()
	q.now = func() time.Time { return now }

	if err := q.Quarantine("host1", time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q.now = func() time.Time { return now.Add(2 * time.Minute) }
	if q.IsQuarantined("host1") {
		t.Fatal("expected quarantine to have expired")
	}
}

func TestRelease_RemovesQuarantine(t *testing.T) {
	q := NewQuarantineManager()
	_ = q.Quarantine("host1", time.Hour)
	q.Release("host1")
	if q.IsQuarantined("host1") {
		t.Fatal("expected host1 to be released")
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	q := NewQuarantineManager()
	q.now = func() time.Time { return now }

	_ = q.Quarantine("host1", time.Minute)
	_ = q.Quarantine("host2", time.Hour)

	q.now = func() time.Time { return now.Add(2 * time.Minute) }
	q.Prune()

	if q.Count() != 1 {
		t.Fatalf("expected 1 active entry after prune, got %d", q.Count())
	}
	if q.IsQuarantined("host1") {
		t.Fatal("expected host1 to be pruned")
	}
	if !q.IsQuarantined("host2") {
		t.Fatal("expected host2 to still be quarantined")
	}
}
