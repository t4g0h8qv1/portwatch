// Package portrank assigns a risk rank to open ports based on exposure
// history and known service classifications.
package portrank

import (
	"sort"
	"sync"
)

// Rank represents the risk level of an open port.
type Rank int

const (
	RankLow Rank = iota
	RankMedium
	RankHigh
	RankCritical
)

func (r Rank) String() string {
	switch r {
	case RankLow:
		return "low"
	case RankMedium:
		return "medium"
	case RankHigh:
		return "high"
	case RankCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Result holds a port and its computed rank.
type Result struct {
	Port int
	Rank Rank
}

// Ranker scores ports by risk.
type Ranker struct {
	mu       sync.RWMutex
	override map[int]Rank
}

// New returns a Ranker with built-in defaults.
func New() *Ranker {
	return &Ranker{override: make(map[int]Rank)}
}

// SetOverride assigns a fixed rank to a specific port.
func (r *Ranker) SetOverride(port int, rank Rank) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.override[port] = rank
}

// Score returns the risk rank for a single port.
func (r *Ranker) Score(port int) Rank {
	r.mu.RLock()
	if rank, ok := r.override[port]; ok {
		r.mu.RUnlock()
		return rank
	}
	r.mu.RUnlock()
	return defaultRank(port)
}

// RankAll scores every port in the slice and returns results sorted
// from highest to lowest risk.
func (r *Ranker) RankAll(ports []int) []Result {
	results := make([]Result, len(ports))
	for i, p := range ports {
		results[i] = Result{Port: p, Rank: r.Score(p)}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Rank > results[j].Rank
	})
	return results
}

func defaultRank(port int) Rank {
	switch {
	case port == 22 || port == 23 || port == 3389:
		return RankCritical
	case port == 21 || port == 25 || port == 445 || port == 1433 || port == 3306:
		return RankHigh
	case port < 1024:
		return RankMedium
	default:
		return RankLow
	}
}
