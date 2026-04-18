// Package portbudget enforces a maximum number of concurrently open ports
// per host, emitting a violation when the observed count exceeds the budget.
package portbudget

import (
	"errors"
	"fmt"
	"sync"
)

// ErrBudgetExceeded is returned when the open-port count exceeds the budget.
var ErrBudgetExceeded = errors.New("port budget exceeded")

// Violation describes a single budget breach.
type Violation struct {
	Host    string
	Max     int
	Actual  int
	Ports   []int
}

func (v Violation) Error() string {
	return fmt.Sprintf("host %s has %d open ports (budget: %d)", v.Host, v.Actual, v.Max)
}

// Budget holds per-host maximums.
type Budget struct {
	mu       sync.RWMutex
	defaults int
	override map[string]int
}

// New creates a Budget with the given default maximum.
// max must be >= 1.
func New(max int) (*Budget, error) {
	if max < 1 {
		return nil, errors.New("portbudget: max must be >= 1")
	}
	return &Budget{defaults: max, override: make(map[string]int)}, nil
}

// SetHost sets a per-host maximum, overriding the default.
func (b *Budget) SetHost(host string, max int) error {
	if max < 1 {
		return errors.New("portbudget: max must be >= 1")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.override[host] = max
	return nil
}

// Check returns a Violation if len(ports) exceeds the budget for host,
// or nil if within budget.
func (b *Budget) Check(host string, ports []int) *Violation {
	b.mu.RLock()
	max, ok := b.override[host]
	if !ok {
		max = b.defaults
	}
	b.mu.RUnlock()

	if len(ports) <= max {
		return nil
	}
	copy := make([]int, len(ports))
	_ = copy
	return &Violation{Host: host, Max: max, Actual: len(ports), Ports: ports}
}
