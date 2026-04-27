package portwatch

import (
	"strings"
	"testing"
	"time"
)

func makeWatchdogStatuses() []WatchdogStatus {
	now := time.Now()
	return []WatchdogStatus{
		{Target: "host1", LastSeen: now, Expired: false},
		{Target: "host2", LastSeen: now.Add(-10 * time.Minute), Expired: true},
	}
}

func TestWriteWatchdogTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteWatchdogTable(&sb, makeWatchdogStatuses())
	out := sb.String()
	for _, hdr := range []string{"TARGET", "LAST SEEN", "STATUS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWriteWatchdogTable_ShowsTarget(t *testing.T) {
	var sb strings.Builder
	WriteWatchdogTable(&sb, makeWatchdogStatuses())
	out := sb.String()
	if !strings.Contains(out, "host1") {
		t.Error("expected host1 in output")
	}
	if !strings.Contains(out, "host2") {
		t.Error("expected host2 in output")
	}
}

func TestWriteWatchdogTable_ShowsExpiredStatus(t *testing.T) {
	var sb strings.Builder
	WriteWatchdogTable(&sb, makeWatchdogStatuses())
	out := sb.String()
	if !strings.Contains(out, "expired") {
		t.Error("expected 'expired' status in output")
	}
	if !strings.Contains(out, "ok") {
		t.Error("expected 'ok' status in output")
	}
}

func TestWriteWatchdogTable_NeverWhenLastSeenZero(t *testing.T) {
	var sb strings.Builder
	WriteWatchdogTable(&sb, []WatchdogStatus{
		{Target: "ghost", LastSeen: time.Time{}, Expired: true},
	})
	out := sb.String()
	if !strings.Contains(out, "never") {
		t.Error("expected 'never' for zero LastSeen")
	}
}

func TestWatchdogSummary_NoTargets(t *testing.T) {
	s := WatchdogSummary(nil)
	if !strings.Contains(s, "no targets") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestWatchdogSummary_AllActive(t *testing.T) {
	now := time.Now()
	statuses := []WatchdogStatus{
		{Target: "a", LastSeen: now, Expired: false},
		{Target: "b", LastSeen: now, Expired: false},
	}
	s := WatchdogSummary(statuses)
	if !strings.Contains(s, "all") || !strings.Contains(s, "active") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestWatchdogSummary_SomeExpired(t *testing.T) {
	now := time.Now()
	statuses := []WatchdogStatus{
		{Target: "a", LastSeen: now, Expired: false},
		{Target: "b", LastSeen: now, Expired: true},
	}
	s := WatchdogSummary(statuses)
	if !strings.Contains(s, "1/2") {
		t.Errorf("unexpected summary: %s", s)
	}
}
