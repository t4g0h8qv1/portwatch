// Package portacl implements access control lists for ports,
// allowing or denying traffic based on port number and direction.
package portacl

import "fmt"

// Direction represents inbound or outbound traffic.
type Direction string

const (
	Inbound  Direction = "inbound"
	Outbound Direction = "outbound"
)

// Action represents the ACL decision.
type Action string

const (
	Allow Action = "allow"
	Deny  Action = "deny"
)

// Rule is a single ACL entry.
type Rule struct {
	Port      int
	Direction Direction
	Action    Action
}

// ACL holds an ordered list of rules.
type ACL struct {
	rules []Rule
}

// New returns an empty ACL.
func New() *ACL {
	return &ACL{}
}

// Add appends a rule to the ACL.
func (a *ACL) Add(r Rule) error {
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("portacl: invalid port %d", r.Port)
	}
	if r.Direction != Inbound && r.Direction != Outbound {
		return fmt.Errorf("portacl: invalid direction %q", r.Direction)
	}
	if r.Action != Allow && r.Action != Deny {
		return fmt.Errorf("portacl: invalid action %q", r.Action)
	}
	a.rules = append(a.rules, r)
	return nil
}

// Evaluate returns the Action for the given port and direction.
// Rules are evaluated in order; the first match wins.
// If no rule matches, Allow is returned.
func (a *ACL) Evaluate(port int, dir Direction) Action {
	for _, r := range a.rules {
		if r.Port == port && r.Direction == dir {
			return r.Action
		}
	}
	return Allow
}

// EvaluateAll returns a map of port→Action for all given ports.
func (a *ACL) EvaluateAll(ports []int, dir Direction) map[int]Action {
	out := make(map[int]Action, len(ports))
	for _, p := range ports {
		out[p] = a.Evaluate(p, dir)
	}
	return out
}
