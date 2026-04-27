package portwatch

import (
	"errors"
	"sync"
	"time"
)

// ErrEmptyQuarantineTarget is returned when an empty target is provided.
var ErrEmptyQuarantineTarget = errors.New("portwatch: quarantine target must not be empty")

// ErrInvalidQuarantineDuration is returned when a zero or negative duration is provided.
var ErrInvalidQuarantineDuration = errors.New("portwatch: quarantine duration must be positive")

// quarantineEntry holds the expiry time for a quarantined target.
type quarantineEntry struct {
	until time.Time
}

// QuarantineManager tracks targets that are temporarily quarantined and
// should be skipped during scanning.
type QuarantineManager struct {
	mu      sync.RWMutex
	entries map[string]quarantineEntry
	now     func() time.Time
}

// NewQuarantineManager returns a new QuarantineManager.
func NewQuarantineManager() *QuarantineManager {
	return &QuarantineManager{
		entries: make(map[string]quarantineEntry),
		now:     time.Now,
	}
}

// Quarantine places target in quarantine for the given duration.
func (q *QuarantineManager) Quarantine(target string, d time.Duration) error {
	if target == "" {
		return ErrEmptyQuarantineTarget
	}
	if d <= 0 {
		return ErrInvalidQuarantineDuration
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries[target] = quarantineEntry{until: q.now().Add(d)}
	return nil
}

// IsQuarantined reports whether target is currently quarantined.
func (q *QuarantineManager) IsQuarantined(target string) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	e, ok := q.entries[target]
	if !ok {
		return false
	}
	return q.now().Before(e.until)
}

// Release removes target from quarantine immediately.
func (q *QuarantineManager) Release(target string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.entries, target)
}

// Prune removes all expired quarantine entries.
func (q *QuarantineManager) Prune() {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := q.now()
	for t, e := range q.entries {
		if !now.Before(e.until) {
			delete(q.entries, t)
		}
	}
}

// Count returns the number of currently active quarantine entries.
func (q *QuarantineManager) Count() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	now := q.now()
	count := 0
	for _, e := range q.entries {
		if now.Before(e.until) {
			count++
		}
	}
	return count
}
