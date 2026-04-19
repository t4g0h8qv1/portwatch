package portwatch

import (
	"errors"
	"testing"
	"time"
)

func TestRecordScan_IncrementsCounter(t *testing.T) {
	var m Metrics
	m.RecordScan()
	m.RecordScan()
	if m.Scans != 2 {
		t.Fatalf("expected 2 scans, got %d", m.Scans)
	}
}

func TestRecordAlert_IncrementsCounter(t *testing.T) {
	var m Metrics
	m.RecordAlert()
	if m.Alerts != 1 {
		t.Fatalf("expected 1 alert, got %d", m.Alerts)
	}
}

func TestRecordError_IncrementsCounter(t *testing.T) {
	var m Metrics
	err := errors.New("fail")
	m.RecordError(err)
	m.RecordError(err)
	if m.Errors != 2 {
		t.Fatalf("expected 2 errors, got %d", m.Errors)
	}
	if m.ConsecErrors != 2 {
		t.Fatalf("expected 2 consec errors, got %d", m.ConsecErrors)
	}
}

func TestRecordScan_ResetsConsecErrors(t *testing.T) {
	var m Metrics
	m.RecordError(errors.New("oops"))
	m.RecordScan()
	if m.ConsecErrors != 0 {
		t.Fatalf("expected consec errors reset, got %d", m.ConsecErrors)
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	var m Metrics
	m.RecordScan()
	snap := m.Snapshot()
	m.RecordScan()
	if snap.Scans != 1 {
		t.Fatalf("snapshot should not reflect later changes")
	}
}

func TestRecordScan_UpdatesLastScanAt(t *testing.T) {
	var m Metrics
	before := time.Now()
	m.RecordScan()
	if m.LastScanAt.Before(before) {
		t.Fatal("LastScanAt not updated")
	}
}
