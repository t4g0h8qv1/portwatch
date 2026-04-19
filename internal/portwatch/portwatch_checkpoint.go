package portwatch

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Checkpoint records the last successful scan time and result summary
// for each target, allowing resumption after restart.
type Checkpoint struct {
	mu      sync.RWMutex
	path    string
	entries map[string]CheckpointEntry
}

// CheckpointEntry holds persisted state for a single target.
type CheckpointEntry struct {
	Target    string    `json:"target"`
	LastScan  time.Time `json:"last_scan"`
	OpenPorts []int     `json:"open_ports"`
	ScanCount int       `json:"scan_count"`
}

// LoadCheckpoint reads a checkpoint file from disk.
// If the file does not exist, an empty Checkpoint is returned.
func LoadCheckpoint(path string) (*Checkpoint, error) {
	cp := &Checkpoint{
		path:    path,
		entries: make(map[string]CheckpointEntry),
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cp, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &cp.entries); err != nil {
		return nil, err
	}
	return cp, nil
}

// Record updates the checkpoint entry for the given target.
func (c *Checkpoint) Record(target string, ports []int) error {
	c.mu.Lock()
	prev := c.entries[target]
	c.entries[target] = CheckpointEntry{
		Target:    target,
		LastScan:  time.Now(),
		OpenPorts: ports,
		ScanCount: prev.ScanCount + 1,
	}
	c.mu.Unlock()
	return c.save()
}

// Get returns the checkpoint entry for a target, and whether it exists.
func (c *Checkpoint) Get(target string) (CheckpointEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[target]
	return e, ok
}

func (c *Checkpoint) save() error {
	c.mu.RLock()
	data, err := json.MarshalIndent(c.entries, "", "  ")
	c.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}
