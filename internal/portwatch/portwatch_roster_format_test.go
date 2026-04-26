package portwatch

import (
	"strings"
	"testing"
	"time"
)

func makeRosterEntries() []RosterEntry {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []RosterEntry{
		{Target: "host1", Active: true, AddedAt: now, LastSeen: now.Add(time.Minute)},
		{Target: "host2", Active: false, AddedAt: now, LastSeen: time.Time{}},
	}
}

func TestWriteRosterTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteRosterTable(&sb, makeRosterEntries())
	out := sb.String()
	for _, h := range []string{"TARGET", "ACTIVE", "ADDED", "LAST SEEN"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestWriteRosterTable_ShowsTargets(t *testing.T) {
	var sb strings.Builder
	WriteRosterTable(&sb, makeRosterEntries())
	out := sb.String()
	if !strings.Contains(out, "host1") {
		t.Error("expected host1 in output")
	}
	if !strings.Contains(out, "host2") {
		t.Error("expected host2 in output")
	}
}

func TestWriteRosterTable_NeverWhenLastSeenZero(t *testing.T) {
	var sb strings.Builder
	WriteRosterTable(&sb, makeRosterEntries())
	if !strings.Contains(sb.String(), "never") {
		t.Error("expected 'never' for zero LastSeen")
	}
}

func TestRosterSummary_NoTargets(t *testing.T) {
	s := RosterSummary(nil)
	if !strings.Contains(s, "no targets") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestRosterSummary_AllActive(t *testing.T) {
	entries := []RosterEntry{
		{Target: "h1", Active: true},
		{Target: "h2", Active: true},
	}
	s := RosterSummary(entries)
	if !strings.Contains(s, "all active") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestRosterSummary_SomeInactive(t *testing.T) {
	entries := []RosterEntry{
		{Target: "h1", Active: true},
		{Target: "h2", Active: false},
	}
	s := RosterSummary(entries)
	if !strings.Contains(s, "inactive") {
		t.Errorf("expected 'inactive' in summary: %q", s)
	}
}
