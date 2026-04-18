package portpolicy_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portpolicy"
)

func makeResults() []portpolicy.Result {
	p := portpolicy.New([]int{80, 443}, []int{22})
	return p.EvaluateAll([]int{80, 22, 9090})
}

func TestWriteTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	portpolicy.WriteTable(&buf, makeResults())
	out := buf.String()
	if !strings.Contains(out, "PORT") || !strings.Contains(out, "STATUS") {
		t.Errorf("expected headers in output: %s", out)
	}
}

func TestWriteTable_ContainsPorts(t *testing.T) {
	var buf bytes.Buffer
	portpolicy.WriteTable(&buf, makeResults())
	out := buf.String()
	for _, want := range []string{"80", "22", "9090"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected port %s in output", want)
		}
	}
}

func TestSummary_Counts(t *testing.T) {
	results := makeResults()
	s := portpolicy.Summary(results)
	if !strings.Contains(s, "1 allowed") {
		t.Errorf("expected 1 allowed in %q", s)
	}
	if !strings.Contains(s, "1 denied") {
		t.Errorf("expected 1 denied in %q", s)
	}
	if !strings.Contains(s, "1 unreviewed") {
		t.Errorf("expected 1 unreviewed in %q", s)
	}
}

func TestSummary_Empty(t *testing.T) {
	s := portpolicy.Summary(nil)
	if s != "0 allowed, 0 denied, 0 unreviewed" {
		t.Errorf("unexpected summary: %s", s)
	}
}
