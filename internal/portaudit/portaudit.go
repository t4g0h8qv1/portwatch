// Package portaudit correlates scan results with a known baseline and
// produces a structured audit record suitable for reporting and alerting.
package portaudit

import (
	"time"

	"github.com/user/portwatch/internal/portstate"
	"github.com/user/portwatch/internal/severity"
)

// Record holds the outcome of a single audit run.
type Record struct {
	Host      string
	ScannedAt time.Time
	New       []int
	Gone      []int
	Stable    []int
	Severity  severity.Level
}

// Auditor compares live scan results against a baseline.
type Auditor struct {
	classifier *severity.Classifier
}

// New returns an Auditor using the provided severity classifier.
func New(c *severity.Classifier) *Auditor {
	return &Auditor{classifier: c}
}

// Run performs the audit and returns a Record.
func (a *Auditor) Run(host string, baseline, current []int) Record {
	diff := portstate.Compare(baseline, current)

	rec := Record{
		Host:      host,
		ScannedAt: time.Now().UTC(),
		New:       diff.New,
		Gone:      diff.Gone,
		Stable:    diff.Stable,
	}

	rec.Severity = a.highest(diff.New)
	return rec
}

func (a *Auditor) highest(ports []int) severity.Level {
	level := severity.Info
	for _, p := range ports {
		if l := a.classifier.Classify(p); l > level {
			level = l
		}
	}
	return level
}

// HasChanges reports whether the record contains any port changes.
func (r Record) HasChanges() bool {
	return len(r.New) > 0 || len(r.Gone) > 0
}
