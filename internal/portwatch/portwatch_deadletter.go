package portwatch

import (
	"sync"
	"time"
)

// DeadLetterEntry holds a scan event that could not be processed.
type DeadLetterEntry struct {
	Target    string
	Err       error
	OccurredAt time.Time
	Attempts  int
}

// DeadLetterQueue stores unprocessable scan events for later inspection.
type DeadLetterQueue struct {
	mu      sync.Mutex
	entries []DeadLetterEntry
	maxSize int
}

// NewDeadLetterQueue creates a DeadLetterQueue with the given capacity.
// maxSize must be >= 1.
func NewDeadLetterQueue(maxSize int) (*DeadLetterQueue, error) {
	if maxSize < 1 {
		return nil, ErrInvalidDeadLetterSize
	}
	return &DeadLetterQueue{maxSize: maxSize}, nil
}

// Push adds an entry to the queue. If the queue is full the oldest entry is
// evicted to make room.
func (q *DeadLetterQueue) Push(target string, err error, attempts int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	entry := DeadLetterEntry{
		Target:     target,
		Err:        err,
		OccurredAt: time.Now(),
		Attempts:   attempts,
	}
	if len(q.entries) >= q.maxSize {
		q.entries = q.entries[1:]
	}
	q.entries = append(q.entries, entry)
}

// All returns a copy of all current entries.
func (q *DeadLetterQueue) All() []DeadLetterEntry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]DeadLetterEntry, len(q.entries))
	copy(out, q.entries)
	return out
}

// Len returns the number of entries currently in the queue.
func (q *DeadLetterQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.entries)
}

// Clear removes all entries from the queue.
func (q *DeadLetterQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries = q.entries[:0]
}
