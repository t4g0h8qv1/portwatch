// Package portsnap captures and compares point-in-time port snapshots.
package portsnap

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds the open ports observed on a host at a specific time.
type Snapshot struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	CapturedAt time.Time `json:"captured_at"`
}

// Take creates a new Snapshot for the given host and ports.
func Take(host string, ports []int) Snapshot {
	return Snapshot{
		Host:      host,
		Ports:     ports,
		CapturedAt: time.Now().UTC(),
	}
}

// Save writes the snapshot to a JSON file at path.
func Save(path string, s Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("portsnap: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from a JSON file at path.
func Load(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, fmt.Errorf("portsnap: file not found: %s", path)
		}
		return Snapshot{}, fmt.Errorf("portsnap: open %s: %w", path, err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Snapshot{}, fmt.Errorf("portsnap: decode: %w", err)
	}
	return s, nil
}

// Diff returns ports opened and closed between old and new snapshots.
func Diff(old, new Snapshot) (opened, closed []int) {
	oldSet := toSet(old.Ports)
	newSet := toSet(new.Ports)
	for p := range newSet {
		if !oldSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range oldSet {
		if !newSet[p] {
			closed = append(closed, p)
		}
	}
	return
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
