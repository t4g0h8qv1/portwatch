// Package alert evaluates scan results against a baseline and notifies on changes.
package alert

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Result holds the diff between a baseline and a current scan.
type Result struct {
	NewPorts  []int
	GonePorts []int
	Changed   bool
}

// Notifier sends alerts when port changes are detected.
type Notifier interface {
	Notify(result Result) error
}

// StdoutNotifier writes alerts to stdout (or a custom writer).
type StdoutNotifier struct {
	out io.Writer
}

// NewStdoutNotifier creates a StdoutNotifier writing to out.
// If out is nil, os.Stdout is used.
func NewStdoutNotifier(out io.Writer) *StdoutNotifier {
	if out == nil {
		out = os.Stdout
	}
	return &StdoutNotifier{out: out}
}

// Notify prints port change information to the configured writer.
func (s *StdoutNotifier) Notify(result Result) error {
	if !result.Changed {
		return nil
	}
	if len(result.NewPorts) > 0 {
		fmt.Fprintf(s.out, "ALERT: new ports detected: %v\n", result.NewPorts)
	}
	if len(result.GonePorts) > 0 {
		fmt.Fprintf(s.out, "ALERT: ports no longer open: %v\n", result.GonePorts)
	}
	return nil
}

// Evaluate compares a current set of open ports against a baseline set.
// It returns a Result describing what changed.
func Evaluate(baseline, current []int) Result {
	baseSet := toSet(baseline)
	currSet := toSet(current)

	var newPorts, gonePorts []int

	for p := range currSet {
		if !baseSet[p] {
			newPorts = append(newPorts, p)
		}
	}
	for p := range baseSet {
		if !currSet[p] {
			gonePorts = append(gonePorts, p)
		}
	}

	sort.Ints(newPorts)
	sort.Ints(gonePorts)

	return Result{
		NewPorts:  newPorts,
		GonePorts: gonePorts,
		Changed:   len(newPorts) > 0 || len(gonePorts) > 0,
	}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
