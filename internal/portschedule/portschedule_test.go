package portschedule_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portschedule"
)

func TestNew_InvalidInterval(t *testing.T) {
	_, err := portschedule.New(0)
	if errtt.Fatal("expected error for zero interval")
	}
}

func TestNew_ValidInterval(t *testing.T) {
	s, err := portschedule.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRegister_And_Get(t *testing.T) {
	s, _ := portschedule.New(time.Minute)
	if err := s.Register("localhost", 0); err != nil {
		t.Fatalf("Register: %v", err)
	}
	e, ok := s.Get("localhost")
	if !ok {
		t.Fatal("expected entry for localhost")
	}
	if e.Interval != time.Minute {
		t.Errorf("interval = %v, want %v", e.Interval, time.Minute)
	}
}

func TestRegister_EmptyHost(t *testing.T) {
	s, _ := portschedule.New(time.Minute)
	if err := s.Register("", 0); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestDue_ImmediatelyAfterRegister(t *testing.T) {
	s, _ := portschedule.New(time.Minute)
	_ = s.Register("host1", 0)
	_ = s.Register("host2", 0)
	due := s.Due()
	if len(due) != 2 {
		t.Errorf("expected 2 due hosts, got %d", len(due))
	}
}

func TestAdvance_UpdatesNextRun(t *testing.T) {
	s, _ := portschedule.New(50 * time.Millisecond)
	_ = s.Register("host1", 0)
	_ = s.Advance("host1")

	due := s.Due()
	for _, h := range due {
		if h == "host1" {
			t.Error("host1 should not be due immediately after Advance")
		}
	}

	time.Sleep(60 * time.Millisecond)
	due = s.Due()
	found := false
	for _, h := range due {
		if h == "host1" {
			found = true
		}
	}
	if !found {
		t.Error("host1 should be due after interval elapsed")
	}
}

func TestAdvance_UnregisteredHost(t *testing.T) {
	s, _ := portschedule.New(time.Minute)
	if err := s.Advance("ghost"); err == nil {
		t.Fatal("expected error for unregistered host")
	}
}

func TestRegister_CustomInterval(t *testing.T) {
	s, _ := portschedule.New(time.Minute)
	_ = s.Register("host1", 5*time.Minute)
	e, _ := s.Get("host1")
	if e.Interval != 5*time.Minute {
		t.Errorf("interval = %v, want 5m", e.Interval)
	}
}
