package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/report"
)

func makeReport() *report.Report {
	return &report.Report{
		Target:    "localhost",
		ScannedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		OpenPorts: []int{22, 80, 443},
		NewPorts:  []int{8080},
		GonePorts: []int{3306},
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	r := makeReport()

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"localhost", "22, 80, 443", "8080", "3306", "Port Watch Report"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)
	r := makeReport()

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{`"target"`, `"localhost"`, `"open_ports"`, `"new_ports"`, "8080"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected JSON output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWrite_NoPorts(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	r := &report.Report{
		Target:    "example.com",
		ScannedAt: time.Now(),
		OpenPorts: []int{},
	}

	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected '(none)' for empty ports, got:\n%s", out)
	}
}

func TestFromAlert(t *testing.T) {
	result := alert.Result{
		NewPorts:  []int{9090},
		GonePorts: []int{21},
	}
	open := []int{22, 9090}

	r := report.FromAlert("192.168.1.1", result, open)

	if r.Target != "192.168.1.1" {
		t.Errorf("expected target 192.168.1.1, got %s", r.Target)
	}
	if len(r.NewPorts) != 1 || r.NewPorts[0] != 9090 {
		t.Errorf("unexpected new ports: %v", r.NewPorts)
	}
	if len(r.GonePorts) != 1 || r.GonePorts[0] != 21 {
		t.Errorf("unexpected gone ports: %v", r.GonePorts)
	}
	if len(r.OpenPorts) != 2 {
		t.Errorf("unexpected open ports: %v", r.OpenPorts)
	}
}

func TestNewWriter_NilOut(t *testing.T) {
	w := report.NewWriter(nil, report.FormatText)
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}
