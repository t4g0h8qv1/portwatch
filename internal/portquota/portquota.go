// Package portquota enforces a maximum number of open ports allowed on a host.
// If the observed port count exceeds the configured limit an error is returned
// so callers can trigger an alert.
package portquota

import (
	"errors"
	"fmt"
)

// ErrQuotaExceeded is returned when the number of open ports exceeds the limit.
var ErrQuotaExceeded = errors.New("port quota exceeded")

// Quota holds the configuration for a port count limit.
type Quota struct {
	max int
}

// New creates a Quota with the given maximum allowed open ports.
// max must be greater than zero.
func New(max int) (*Quota, error) {
	if max <= 0 {
		return nil, fmt.Errorf("portquota: max must be greater than zero, got %d", max)
	}
	return &Quota{max: max}, nil
}

// Check returns nil when the number of open ports is within the allowed limit.
// It returns ErrQuotaExceeded (with a descriptive message) when the count
// exceeds the maximum.
func (q *Quota) Check(openPorts []int) error {
	count := len(openPorts)
	if count > q.max {
		return fmt.Errorf("%w: %d open ports found, limit is %d", ErrQuotaExceeded, count, q.max)
	}
	return nil
}

// Max returns the configured maximum number of open ports.
func (q *Quota) Max() int {
	return q.max
}

// Remaining returns how many additional ports may be open before the quota is
// breached. A negative value indicates the quota has already been exceeded.
func (q *Quota) Remaining(openPorts []int) int {
	return q.max - len(openPorts)
}
