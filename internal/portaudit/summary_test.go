package portaudit_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portaudit"
	"github.com/user/portwatch/internal/severity"
)

func TestSummary_NoChanges(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("myhost", []int{80}, []int{80})
	s := rec.Summary()
	if !strings.Contains(s, "no changes") {
		t.Fatalf("expected 'no changes' in summary, got: %s", s)
	}
	if !strings.Contains(s, "myhost") {
		t.Fatalf("expected host in summary, got: %s", s)
	}
}

func TestSummary_WithChanges(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("myhost", []int{80}, []int{80, 9090})
	s := rec.Summary()
	if !strings.Contains(s, "+1") {
		t.Fatalf("expected '+1' in summary, got: %s", s)
	}
	if !strings.Contains(s, "myhost") {
		t.Fatalf("expected host in summary, got: %s", s)
	}
}

func TestSummary_IncludesSeverity(t *testing.T) {
	a := newAuditor(t)
	rec := a.Run("h", []int{}, []int{9999})
	s := rec.Summary()
	if !strings.Contains(s, severity.Info.String()) {
		t.Fatalf("expected severity in summary, got: %s", s)
	}
}
