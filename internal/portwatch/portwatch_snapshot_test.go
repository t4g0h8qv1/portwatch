package portwatch

import (
	"bytes"
	"strings"
	"testing"
)

func TestSnapshotStore_RecordAndGet(t *testing.T) {
	s := NewSnapshotStore()
	s.Record("localhost", []int{80, 443})
	e, ok := s.Get("localhost")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if e.Target != "localhost" {
		t.Errorf("expected target localhost, got %s", e.Target)
	}
	if len(e.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(e.Ports))
	}
}

func TestSnapshotStore_Get_Missing(t *testing.T) {
	s := NewSnapshotStore()
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected no snapshot for unknown target")
	}
}

func TestSnapshotStore_Record_EmptyTarget(t *testing.T) {
	s := NewSnapshotStore()
	s.Record("", []int{22})
	if len(s.Targets()) != 0 {
		t.Fatal("empty target should not be stored")
	}
}

func TestSnapshotStore_Record_Overwrites(t *testing.T) {
	s := NewSnapshotStore()
	s.Record("host", []int{22})
	s.Record("host", []int{80, 443, 8080})
	e, _ := s.Get("host")
	if len(e.Ports) != 3 {
		t.Errorf("expected 3 ports after overwrite, got %d", len(e.Ports))
	}
}

func TestSnapshotStore_Targets(t *testing.T) {
	s := NewSnapshotStore()
	s.Record("a", []int{1})
	s.Record("b", []int{2})
	if len(s.Targets()) != 2 {
		t.Errorf("expected 2 targets, got %d", len(s.Targets()))
	}
}

func TestWriteSnapshotTable_ContainsHeaders(t *testing.T) {
	s := NewSnapshotStore()
	s.Record("myhost", []int{22, 80})
	var buf bytes.Buffer
	WriteSnapshotTable(&buf, s)
	out := buf.String()
	for _, h := range []string{"TARGET", "PORTS", "TAKEN AT"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestWriteSnapshotTable_ShowsTarget(t *testing.T) {
	s := NewSnapshotStore()
	s.Record("scanme.example.com", []int{443})
	var buf bytes.Buffer
	WriteSnapshotTable(&buf, s)
	if !strings.Contains(buf.String(), "scanme.example.com") {
		t.Error("expected target name in table output")
	}
}
