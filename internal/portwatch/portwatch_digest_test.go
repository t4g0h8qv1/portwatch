package portwatch

import (
	"testing"
	"time"
)

var digestNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestNewDigestManager_Empty(t *testing.T) {
	dm := NewDigestManager()
	if len(dm.Targets()) != 0 {
		t.Fatal("expected no targets")
	}
}

func TestRecord_ChangedOnFirstCall(t *testing.T) {
	dm := NewDigestManager()
	_, changed := dm.Record("host-a", []int{80, 443}, digestNow)
	if !changed {
		t.Fatal("expected changed=true on first record")
	}
}

func TestRecord_UnchangedWhenSamePorts(t *testing.T) {
	dm := NewDigestManager()
	dm.Record("host-a", []int{80, 443}, digestNow)
	_, changed := dm.Record("host-a", []int{443, 80}, digestNow.Add(time.Minute))
	if changed {
		t.Fatal("expected changed=false for same port set (different order)")
	}
}

func TestRecord_ChangedWhenPortsAlter(t *testing.T) {
	dm := NewDigestManager()
	dm.Record("host-a", []int{80}, digestNow)
	_, changed := dm.Record("host-a", []int{80, 8080}, digestNow.Add(time.Minute))
	if !changed {
		t.Fatal("expected changed=true when port set grows")
	}
}

func TestGet_MissingTarget(t *testing.T) {
	dm := NewDigestManager()
	_, ok := dm.Get("ghost")
	if ok {
		t.Fatal("expected ok=false for unknown target")
	}
}

func TestGet_ReturnsStoredEntry(t *testing.T) {
	dm := NewDigestManager()
	digest, _ := dm.Record("host-b", []int{22, 80}, digestNow)
	e, ok := dm.Get("host-b")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Digest != digest {
		t.Fatalf("digest mismatch: got %s want %s", e.Digest, digest)
	}
	if e.PortCount != 2 {
		t.Fatalf("expected PortCount 2, got %d", e.PortCount)
	}
}

func TestTargets_SortedOrder(t *testing.T) {
	dm := NewDigestManager()
	dm.Record("zebra", []int{80}, digestNow)
	dm.Record("alpha", []int{443}, digestNow)
	dm.Record("mango", []int{22}, digestNow)

	targets := dm.Targets()
	expected := []string{"alpha", "mango", "zebra"}
	for i, want := range expected {
		if targets[i] != want {
			t.Fatalf("targets[%d] = %s, want %s", i, targets[i], want)
		}
	}
}

func TestComputeDigest_DeterministicAndOrderIndependent(t *testing.T) {
	a := computeDigest([]int{80, 443, 22})
	b := computeDigest([]int{22, 443, 80})
	if a != b {
		t.Fatalf("expected same digest for same ports in different order: %s vs %s", a, b)
	}
}

func TestComputeDigest_EmptyPorts(t *testing.T) {
	d := computeDigest([]int{})
	if d == "" {
		t.Fatal("expected non-empty digest for empty port list")
	}
}
