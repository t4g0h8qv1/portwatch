package portwatch

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// DefaultShardConfig returns a ShardConfig with sensible defaults.
func DefaultShardConfig() ShardConfig {
	return ShardConfig{
		ShardCount: 4,
	}
}

// ShardConfig holds configuration for the ShardManager.
type ShardConfig struct {
	ShardCount int
}

// shardEntry holds the targets assigned to a single shard.
type shardEntry struct {
	targets []string
}

// ShardManager distributes scan targets across a fixed number of shards.
type ShardManager struct {
	mu     sync.RWMutex
	shards []shardEntry
	index  map[string]int // target -> shard id
	cfg    ShardConfig
}

// NewShardManager creates a ShardManager with the given config.
func NewShardManager(cfg ShardConfig) (*ShardManager, error) {
	if cfg.ShardCount < 1 {
		return nil, errors.New("shard count must be at least 1")
	}
	return &ShardManager{
		shards: make([]shardEntry, cfg.ShardCount),
		index:  make(map[string]int),
		cfg:    cfg,
	}, nil
}

// Assign places a target into a shard deterministically by hash.
// If the target is already assigned, the existing shard id is returned.
func (m *ShardManager) Assign(target string) (int, error) {
	if target == "" {
		return 0, errors.New("target must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if id, ok := m.index[target]; ok {
		return id, nil
	}
	id := shardID(target, m.cfg.ShardCount)
	m.shards[id].targets = append(m.shards[id].targets, target)
	m.index[target] = id
	return id, nil
}

// ShardOf returns the shard id for a target, or an error if not assigned.
func (m *ShardManager) ShardOf(target string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	id, ok := m.index[target]
	if !ok {
		return 0, fmt.Errorf("target %q not assigned", target)
	}
	return id, nil
}

// Targets returns a sorted list of targets in the given shard.
func (m *ShardManager) Targets(shardID int) ([]string, error) {
	if shardID < 0 || shardID >= m.cfg.ShardCount {
		return nil, fmt.Errorf("shard id %d out of range", shardID)
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, len(m.shards[shardID].targets))
	copy(out, m.shards[shardID].targets)
	sort.Strings(out)
	return out, nil
}

// ShardCount returns the total number of shards.
func (m *ShardManager) ShardCount() int { return m.cfg.ShardCount }

// shardID computes a stable shard index for a target string.
func shardID(target string, n int) int {
	h := 0
	for _, c := range target {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h % n
}
