package portwatch

import (
	"testing"
)

func TestNewTopologyManager_Empty(t *testing.T) {
	tm := NewTopologyManager()
	if tm == nil {
		t.Fatal("expected non-nil TopologyManager")
	}
}

func TestAddEdge_And_Neighbors(t *testing.T) {
	tm := NewTopologyManager()

	tm.AddEdge("host-a", "host-b")

	neighbors := tm.Neighbors("host-a")
	if len(neighbors) != 1 || neighbors[0] != "host-b" {
		t.Errorf("expected [host-b], got %v", neighbors)
	}
}

func TestAddEdge_Bidirectional(t *testing.T) {
	tm := NewTopologyManager()

	tm.AddEdge("host-a", "host-b")

	neighborsA := tm.Neighbors("host-a")
	neighborsB := tm.Neighbors("host-b")

	if len(neighborsA) != 1 || neighborsA[0] != "host-b" {
		t.Errorf("expected host-a neighbors [host-b], got %v", neighborsA)
	}
	if len(neighborsB) != 1 || neighborsB[0] != "host-a" {
		t.Errorf("expected host-b neighbors [host-a], got %v", neighborsB)
	}
}

func TestAddEdge_EmptyTarget_Ignored(t *testing.T) {
	tm := NewTopologyManager()

	tm.AddEdge("", "host-b")
	tm.AddEdge("host-a", "")

	if len(tm.Neighbors("host-a")) != 0 {
		t.Error("expected no neighbors for host-a after invalid edge")
	}
}

func TestAddEdge_Idempotent(t *testing.T) {
	tm := NewTopologyManager()

	tm.AddEdge("host-a", "host-b")
	tm.AddEdge("host-a", "host-b")

	neighbors := tm.Neighbors("host-a")
	if len(neighbors) != 1 {
		t.Errorf("expected 1 neighbor after duplicate AddEdge, got %d", len(neighbors))
	}
}

func TestRemoveEdge_RemovesConnection(t *testing.T) {
	tm := NewTopologyManager()

	tm.AddEdge("host-a", "host-b")
	tm.RemoveEdge("host-a", "host-b")

	if len(tm.Neighbors("host-a")) != 0 {
		t.Error("expected no neighbors after RemoveEdge")
	}
	if len(tm.Neighbors("host-b")) != 0 {
		t.Error("expected no neighbors for host-b after RemoveEdge")
	}
}

func TestNeighbors_UnknownTarget(t *testing.T) {
	tm := NewTopologyManager()

	neighbors := tm.Neighbors("unknown")
	if neighbors == nil {
		t.Error("expected non-nil slice for unknown target")
	}
	if len(neighbors) != 0 {
		t.Errorf("expected empty slice, got %v", neighbors)
	}
}

func TestTargets_ReturnsAllNodes(t *testing.T) {
	tm := NewTopologyManager()

	tm.AddEdge("host-a", "host-b")
	tm.AddEdge("host-b", "host-c")

	targets := tm.Targets()
	if len(targets) != 3 {
		t.Errorf("expected 3 targets, got %d: %v", len(targets), targets)
	}
}

func TestTargets_Empty(t *testing.T) {
	tm := NewTopologyManager()

	targets := tm.Targets()
	if len(targets) != 0 {
		t.Errorf("expected empty targets, got %v", targets)
	}
}
