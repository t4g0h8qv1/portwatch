package portreport_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portreport"
)

func makeEntries() []portreport.Entry {
	return []portreport.Entry{
		{Port: 443, Label: "https", Status: "open", Severity: "low"},
		{Port: 22, Label: "ssh", Status: "new", Severity: "critical"},
		{Port: 8080, Label: "http-alt", Status: "closed", Severity: "medium"},
	}
}

func TestNew_SortsByPort(t *testing.T) {
	r := portreport.New("localhost", makeEntries())
	if r.Entries[0].Port != 22 || r.Entries[1].Port != 443 || r.Entries[2].Port != 8080 {
		t.Fatalf("expected ports sorted ascending, got %v", r.Entries)
	}
}

func TestNew_SetsHost(t *testing.T) {
	r := portreport.New("10.0.0.1", makeEntries())
	if r.Host != "10.0.0.1" {
		t.Fatalf("expected host 10.0.0.1, got %s", r.Host)
	}
}

func TestCountByStatus(t *testing.T) {
	r := portreport.New("localhost", makeEntries())
	counts := r.CountByStatus()
	if counts["new"] != 1 || counts["closed"] != 1 || counts["open"] != 1 {
		t.Fatalf("unexpected counts: %v", counts)
	}
}

func TestCountByStatus_Empty(t *testing.T) {
	r := portreport.New("localhost", nil)
	if len(r.CountByStatus()) != 0 {
		t.Fatal("expected empty counts")
	}
}

func TestWriteSummary_ContainsHost(t *testing.T) {
	r := portreport.New("myhost", makeEntries())
	var buf bytes.Buffer
	r.WriteSummary(&buf)
	if !strings.Contains(buf.String(), "myhost") {
		t.Fatalf("summary missing host: %s", buf.String())
	}
}

func TestWriteSummary_ContainsPortLines(t *testing.T) {
	r := portreport.New("localhost", makeEntries())
	var buf bytes.Buffer
	r.WriteSummary(&buf)
	out := buf.String()
	for _, want := range []string{"22", "443", "8080", "ssh", "https"} {
		if !strings.Contains(out, want) {
			t.Errorf("summary missing %q", want)
		}
	}
}
