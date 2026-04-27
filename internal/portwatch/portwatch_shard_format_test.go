package portwatch

import (
	"strings"
	"testing"
)

func TestWriteShardTable_ContainsHeaders(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 2})
	var sb strings.Builder
	WriteShardTable(&sb, m)
	out := sb.String()
	if !strings.Contains(out, "SHARD") {
		t.Errorf("expected SHARD header, got: %s", out)
	}
	if !strings.Contains(out, "TARGETS") {
		t.Errorf("expected TARGETS header, got: %s", out)
	}
}

func TestWriteShardTable_ShowsShardRows(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 3})
	m.Assign("host-1") //nolint:errcheck
	m.Assign("host-2") //nolint:errcheck
	var sb strings.Builder
	WriteShardTable(&sb, m)
	out := sb.String()
	// Three shard rows (0, 1, 2) plus header
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (header + 3 shards), got %d: %s", len(lines), out)
	}
}

func TestShardSummary_NoTargets(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 4})
	s := ShardSummary(m)
	if !strings.Contains(s, "no targets") {
		t.Errorf("expected 'no targets' in summary, got: %s", s)
	}
}

func TestShardSummary_WithTargets(t *testing.T) {
	m, _ := NewShardManager(ShardConfig{ShardCount: 2})
	m.Assign("a") //nolint:errcheck
	m.Assign("b") //nolint:errcheck
	m.Assign("c") //nolint:errcheck
	s := ShardSummary(m)
	if !strings.Contains(s, "3 targets") {
		t.Errorf("expected '3 targets' in summary, got: %s", s)
	}
	if !strings.Contains(s, "2 shards") {
		t.Errorf("expected '2 shards' in summary, got: %s", s)
	}
}
