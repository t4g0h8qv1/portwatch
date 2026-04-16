package portdiff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portdiff"
)

func TestFormat_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	portdiff.Format(&buf, portdiff.Result{})
	if !strings.Contains(buf.String(), "no port changes") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFormat_ShowsOpenedAndClosed(t *testing.T) {
	r := portdiff.Compute([]int{80}, []int{443})
	var buf bytes.Buffer
	portdiff.Format(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "+ port 443 opened") {
		t.Errorf("missing opened line in: %q", out)
	}
	if !strings.Contains(out, "- port 80 closed") {
		t.Errorf("missing closed line in: %q", out)
	}
}

func TestSummary_NoChanges(t *testing.T) {
	s := portdiff.Summary(portdiff.Result{})
	if s != "no changes" {
		t.Errorf("expected 'no changes', got %q", s)
	}
}

func TestSummary_OnlyOpened(t *testing.T) {
	r := portdiff.Compute([]int{}, []int{22, 80})
	s := portdiff.Summary(r)
	if s != "2 opened" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_Both(t *testing.T) {
	r := portdiff.Compute([]int{22}, []int{80})
	s := portdiff.Summary(r)
	if !strings.Contains(s, "opened") || !strings.Contains(s, "closed") {
		t.Errorf("unexpected summary: %q", s)
	}
}
