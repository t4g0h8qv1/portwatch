package portwatch

import (
	"testing"
	"time"
)

func TestRecordScan_IncrementsCounter(t *testing.T) {
	m := &Metrics{}
	m.RecordScan(5)
	s := m.Snapshot()
	if s.ScansTotal != 1 {
		t.Fatalf("expected ScansTotal=1, got %d", s.ScansTotal)
	}
	if s.OpenPortCount != 5 {
		t.Fatalf("expected OpenPortCount=5, got %d", s.OpenPortCount)
	}
	if s.LastScanAt.IsZero() {
		t.Fatal("expected LastScanAt to be set")
	}
}

func TestRecordAlert_IncrementsCounter(t *testing.T) {
	m := &Metrics{}
	m.RecordAlert()
	m.RecordAlert()
	s := m.Snapshot()
	if s.AlertsTotal != 2 {
		t.Fatalf("expected AlertsTotal=2, got %d", s.AlertsTotal)
	}
	if s.LastAlertAt.IsZero() {
		t.Fatal("expected LastAlertAt to be set")
	}
}

func TestRecordError_IncrementsCounter(t *testing.T) {
	m := &Metrics{}
	m.RecordError()
	s := m.Snapshot()
	if s.ErrorsTotal != 1 {
		t.Fatalf("expected ErrorsTotal=1, got %d", s.ErrorsTotal)
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	m := &Metrics{}
	m.RecordScan(3)
	s1 := m.Snapshot()
	m.RecordScan(7)
	s2 := m.Snapshot()
	if s1.ScansTotal != 1 {
		t.Fatalf("snapshot should not be affected by later mutations")
	}
	if s2.ScansTotal != 2 {
		t.Fatalf("expected ScansTotal=2, got %d", s2.ScansTotal)
	}
}

func TestRecordScan_UpdatesLastScanAt(t *testing.T) {
	m := &Metrics{}
	before := time.Now()
	m.RecordScan(0)
	after := time.Now()
	s := m.Snapshot()
	if s.LastScanAt.Before(before) || s.LastScanAt.After(after) {
		t.Fatal("LastScanAt out of expected range")
	}
}
