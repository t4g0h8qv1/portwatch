package portping_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
)

func makeResults() []portping.Result {
	return []portping.Result{
		{Host: "127.0.0.1", Port: 80, Latency: 2 * time.Millisecond},
		{Host: "127.0.0.1", Port: 9999, Err: fmt.Errorf("refused")},
	}
}

func TestWriteTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	results := []portping.Result{
		{Host: "localhost", Port: 22, Latency: 1 * time.Millisecond},
	}
	if err := portping.WriteTable(&buf, results); err != nil {
		t.Fatalf("WriteTable: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"HOST", "PORT", "LATENCY", "STATUS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestWriteTable_ShowsError(t *testing.T) {
	var buf bytes.Buffer
	results := []portping.Result{
		{Host: "localhost", Port: 9999, Err: fmt.Errorf("connection refused")},
	}
	portping.WriteTable(&buf, results)
	if !strings.Contains(buf.String(), "error") {
		t.Error("expected error in output")
	}
}

func TestSummary_AllReachable(t *testing.T) {
	results := []portping.Result{
		{Host: "h", Port: 80, Latency: time.Millisecond},
		{Host: "h", Port: 443, Latency: time.Millisecond},
	}
	s := portping.Summary(results)
	if !strings.HasPrefix(s, "2/2") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_Empty(t *testing.T) {
	s := portping.Summary(nil)
	if s != "no ports probed" {
		t.Errorf("unexpected: %s", s)
	}
}
