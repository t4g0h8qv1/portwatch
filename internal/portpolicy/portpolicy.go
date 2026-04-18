// Package portpolicy evaluates open ports against a defined policy,
// classifying each port as allowed, denied, or unreviewed.
package portpolicy

import "fmt"

// Status represents the policy verdict for a port.
type Status string

const (
	Allowed    Status = "allowed"
	Denied     Status = "denied"
	Unreviewed Status = "unreviewed"
)

// Result holds the policy evaluation result for a single port.
type Result struct {
	Port   int
	Status Status
	Reason string
}

// Policy holds sets of explicitly allowed and denied ports.
type Policy struct {
	allowed map[int]struct{}
	denied  map[int]struct{}
}

// New creates a Policy from allowed and denied port lists.
// Denied takes precedence over allowed.
func New(allowed, denied []int) *Policy {
	p := &Policy{
		allowed: make(map[int]struct{}, len(allowed)),
		denied:  make(map[int]struct{}, len(denied)),
	}
	for _, port := range allowed {
		p.allowed[port] = struct{}{}
	}
	for _, port := range denied {
		p.denied[port] = struct{}{}
	}
	return p
}

// Evaluate checks a single port against the policy.
func (p *Policy) Evaluate(port int) Result {
	if _, ok := p.denied[port]; ok {
		return Result{Port: port, Status: Denied, Reason: fmt.Sprintf("port %d is explicitly denied", port)}
	}
	if _, ok := p.allowed[port]; ok {
		return Result{Port: port, Status: Allowed, Reason: fmt.Sprintf("port %d is explicitly allowed", port)}
	}
	return Result{Port: port, Status: Unreviewed, Reason: fmt.Sprintf("port %d has no policy entry", port)}
}

// EvaluateAll evaluates a slice of ports and returns all results.
func (p *Policy) EvaluateAll(ports []int) []Result {
	results := make([]Result, len(ports))
	for i, port := range ports {
		results[i] = p.Evaluate(port)
	}
	return results
}

// Violations returns only the denied results from EvaluateAll.
func (p *Policy) Violations(ports []int) []Result {
	var out []Result
	for _, r := range p.EvaluateAll(ports) {
		if r.Status == Denied {
			out = append(out, r)
		}
	}
	return out
}
