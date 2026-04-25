package portwatch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var scorecardNow = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestScorecardManager_RecordScan(t *testing.T) {
	sm := NewScorecardManager()
	sm.RecordScan("host-a", scorecardNow)
	sm.RecordScan("host-a", scorecardNow.Add(time.Minute))

	e, ok := sm.Get("host-a")
	if !ok {
		t.Fatal("expected entry for host-a")
	}
	if e.TotalScans != 2 {
		t.Errorf("TotalScans: got %d, want 2", e.TotalScans)
	}
	if !e.LastScan.Equal(scorecardNow.Add(time.Minute)) {
		t.Errorf("LastScan: got %v", e.LastScan)
	}
}

func TestScorecardManager_RecordAlert(t *testing.T) {
	sm := NewScorecardManager()
	sm.RecordAlert("host-b", scorecardNow)

	e, ok := sm.Get("host-b")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.TotalAlerts != 1 {
		t.Errorf("TotalAlerts: got %d, want 1", e.TotalAlerts)
	}
	if !e.LastAlert.Equal(scorecardNow) {
		t.Errorf("LastAlert: got %v", e.LastAlert)
	}
}

func TestScorecardManager_RecordError(t *testing.T) {
	sm := NewScorecardManager()
	sm.RecordError("host-c")
	sm.RecordError("host-c")

	e, ok := sm.Get("host-c")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.TotalErrors != 2 {
		t.Errorf("TotalErrors: got %d, want 2", e.TotalErrors)
	}
}

func TestScorecardManager_Get_Missing(t *testing.T) {
	sm := NewScorecardManager()
	_, ok := sm.Get("unknown")
	if ok {
		t.Error("expected false for unknown target")
	}
}

func TestScorecardManager_All_SortedByTarget(t *testing.T) {
	sm := NewScorecardManager()
	sm.RecordScan("zz", scorecardNow)
	sm.RecordScan("aa", scorecardNow)
	sm.RecordScan("mm", scorecardNow)

	all := sm.All()
	if len(all) != 3 {
		t.Fatalf("len: got %d, want 3", len(all))
	}
	if all[0].Target != "aa" || all[1].Target != "mm" || all[2].Target != "zz" {
		t.Errorf("unexpected order: %v %v %v", all[0].Target, all[1].Target, all[2].Target)
	}
}

func TestWriteScorecardTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	WriteScorecardTable(&buf, nil)
	out := buf.String()
	for _, h := range []string{"TARGET", "SCANS", "ALERTS", "ERRORS", "LAST SCAN"} {
		if !strings.Contains(out, h) {
			t.Errorf("header %q not found in output", h)
		}
	}
}

func TestWriteScorecardTable_ShowsEntry(t *testing.T) {
	sm := NewScorecardManager()
	sm.RecordScan("myhost", scorecardNow)
	sm.RecordAlert("myhost", scorecardNow)
	sm.RecordError("myhost")

	var buf bytes.Buffer
	WriteScorecardTable(&buf, sm.All())
	out := buf.String()
	if !strings.Contains(out, "myhost") {
		t.Error("expected target name in output")
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Error("expected date in output")
	}
}
