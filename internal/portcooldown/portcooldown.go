// Package portcooldown tracks how recently a port was first seen open
// and enforces a minimum observation window before alerts are raised.
package portcooldown

import (
	"sync"
	"time"
)

// Tracker records the first-seen timestamp for each port on a given host
// and reports whether the port has been open long enough to alert.
type Tracker struct {
	mu       sync.Mutex
	window   time.Duration
	firstSeen map[string]map[int]time.Time
	now      func() time.Time
}

// New returns a Tracker that suppresses alerts until a port has been
// continuously observed for at least window.
func New(window time.Duration) (*Tracker, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Tracker{
		window:    window,
		firstSeen: make(map[string]map[int]time.Time),
		now:       time.Now,
	}, nil
}

// Observe records that port is open on host. It is idempotent: calling it
// multiple times does not reset the first-seen clock.
func (t *Tracker) Observe(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.firstSeen[host] == nil {
		t.firstSeen[host] = make(map[int]time.Time)
	}
	if _, ok := t.firstSeen[host][port]; !ok {
		t.firstSeen[host][port] = t.now()
	}
}

// Forget removes the first-seen record for port on host (e.g. when the port
// closes so the clock resets if it reopens later).
func (t *Tracker) Forget(host string, port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.firstSeen[host] != nil {
		delete(t.firstSeen[host], port)
	}
}

// Ready reports whether port on host has been observed for at least the
// configured window and should therefore trigger an alert.
func (t *Tracker) Ready(host string, port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	ports, ok := t.firstSeen[host]
	if !ok {
		return false
	}
	first, ok := ports[port]
	if !ok {
		return false
	}
	return t.now().Sub(first) >= t.window
}
