package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Baseline represents a saved snapshot of expected open ports for a host.
type Baseline struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New creates a new Baseline for the given host and port list.
func New(host string, ports []int) *Baseline {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)
	now := time.Now().UTC()
	return &Baseline{
		Host:      host,
		Ports:     sorted,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Save writes the baseline to a JSON file at the given path.
func (b *Baseline) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// Load reads a baseline from a JSON file at the given path.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: open file: %w", err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode: %w", err)
	}
	return &b, nil
}

// Diff compares the baseline ports against a current set of open ports.
// It returns ports that are new (not in baseline) and ports that are missing
// (in baseline but no longer open).
func (b *Baseline) Diff(current []int) (newPorts, missingPorts []int) {
	baseSet := make(map[int]struct{}, len(b.Ports))
	for _, p := range b.Ports {
		baseSet[p] = struct{}{}
	}
	currentSet := make(map[int]struct{}, len(current))
	for _, p := range current {
		currentSet[p] = struct{}{}
	}
	for _, p := range current {
		if _, ok := baseSet[p]; !ok {
			newPorts = append(newPorts, p)
		}
	}
	for _, p := range b.Ports {
		if _, ok := currentSet[p]; !ok {
			missingPorts = append(missingPorts, p)
		}
	}
	sort.Ints(newPorts)
	sort.Ints(missingPorts)
	return
}
