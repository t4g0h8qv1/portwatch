package portwatch

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewAuditLog_InvalidMax(t *testing.T) {
	_, err := NewAuditLog(0)
	if err == nil {
		t.Fatal("expected error for maxLen=0")
	}
}

func TestNewAuditLog_Valid(t *testing.T) {
	al, err := NewAuditLog(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if al == nil {
		t.Fatal("expected non-nil AuditLog")
	}
}

func TestRecord_And_All(t *testing.T) {
	al, _ := NewAuditLog(5)
	al.Record(AuditEvent{Target: "host-a", Ports: []int{80, 443}})
	al.Record(AuditEvent{Target: "host-b", Ports: []int{22}})
	events := al.All()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Target != "host-a" {
		t.Errorf("expected host-a, got %s", events[0].Target)
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	al, _ := NewAuditLog(3)
	for i := 0; i < 5; i++ {
		al.Record(AuditEvent{Target: "host", Ports: []int{i}})
	}
	events := al.All()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	// oldest evicted: first retained event has port 2
	if events[0].Ports[0] != 2 {
		t.Errorf("expected port 2 in oldest retained, got %d", events[0].Ports[0])
	}
}

func TestRecord_SetsTimestampWhenZero(t *testing.T) {
	al, _ := NewAuditLog(5)
	before := time.Now()
	al.Record(AuditEvent{Target: "host-a"})
	after := time.Now()
	events := al.All()
	ts := events[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range", ts)
	}
}

func TestLast_Found(t *testing.T) {
	al, _ := NewAuditLog(10)
	al.Record(AuditEvent{Target: "alpha", Opened: []int{80}})
	al.Record(AuditEvent{Target: "beta", Opened: []int{22}})
	al.Record(AuditEvent{Target: "alpha", Opened: []int{443}})
	ev, ok := al.Last("alpha")
	if !ok {
		t.Fatal("expected to find last event for alpha")
	}
	if len(ev.Opened) != 1 || ev.Opened[0] != 443 {
		t.Errorf("expected opened=[443], got %v", ev.Opened)
	}
}

func TestLast_Missing(t *testing.T) {
	al, _ := NewAuditLog(5)
	_, ok := al.Last("ghost")
	if ok {
		t.Fatal("expected not found for unknown target")
	}
}

func TestWriteAuditTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteAuditTable(&sb, []AuditEvent{})
	out := sb.String()
	for _, hdr := range []string{"TARGET", "OPENED", "CLOSED", "TIMESTAMP"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWriteAuditTable_ShowsError(t *testing.T) {
	var sb strings.Builder
	WriteAuditTable(&sb, []AuditEvent{
		{Target: "host-x", Timestamp: time.Now(), Error: errors.New("timeout")},
	})
	out := sb.String()
	if !strings.Contains(out, "timeout") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}
