package portwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewSignalManager_InvalidMaxAge(t *testing.T) {
	_, err := NewSignalManager(0)
	if err == nil {
		t.Fatal("expected error for zero maxAge")
	}
}

func TestNewSignalManager_Valid(t *testing.T) {
	m, err := NewSignalManager(time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestRecord_And_All(t *testing.T) {
	m, _ := NewSignalManager(time.Hour)
	if err := m.Record("host1", 80, SignalOpened); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := m.Record("host1", 443, SignalStable); err != nil {
		t.Fatalf("Record: %v", err)
	}
	sigs := m.All("host1")
	if len(sigs) != 2 {
		t.Fatalf("expected 2 signals, got %d", len(sigs))
	}
	if sigs[0].Port != 80 || sigs[0].Kind != SignalOpened {
		t.Errorf("unexpected first signal: %+v", sigs[0])
	}
}

func TestRecord_EmptyTarget_ReturnsError(t *testing.T) {
	m, _ := NewSignalManager(time.Hour)
	if err := m.Record("", 80, SignalOpened); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestAll_Missing_ReturnsEmpty(t *testing.T) {
	m, _ := NewSignalManager(time.Hour)
	if got := m.All("ghost"); len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestTargets_SortedAlphabetically(t *testing.T) {
	m, _ := NewSignalManager(time.Hour)
	_ = m.Record("zebra", 22, SignalOpened)
	_ = m.Record("alpha", 80, SignalClosed)
	targets := m.Targets()
	if len(targets) != 2 || targets[0] != "alpha" || targets[1] != "zebra" {
		t.Errorf("unexpected targets order: %v", targets)
	}
}

func TestWriteSignalTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	WriteSignalTable(&buf, []ScanSignal{
		{Target: "host1", Port: 22, Kind: SignalOpened, ObservedAt: time.Now()},
	})
	out := buf.String()
	for _, hdr := range []string{"TARGET", "PORT", "KIND", "OBSERVED AT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output:\n%s", hdr, out)
		}
	}
}

func TestSignalSummary_NoSignals(t *testing.T) {
	if got := SignalSummary(nil); got != "no signals recorded" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSignalSummary_WithSignals(t *testing.T) {
	sigs := []ScanSignal{
		{Kind: SignalOpened},
		{Kind: SignalClosed},
		{Kind: SignalOpened},
	}
	got := SignalSummary(sigs)
	if !strings.Contains(got, "3 signal(s)") {
		t.Errorf("unexpected summary: %q", got)
	}
	if !strings.Contains(got, "2 opened") {
		t.Errorf("expected 2 opened in summary: %q", got)
	}
}
