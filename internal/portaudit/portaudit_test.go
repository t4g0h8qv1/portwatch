package portaudit_test

import (
	"testing"

	"github.com/user/portwatch/internal/portaudit"
	"github.com/user/portwatch/internal/severity"
)

func newAuditor(t *testing.T) *portaudit.Auditor {
	t.Helper()
	c, err := severity.New(nil)
	if err != nil {
		t.Fatalf("severity.New: %v", err)
	}
	return portaudit.New(c)
}

func TestRun_NoChanges(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("localhost", []int{80, 443}, []int{80, 443})
	if rec.HasChanges() {
		t.Fatal("expected no changes")
	}
	if len(rec.Stable) != 2 {
		t.Fatalf("stable: want 2, got %d", len(rec.Stable))
	}
}

func TestRun_NewPorts(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("localhost", []int{80}, []int{80, 8080})
	if !rec.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(rec.New) != 1 || rec.New[0] != 8080 {
		t.Fatalf("new ports: want [8080], got %v", rec.New)
	}
}

func TestRun_GonePorts(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("localhost", []int{80, 443}, []int{80})
	if !rec.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(rec.Gone) != 1 || rec.Gone[0] != 443 {
		t.Fatalf("gone ports: want [443], got %v", rec.Gone)
	}
}

func TestRun_SeverityElevated(t *testing.T) {
	a := newAuditor(t)
	// port 22 is critical in default classifier
	rec := a.Run("host", []int{}, []int{22})
	if rec.Severity < severity.Warning {
		t.Fatalf("expected elevated severity, got %v", rec.Severity)
	}
}

func TestRun_HostRecorded(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("192.168.1.1", nil, nil)
	if rec.Host != "192.168.1.1" {
		t.Fatalf("host: want 192.168.1.1, got %s", rec.Host)
	}
}

func TestRun_ScannedAtSet(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("h", nil, nil)
	if rec.ScannedAt.IsZero() {
		t.Fatal("ScannedAt should not be zero")
	}
}
