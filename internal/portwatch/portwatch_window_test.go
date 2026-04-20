package portwatch

import (
	"testing"
	"time"
)

func fixedWindowNow(h, min int) func() time.Time {
	return func() time.Time {
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return midnight.Add(time.Duration(h)*time.Hour + time.Duration(min)*time.Minute)
	}
}

func TestNewScanWindowManager_Empty(t *testing.T) {
	m := NewScanWindowManager()
	if len(m.Targets()) != 0 {
		t.Fatal("expected no targets")
	}
}

func TestSet_InvalidTarget(t *testing.T) {
	m := NewScanWindowManager()
	err := m.Set("", WindowConfig{Start: time.Hour, End: 2 * time.Hour})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestSet_InvalidWindow(t *testing.T) {
	m := NewScanWindowManager()
	err := m.Set("host", WindowConfig{Start: 2 * time.Hour, End: time.Hour})
	if err == nil {
		t.Fatal("expected error when end <= start")
	}
}

func TestAllowed_NoWindow(t *testing.T) {
	m := NewScanWindowManager()
	if !m.Allowed("host") {
		t.Fatal("expected allowed when no window registered")
	}
}

func TestAllowed_WithinWindow(t *testing.T) {
	m := NewScanWindowManager()
	m.now = fixedWindowNow(10, 0) // 10:00
	_ = m.Set("host", WindowConfig{Start: 9 * time.Hour, End: 11 * time.Hour})
	if !m.Allowed("host") {
		t.Fatal("expected allowed within window")
	}
}

func TestAllowed_OutsideWindow(t *testing.T) {
	m := NewScanWindowManager()
	m.now = fixedWindowNow(8, 0) // 08:00
	_ = m.Set("host", WindowConfig{Start: 9 * time.Hour, End: 11 * time.Hour})
	if m.Allowed("host") {
		t.Fatal("expected denied outside window")
	}
}

func TestRemove_ClearsWindow(t *testing.T) {
	m := NewScanWindowManager()
	m.now = fixedWindowNow(8, 0)
	_ = m.Set("host", WindowConfig{Start: 9 * time.Hour, End: 11 * time.Hour})
	m.Remove("host")
	if !m.Allowed("host") {
		t.Fatal("expected allowed after window removed")
	}
}

func TestTargets_Listed(t *testing.T) {
	m := NewScanWindowManager()
	_ = m.Set("a", WindowConfig{Start: time.Hour, End: 2 * time.Hour})
	_ = m.Set("b", WindowConfig{Start: time.Hour, End: 2 * time.Hour})
	if len(m.Targets()) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(m.Targets()))
	}
}
