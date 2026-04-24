package portwatch

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// AuditEvent records a single scan audit entry for a target.
type AuditEvent struct {
	Target    string
	Timestamp time.Time
	Ports     []int
	Opened    []int
	Closed    []int
	Error     error
}

// AuditLog maintains an ordered, bounded in-memory audit trail of scan events.
type AuditLog struct {
	mu     sync.Mutex
	events []AuditEvent
	maxLen int
}

// NewAuditLog creates an AuditLog that retains at most maxLen events.
// Returns an error if maxLen is less than 1.
func NewAuditLog(maxLen int) (*AuditLog, error) {
	if maxLen < 1 {
		return nil, fmt.Errorf("portwatch: audit log maxLen must be >= 1, got %d", maxLen)
	}
	return &AuditLog{maxLen: maxLen}, nil
}

// Record appends an AuditEvent, evicting the oldest entry when the log is full.
func (a *AuditLog) Record(ev AuditEvent) {
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now()
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.events) >= a.maxLen {
		a.events = a.events[1:]
	}
	a.events = append(a.events, ev)
}

// All returns a shallow copy of all recorded events, oldest first.
func (a *AuditLog) All() []AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]AuditEvent, len(a.events))
	copy(out, a.events)
	return out
}

// Last returns the most recent event for the given target, and whether one was found.
func (a *AuditLog) Last(target string) (AuditEvent, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := len(a.events) - 1; i >= 0; i-- {
		if a.events[i].Target == target {
			return a.events[i], true
		}
	}
	return AuditEvent{}, false
}

// WriteAuditTable writes a human-readable audit table to w.
func WriteAuditTable(w io.Writer, events []AuditEvent) {
	fmt.Fprintf(w, "%-20s %-8s %-8s %s\n", "TARGET", "OPENED", "CLOSED", "TIMESTAMP")
	for _, ev := range events {
		errStr := ""
		if ev.Error != nil {
			errStr = " err=" + ev.Error.Error()
		}
		fmt.Fprintf(w, "%-20s %-8d %-8d %s%s\n",
			ev.Target,
			len(ev.Opened),
			len(ev.Closed),
			ev.Timestamp.Format(time.RFC3339),
			errStr,
		)
	}
}
