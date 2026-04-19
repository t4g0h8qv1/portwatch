package portwatch_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwatch"
)

func TestWriteRunnerResult_ContainsHeaders(t *testing.T) {
	buf := &bytes.Buffer{}
	portwatch.WriteRunnerResult(buf, portwatch.RunnerResult{})
	out := buf.String()
	for _, h := range []string{"FIELD", "VALUE"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestWriteRunnerResult_NeverWhenZeroTime(t *testing.T) {
	buf := &bytes.Buffer{}
	portwatch.WriteRunnerResult(buf, portwatch.RunnerResult{ScansCompleted: 1})
	if !strings.Contains(buf.String(), "never") {
		t.Error("expected 'never' for zero LastScan")
	}
}

func TestWriteRunnerResult_ShowsLastScan(t *testing.T) {
	buf := &bytes.Buffer{}
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	portwatch.WriteRunnerResult(buf, portwatch.RunnerResult{ScansCompleted: 3, LastScan: ts})
	if !strings.Contains(buf.String(), "2024-06-01") {
		t.Error("expected formatted time in output")
	}
}

func TestRunnerSummary_NoScans(t *testing.T) {
	s := portwatch.RunnerSummary(portwatch.RunnerResult{})
	if s != "no scans completed" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestRunnerSummary_WithErrors(t *testing.T) {
	res := portwatch.RunnerResult{
		ScansCompleted: 5,
		Errors:         2,
		LastScan:       time.Now(),
	}
	s := portwatch.RunnerSummary(res)
	if !strings.Contains(s, "2 error(s)") {
		t.Errorf("expected error count in summary: %q", s)
	}
}

func TestRunnerSummary_NoErrors(t *testing.T) {
	res := portwatch.RunnerResult{
		ScansCompleted: 4,
		LastScan:       time.Now(),
	}
	s := portwatch.RunnerSummary(res)
	if strings.Contains(s, "error") {
		t.Errorf("unexpected error mention in summary: %q", s)
	}
}
