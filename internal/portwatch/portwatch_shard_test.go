package portwatch

import (
	"testing"
)

func TestDefaultShardConfig_Defaults(t *testing.T) {
	cfg := DefaultShardConfig()
	if cfg.ShardCount != 4 {
		t.Fatalf("expected ShardCount=4, got %d", cfg.ShardCount)
	}
}

func TestNewShardManager_InvalidCount(t *testing.T) {
	_, err := NewShardManager(ShardConfig{ShardCount: 0})
	if err == nil {
		t.Fatal("expected error for ShardCount=0")
	}
}

func TestNewShardManager_Valid(t *testing.T) {
	m, err := NewShardManager(ShardConfig{ShardCount: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ShardCount() != 3 {
		t.Fatalf("expected 3 shards, got %d", m.ShardCount())
	}
}

func TestAssign_EmptyTarget(t *testing.T) {
	m, _ := NewShardManager(DefaultShardConfig())
	_, err := m.Assign("")
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestAssign_ReturnsConsistentShard(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 4})
	id1, err := m.Assign("host-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id2, err := m.Assign("host-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("expected same shard, got %d and %d", id1, id2)
	}
}

func TestShardOf_NotAssigned(t *testing.T) {
	m, _ := NewShardManager(DefaultShardConfig())
	_, err := m.ShardOf("unknown")
	if err == nil {
		t.Fatal("expected error for unassigned target")
	}
}

func TestShardOf_AssignedTarget(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 4})
	assigned, _ := m.Assign("host-b")
	got, err := m.ShardOf("host-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != assigned {
		t.Fatalf("expected shard %d, got %d", assigned, got)
	}
}

func TestTargets_InvalidShard(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 2})
	_, err := m.Targets(5)
	if err == nil {
		t.Fatal("expected error for out-of-range shard")
	}
}

func TestTargets_ReturnsSorted(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 1})
	hosts := []string{"zebra", "alpha", "mango"}
	for _, h := range hosts {
		m.Assign(h) //nolint:errcheck
	}
	targets, err := m.Targets(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(targets); i++ {
		if targets[i-1] > targets[i] {
			t.Fatalf("targets not sorted: %v", targets)
		}
	}
}
