package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestEvaluate_NoChanges(t *testing.T) {
	result := alert.Evaluate([]int{22, 80}, []int{22, 80})
	if result.Changed {
		t.Error("expected no changes")
	}
	if len(result.NewPorts) != 0 {
		t.Errorf("expected no new ports, got %v", result.NewPorts)
	}
	if len(result.GonePorts) != 0 {
		t.Errorf("expected no gone ports, got %v", result.GonePorts)
	}
}

func TestEvaluate_NewPorts(t *testing.T) {
	result := alert.Evaluate([]int{22}, []int{22, 8080})
	if !result.Changed {
		t.Error("expected changes")
	}
	if len(result.NewPorts) != 1 || result.NewPorts[0] != 8080 {
		t.Errorf("expected new port 8080, got %v", result.NewPorts)
	}
	if len(result.GonePorts) != 0 {
		t.Errorf("expected no gone ports, got %v", result.GonePorts)
	}
}

func TestEvaluate_GonePorts(t *testing.T) {
	result := alert.Evaluate([]int{22, 3306}, []int{22})
	if !result.Changed {
		t.Error("expected changes")
	}
	if len(result.GonePorts) != 1 || result.GonePorts[0] != 3306 {
		t.Errorf("expected gone port 3306, got %v", result.GonePorts)
	}
}

func TestEvaluate_BothChanges(t *testing.T) {
	result := alert.Evaluate([]int{22, 3306}, []int{22, 9090})
	if !result.Changed {
		t.Error("expected changes")
	}
	if len(result.NewPorts) != 1 || result.NewPorts[0] != 9090 {
		t.Errorf("unexpected new ports: %v", result.NewPorts)
	}
	if len(result.GonePorts) != 1 || result.GonePorts[0] != 3306 {
		t.Errorf("unexpected gone ports: %v", result.GonePorts)
	}
}

func TestStdoutNotifier_Output(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewStdoutNotifier(&buf)

	result := alert.Result{
		NewPorts:  []int{8080},
		GonePorts: []int{3306},
		Changed:   true,
	}

	if err := n.Notify(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "3306") {
		t.Errorf("expected 3306 in output, got: %s", out)
	}
}

func TestStdoutNotifier_NoOutputWhenUnchanged(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewStdoutNotifier(&buf)

	result := alert.Result{Changed: false}
	if err := n.Notify(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for unchanged result, got: %s", buf.String())
	}
}
