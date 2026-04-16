// Package portlabel maps port numbers to human-readable service names.
package portlabel

import "fmt"

// Label holds a port number and its resolved service name.
type Label struct {
	Port    int
	Service string
}

// Labeler resolves port numbers to service names.
type Labeler struct {
	custom map[int]string
}

// New returns a Labeler seeded with well-known IANA service names.
func New(custom map[int]string) *Labeler {
	l := &Labeler{custom: make(map[int]string)}
	for k, v := range custom {
		l.custom[k] = v
	}
	return l
}

// Resolve returns the service name for port. Custom entries take
// precedence over built-in names. Unknown ports return "unknown".
func (l *Labeler) Resolve(port int) string {
	if name, ok := l.custom[port]; ok {
		return name
	}
	if name, ok := builtIn[port]; ok {
		return name
	}
	return "unknown"
}

// Label returns a Label for the given port.
func (l *Labeler) Label(port int) Label {
	return Label{Port: port, Service: l.Resolve(port)}
}

// String formats a Label as "port/service".
func (lb Label) String() string {
	return fmt.Sprintf("%d/%s", lb.Port, lb.Service)
}

// LabelAll resolves a slice of ports into Labels.
func (l *Labeler) LabelAll(ports []int) []Label {
	out := make([]Label, len(ports))
	for i, p := range ports {
		out[i] = l.Label(p)
	}
	return out
}
