package portwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestDefaultHeartbeatConfig_Defaults(t *testing.T) {
	cfg := DefaultHeartbeatConfig()
	if cfg.Interval != 5*time.Minute {
		t.Fatalf("expected 5m, got %v", cfg.Interval)
	}
}

func TestNewHeartbeatManager_InvalidInterval(t *testing.T) {
	_, err := NewHeartbeatManager(HeartbeatConfig{Interval: 0})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewHeartbeatManager_Valid(t *testing.T) {
	hm, err := NewHeartbeatManager(DefaultHeartbeatConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hm == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestBeat_And_Last(t *testing.T) {
	hm, _ := NewHeartbeatManager(DefaultHeartbeatConfig())
	now := time.Now()
	hm.now = func() time.Time { return now }

	hm.Beat("host1")
	last, ok := hm.Last("host1")
	if !ok {
		t.Fatal("expected heartbeat to be recorded")
	}
	if !last.Equal(now) {
		t.Fatalf("expected %v, got %v", now, last)
	}
}

func TestLast_Missing(t *testing.T) {
	hm, _ := NewHeartbeatManager(DefaultHeartbeatConfig())
	_, ok := hm.Last("unknown")
	if ok {
		t.Fatal("expected no heartbeat for unknown target")
	}
}

func TestIsSilent_NoObservation(t *testing.T) {
	hm, _ := NewHeartbeatManager(DefaultHeartbeatConfig())
	if !hm.IsSilent("host1") {
		t.Fatal("expected target with no heartbeat to be silent")
	}
}

func TestIsSilent_RecentBeat(t *testing.T) {
	hm, _ := NewHeartbeatManager(HeartbeatConfig{Interval: time.Minute})
	now := time.Now()
	hm.now = func() time.Time { return now }
	hm.Beat("host1")
	if hm.IsSilent("host1") {
		t.Fatal("expected recent beat to not be silent")
	}
}

func TestIsSilent_ExpiredBeat(t *testing.T) {
	hm, _ := NewHeartbeatManager(HeartbeatConfig{Interval: time.Minute})
	past := time.Now().Add(-2 * time.Minute)
	hm.now = func() time.Time { return past }
	hm.Beat("host1")
	hm.now = func() time.Time { return time.Now() }
	if !hm.IsSilent("host1") {
		t.Fatal("expected expired beat to be silent")
	}
}

func TestWriteHeartbeatTable_NoData(t *testing.T) {
	hm, _ := NewHeartbeatManager(DefaultHeartbeatConfig())
	var buf bytes.Buffer
	WriteHeartbeatTable(&buf, hm)
	if !strings.Contains(buf.String(), "no heartbeat data") {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestWriteHeartbeatTable_ContainsTarget(t *testing.T) {
	hm, _ := NewHeartbeatManager(DefaultHeartbeatConfig())
	hm.Beat("db.internal")
	var buf bytes.Buffer
	WriteHeartbeatTable(&buf, hm)
	out := buf.String()
	if !strings.Contains(out, "db.internal") {
		t.Fatalf("expected target in output, got: %s", out)
	}
	if !strings.Contains(out, "TARGET") {
		t.Fatalf("expected header in output, got: %s", out)
	}
}
