package portwatch

import (
	"bytes"
	"strings"
	"testing"
)

func TestHasChanges_True(t *testing.T) {
	e := ChangeEvent{Host: "localhost", Opened: []int{8080}}
	if !e.HasChanges() {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestHasChanges_False(t *testing.T) {
	e := ChangeEvent{Host: "localhost"}
	if e.HasChanges() {
		t.Fatal("expected HasChanges to be false")
	}
}

func TestSummary_NoChanges(t *testing.T) {
	e := ChangeEvent{Host: "host1"}
	got := e.Summary()
	if got != "host1: no changes" {
		t.Fatalf("unexpected summary: %q", got)
	}
}

func TestSummary_Opened(t *testing.T) {
	e := ChangeEvent{Host: "host1", Opened: []int{22, 80}}
	got := e.Summary()
	if !strings.Contains(got, "2 opened") {
		t.Fatalf("expected '2 opened' in summary, got: %q", got)
	}
}

func TestSummary_Both(t *testing.T) {
	e := ChangeEvent{Host: "host1", Opened: []int{443}, Closed: []int{8080}}
	got := e.Summary()
	if !strings.Contains(got, "1 opened") || !strings.Contains(got, "1 closed") {
		t.Fatalf("unexpected summary: %q", got)
	}
}

func TestWriteEvent_Output(t *testing.T) {
	var buf bytes.Buffer
	e := ChangeEvent{Host: "myhost", Opened: []int{9090}, Closed: []int{3306}}
	WriteEvent(&buf, e)
	out := buf.String()
	if !strings.Contains(out, "+ 9090") {
		t.Errorf("expected opened port in output, got: %q", out)
	}
	if !strings.Contains(out, "- 3306") {
		t.Errorf("expected closed port in output, got: %q", out)
	}
}
