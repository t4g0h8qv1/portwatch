package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/output"
)

var fixedTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func makeResult(newPorts, gonePorts []int) output.Result {
	return output.Result{
		Host:      "localhost",
		OpenPorts: []int{22, 80, 443},
		NewPorts:  newPorts,
		GonePorts: gonePorts,
		ScannedAt: fixedTime,
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatText)
	if err := w.Write(makeResult([]int{8080}, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected host in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected new port 8080 in output, got: %s", out)
	}
	if strings.Contains(out, "Gone ports") {
		t.Errorf("did not expect 'Gone ports' line when empty")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)
	r := makeResult(nil, []int{9000})
	if err := w.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got output.Result
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", got.Host)
	}
	if len(got.GonePorts) != 1 || got.GonePorts[0] != 9000 {
		t.Errorf("expected gone port 9000, got %v", got.GonePorts)
	}
}

func TestWrite_NoChanges_TextOmitsChangeLines(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatText)
	if err := w.Write(makeResult(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "New ports") || strings.Contains(out, "Gone ports") {
		t.Errorf("expected no change lines, got: %s", out)
	}
}

func TestFormatDefault(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.Format("unknown"))
	if err := w.Write(makeResult(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output for unknown format (should default to text)")
	}
}
