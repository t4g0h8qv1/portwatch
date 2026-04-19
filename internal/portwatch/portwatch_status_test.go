package portwatch

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestIsHealthy_NoError(t *testing.T) {
	s := Status{}
	if !s.IsHealthy() {
		t.Fatal("expected healthy when LastError is nil")
	}
}

func TestIsHealthy_WithError(t *testing.T) {
	s := Status{LastError: errors.New("boom")}
	if s.IsHealthy() {
		t.Fatal("expected unhealthy when LastError is set")
	}
}

func TestWriteStatus_ContainsTarget(t *testing.T) {
	var b strings.Builder
	s := Status{
		Target:     "localhost",
		Ports:      []int{22, 80, 443},
		ScanCount:  5,
		AlertCount: 1,
		ErrorCount: 0,
		UpSince:    time.Now(),
		LastScan:   time.Now(),
	}
	WriteStatus(&b, s)
	out := b.String()
	for _, want := range []string{"localhost", "3", "5", "1", "yes"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWriteStatus_UnhealthyShowsError(t *testing.T) {
	var b strings.Builder
	s := Status{LastError: errors.New("scan failed")}
	WriteStatus(&b, s)
	if !strings.Contains(b.String(), "scan failed") {
		t.Errorf("expected error message in output")
	}
}

func TestWriteStatus_NeverWhenZeroTime(t *testing.T) {
	var b strings.Builder
	WriteStatus(&b, Status{})
	out := b.String()
	if !strings.Contains(out, "never") {
		t.Errorf("expected 'never' for zero times, got:\n%s", out)
	}
}
