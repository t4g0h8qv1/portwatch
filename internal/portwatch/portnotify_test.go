package portwatch

import (
	"bytes"
	"strings"
	"testing"
)

func TestHasChanges_True(t *testing.T) {
	e := Event{Host: "h", Opened: []int{80}}
	if !HasChanges(e) {
		t.Fatal("expected changes")
	}
}

func TestHasChanges_False(t *testing.T) {
	e := Event{Host: "h"}
	if HasChanges(e) {
		t.Fatal("expected no changes")
	}
}

func TestSummary_NoChanges(t *testing.T) {
	e := Event{Host: "localhost"}
	if !strings.Contains(Summary(e), "no changes") {
		t.Fatalf("unexpected summary: %s", Summary(e))
	}
}

func TestSummary_Opened(t *testing.T) {
	e := Event{Host: "localhost", Opened: []int{443}}
	s := Summary(e)
	if !strings.Contains(s, "+1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestSummary_Both(t *testing.T) {
	e := Event{Host: "h", Opened: []int{80, 443}, Closed: []int{8080}}
	s := Summary(e)
	if !strings.Contains(s, "+2") || !strings.Contains(s, "-1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestWriteEvent_Output(t *testing.T) {
	var buf bytes.Buffer
	e := Event{Host: "myhost", Opened: []int{22}, Closed: []int{8080}}
	WriteEvent(&buf, e)
	out := buf.String()
	if !strings.Contains(out, "myhost") {
		t.Error("missing host")
	}
	if !strings.Contains(out, "+ 22") {
		t.Error("missing opened port")
	}
	if !strings.Contains(out, "- 8080") {
		t.Error("missing closed port")
	}
}
