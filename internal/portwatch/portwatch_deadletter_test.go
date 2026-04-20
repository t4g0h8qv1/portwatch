package portwatch

import (
	"errors"
	"strings"
	"testing"
)

func TestNewDeadLetterQueue_InvalidSize(t *testing.T) {
	_, err := NewDeadLetterQueue(0)
	if !errors.Is(err, ErrInvalidDeadLetterSize) {
		t.Fatalf("expected ErrInvalidDeadLetterSize, got %v", err)
	}
}

func TestNewDeadLetterQueue_Valid(t *testing.T) {
	q, err := NewDeadLetterQueue(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Len() != 0 {
		t.Fatalf("expected empty queue")
	}
}

func TestPush_And_All(t *testing.T) {
	q, _ := NewDeadLetterQueue(5)
	q.Push("host-a", errors.New("timeout"), 1)
	q.Push("host-b", errors.New("refused"), 3)
	entries := q.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Target != "host-a" {
		t.Errorf("unexpected target: %s", entries[0].Target)
	}
}

func TestPush_EvictsOldestWhenFull(t *testing.T) {
	q, _ := NewDeadLetterQueue(2)
	q.Push("first", errors.New("e"), 1)
	q.Push("second", errors.New("e"), 1)
	q.Push("third", errors.New("e"), 1)
	entries := q.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Target != "second" {
		t.Errorf("expected oldest evicted, got %s", entries[0].Target)
	}
}

func TestClear_EmptiesQueue(t *testing.T) {
	q, _ := NewDeadLetterQueue(5)
	q.Push("host", errors.New("err"), 1)
	q.Clear()
	if q.Len() != 0 {
		t.Fatalf("expected empty queue after Clear")
	}
}

func TestWriteDeadLetterTable_ContainsHeaders(t *testing.T) {
	q, _ := NewDeadLetterQueue(5)
	q.Push("myhost", errors.New("scan failed"), 2)
	var sb strings.Builder
	WriteDeadLetterTable(&sb, q.All())
	out := sb.String()
	for _, hdr := range []string{"TARGET", "ATTEMPTS", "OCCURRED AT", "ERROR"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
	if !strings.Contains(out, "myhost") {
		t.Errorf("expected target in output")
	}
}

func TestDeadLetterSummary_Empty(t *testing.T) {
	s := DeadLetterSummary(nil)
	if !strings.Contains(s, "empty") {
		t.Errorf("expected 'empty' in summary, got %q", s)
	}
}

func TestDeadLetterSummary_WithEntries(t *testing.T) {
	q, _ := NewDeadLetterQueue(5)
	q.Push("h", errors.New("e"), 1)
	s := DeadLetterSummary(q.All())
	if !strings.Contains(s, "1 unprocessed") {
		t.Errorf("unexpected summary: %q", s)
	}
}
