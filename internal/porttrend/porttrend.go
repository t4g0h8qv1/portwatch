// Package porttrend tracks how frequently ports appear across scan history.
package porttrend

import (
	"sort"

	"github.com/example/portwatch/internal/history"
)

// Entry holds trend data for a single port.
type Entry struct {
	Port      int
	SeenCount int
	Frequency float64 // fraction of scans in which port was open
}

// Analyzer computes port frequency trends from scan history.
type Analyzer struct {
	records []history.Record
}

// New returns an Analyzer loaded with the provided records.
func New(records []history.Record) *Analyzer {
	return &Analyzer{records: records}
}

// Analyze returns trend entries sorted by frequency descending.
func (a *Analyzer) Analyze() []Entry {
	if len(a.records) == 0 {
		return nil
	}

	counts := make(map[int]int)
	for _, r := range a.records {
		for _, p := range r.OpenPorts {
			counts[p]++
		}
	}

	total := float64(len(a.records))
	entries := make([]Entry, 0, len(counts))
	for port, seen := range counts {
		entries = append(entries, Entry{
			Port:      port,
			SeenCount: seen,
			Frequency: float64(seen) / total,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Frequency != entries[j].Frequency {
			return entries[i].Frequency > entries[j].Frequency
		}
		return entries[i].Port < entries[j].Port
	})

	return entries
}

// Unstable returns ports that appear in fewer than threshold fraction of scans.
func (a *Analyzer) Unstable(threshold float64) []Entry {
	all := a.Analyze()
	var out []Entry
	for _, e := range all {
		if e.Frequency < threshold {
			out = append(out, e)
		}
	}
	return out
}
