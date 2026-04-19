package portwatch

import (
	"strings"
	"testing"
	"time"
)

func TestWriteMetrics_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	WriteMetrics(&sb, Metrics{})
	out := sb.String()
	for _, h := range []string{"METRIC", "VALUE"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestWriteMetrics_ShowsCounters(t *testing.T) {
	var sb strings.Builder
	m := Metrics{ScansTotal: 4, AlertsTotal: 2, ErrorsTotal: 1, OpenPortCount: 7}
	WriteMetrics(&sb, m)
	out := sb.String()
	for _, want := range []string{"4", "2", "1", "7"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteMetrics_NeverForZeroTime(t *testing.T) {
	var sb strings.Builder
	WriteMetrics(&sb, Metrics{})
	out := sb.String()
	if !strings.Contains(out, "never") {
		t.Error("expected 'never' for zero time fields")
	}
}

func TestWriteMetrics_FormatsTime(t *testing.T) {
	var sb strings.Builder
	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	m := Metrics{LastScanAt: fixed}
	WriteMetrics(&sb, m)
	out := sb.String()
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Errorf("expected formatted time in output, got:\n%s", out)
	}
}
