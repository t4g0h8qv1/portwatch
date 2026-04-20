package portwatch

import (
	"strings"
	"testing"
	"time"
)

func TestWriteWindowTable_ContainsHeaders(t *testing.T) {
	m := NewScanWindowManager()
	var sb strings.Builder
	WriteWindowTable(&sb, m)
	if !strings.Contains(sb.String(), "TARGET") {
		t.Error("expected TARGET header")
	}
	if !strings.Contains(sb.String(), "WINDOW START") {
		t.Error("expected WINDOW START header")
	}
}

func TestWriteWindowTable_ShowsTarget(t *testing.T) {
	m := NewScanWindowManager()
	_ = m.Set("myhost", WindowConfig{Start: 8 * time.Hour, End: 18 * time.Hour})
	var sb strings.Builder
	WriteWindowTable(&sb, m)
	if !strings.Contains(sb.String(), "myhost") {
		t.Error("expected myhost in output")
	}
	if !strings.Contains(sb.String(), "08:00") {
		t.Error("expected 08:00 start in output")
	}
}

func TestWindowSummary_NoTargets(t *testing.T) {
	m := NewScanWindowManager()
	s := WindowSummary(m)
	if !strings.Contains(s, "no scan windows") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestWindowSummary_AllAllowed(t *testing.T) {
	m := NewScanWindowManager()
	m.now = fixedWindowNow(10, 0)
	_ = m.Set("a", WindowConfig{Start: 9 * time.Hour, End: 11 * time.Hour})
	_ = m.Set("b", WindowConfig{Start: 9 * time.Hour, End: 11 * time.Hour})
	s := WindowSummary(m)
	if !strings.Contains(s, "2/2") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestWindowSummary_SomeBlocked(t *testing.T) {
	m := NewScanWindowManager()
	m.now = fixedWindowNow(8, 0)
	_ = m.Set("a", WindowConfig{Start: 9 * time.Hour, End: 11 * time.Hour})
	_ = m.Set("b", WindowConfig{Start: 7 * time.Hour, End: 9 * time.Hour})
	s := WindowSummary(m)
	if !strings.Contains(s, "1/2") {
		t.Errorf("unexpected summary: %s", s)
	}
}
