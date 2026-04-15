// Package throttle provides a simple per-key throttle that prevents
// the same alert key from firing more than once within a cooldown window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last fire time for each key and suppresses
// repeated firings within the cooldown duration.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New creates a new Throttle with the given cooldown duration.
// Returns an error if cooldown is not positive.
func New(cooldown time.Duration) (*Throttle, error) {
	if cooldown <= 0 {
		return nil, ErrInvalidCooldown
	}
	return &Throttle{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}, nil
}

// Allow returns true if the key has not fired within the cooldown window,
// and records the current time as the last fire time for the key.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded fire time for a specific key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Prune removes all entries whose last fire time is older than the cooldown,
// freeing memory for keys that are no longer active.
func (t *Throttle) Prune() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	for k, last := range t.last {
		if now.Sub(last) >= t.cooldown {
			delete(t.last, k)
		}
	}
}
