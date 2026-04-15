// Package history provides persistent scan history tracking for portwatch.
// It records timestamped scan results so that trends and changes over time
// can be reviewed and audited.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single recorded scan result.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
}

// History holds an ordered list of scan entries.
type History struct {
	Entries []Entry `json:"entries"`
}

// Record appends a new entry to the history and persists it to path.
func Record(path string, host string, ports []int) error {
	h, err := Load(path)
	if err != nil {
		return fmt.Errorf("history: load: %w", err)
	}
	h.Entries = append(h.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Ports:     ports,
	})
	return save(path, h)
}

// Load reads history from path. If the file does not exist an empty History
// is returned without error.
func Load(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &History{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read %s: %w", path, err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("history: unmarshal: %w", err)
	}
	return &h, nil
}

// Last returns the most recent Entry, and false if no entries exist.
func (h *History) Last() (Entry, bool) {
	if len(h.Entries) == 0 {
		return Entry{}, false
	}
	return h.Entries[len(h.Entries)-1], true
}

func save(path string, h *History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %s: %w", path, err)
	}
	return nil
}
