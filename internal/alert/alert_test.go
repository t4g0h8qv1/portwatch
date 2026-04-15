package alert_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
)

// captureNotifier records the last alert it received.
type captureNotifier struct {
	Called bool
	Last   alert.Alert
}

func (c *captureNotifier) Notify(a alert.Alert) error {
	c.Called = true
	c.Last = a
	return nil
}

func TestEvaluate_NoChanges(t *testing.T) {
	n := &captureNotifier{}
	sentAlert, err := alert.Evaluate(n, baseline.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sentAlert {
		t.Error("expected no alert for empty diff")
	}
	if n.Called {
		t.Error("notifier should not have been called")
	}
}

func TestEvaluate_NewPorts(t *testing.T) {
	n := &captureNotifier{}
	diff := baseline.Diff{New: []int{8080, 9090}}
	sentAlert, err := alert.Evaluate(n, diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sentAlert {
		t.Error("expected alert to be sent")
	}
	if n.Last.Level != alert.LevelError {
		t.Errorf("expected ERROR level, got %s", n.Last.Level)
	}
	if len(n.Last.NewPorts) != 2 {
		t.Errorf("expected 2 new ports, got %d", len(n.Last.NewPorts))
	}
}

func TestEvaluate_GonePorts(t *testing.T) {
	n := &captureNotifier{}
	diff := baseline.Diff{Gone: []int{22}}
	sentAlert, err := alert.Evaluate(n, diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sentAlert {
		t.Error("expected alert to be sent")
	}
	if n.Last.Level != alert.LevelWarn {
		t.Errorf("expected WARN level, got %s", n.Last.Level)
	}
}

func TestStdoutNotifier_Output(t *testing.T) {
	var buf strings.Builder
	n := &alert.StdoutNotifier{Writer: &buf}
	a := alert.Alert{
		Level:    alert.LevelError,
		Message:  "Port state change detected",
		NewPorts: []int{4444},
	}
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
	if !strings.Contains(out, "4444") {
		t.Errorf("expected port 4444 in output, got: %s", out)
	}
}
