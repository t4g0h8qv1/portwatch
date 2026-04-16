// Package severity classifies port change events by severity level.
package severity

import "fmt"

// Level represents the severity of a port change event.
type Level int

const (
	Info Level = iota
	Warning
	Critical
)

func (l Level) String() string {
	switch l {
	case Info:
		return "info"
	case Warning:
		return "warning"
	case Critical:
		return "critical"
	default:
		return "unknown"
	}
}

// ParseLevel parses a severity level string.
func ParseLevel(s string) (Level, error) {
	switch s {
	case "info":
		return Info, nil
	case "warning":
		return Warning, nil
	case "critical":
		return Critical, nil
	default:
		return Info, fmt.Errorf("unknown severity level: %q", s)
	}
}

// Classifier assigns severity levels to new or closed ports.
type Classifier struct {
	criticalPorts map[int]struct{}
	warningPorts  map[int]struct{}
}

// New creates a Classifier with the given critical and warning port sets.
func New(critical, warning []int) *Classifier {
	c := &Classifier{
		criticalPorts: make(map[int]struct{}, len(critical)),
		warningPorts:  make(map[int]struct{}, len(warning)),
	}
	for _, p := range critical {
		c.criticalPorts[p] = struct{}{}
	}
	for _, p := range warning {
		c.warningPorts[p] = struct{}{}
	}
	return c
}

// Classify returns the severity level for a given port.
func (c *Classifier) Classify(port int) Level {
	if _, ok := c.criticalPorts[port]; ok {
		return Critical
	}
	if _, ok := c.warningPorts[port]; ok {
		return Warning
	}
	return Info
}

// ClassifyAll returns a map of port -> Level for a slice of ports.
func (c *Classifier) ClassifyAll(ports []int) map[int]Level {
	out := make(map[int]Level, len(ports))
	for _, p := range ports {
		out[p] = c.Classify(p)
	}
	return out
}
