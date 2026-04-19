package portwatch

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestHealth_OK(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "localhost"
	m := Metrics{ConsecErrors: 0, LastScanAt: time.Now()}
	r := Health(cfg, m)
	if r.Status != HealthOK {
		t.Fatalf("expected ok, got %s", r.Status)
	}
}

func TestHealth_Degraded(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "localhost"
	m := Metrics{ConsecErrors: 3, LastError: errors.New("timeout")}
	r := Health(cfg, m)
	if r.Status != HealthDegraded {
		t.Fatalf("expected degraded, got %s", r.Status)
	}
	if r.LastError == nil {
		t.Fatal("expected last error to be set")
	}
}

func TestHealth_TargetPropagated(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "192.0.2.1"
	r := Health(cfg, Metrics{})
	if r.Target != "192.0.2.1" {
		t.Fatalf("unexpected target %q", r.Target)
	}
}

func TestWriteHealth_ContainsTarget(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "myhost"
	r := Health(cfg, Metrics{})
	var buf bytes.Buffer
	WriteHealth(&buf, r)
	if !strings.Contains(buf.String(), "myhost") {
		t.Fatalf("output missing target: %s", buf.String())
	}
}

func TestWriteHealth_NeverWhenZeroTime(t *testing.T) {
	r := HealthReport{Target: "x", LastScanAt: time.Time{}}
	var buf bytes.Buffer
	WriteHealth(&buf, r)
	if !strings.Contains(buf.String(), "never") {
		t.Fatalf("expected 'never' for zero time: %s", buf.String())
	}
}

func TestWriteHealth_ShowsError(t *testing.T) {
	r := HealthReport{
		Target:    "x",
		Status:    HealthDegraded,
		LastError: errors.New("connection refused"),
	}
	var buf bytes.Buffer
	WriteHealth(&buf, r)
	if !strings.Contains(buf.String(), "connection refused") {
		t.Fatalf("expected error in output: %s", buf.String())
	}
}

func TestHealthStatus_String(t *testing.T) {
	cases := []struct {
		s    HealthStatus
		want string
	}{
		{HealthOK, "ok"},
		{HealthDegraded, "degraded"},
		{HealthUnknown, "unknown"},
	}
	for _, c := range cases {
		if c.s.String() != c.want {
			t.Errorf("got %q want %q", c.s.String(), c.want)
		}
	}
}
