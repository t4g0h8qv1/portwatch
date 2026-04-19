package portwatch

import (
	"errors"
	"sync"
	"time"
)

// ScanLimiter enforces a minimum interval between scans for a given target.
type ScanLimiter struct {
	mu       sync.Mutex
	last     map[string]time.Time
	minGap   time.Duration
	nowFunc  func() time.Time
}

// ErrTooSoon is returned when a scan is attempted before the minimum gap has elapsed.
var ErrTooSoon = errors.New("scan attempted too soon after previous scan")

// NewScanLimiter creates a ScanLimiter with the given minimum gap between scans.
func NewScanLimiter(minGap time.Duration) (*ScanLimiter, error) {
	if minGap <= 0 {
		return nil, errors.New("minGap must be positive")
	}
	return &ScanLimiter{
		last:    make(map[string]time.Time),
		minGap:  minGap,
		nowFunc: time.Now,
	}, nil
}

// Allow returns nil if a scan for target is permitted, or ErrTooSoon if not.
func (l *ScanLimiter) Allow(target string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.nowFunc()
	if t, ok := l.last[target]; ok {
		if now.Sub(t) < l.minGap {
			return ErrTooSoon
		}
	}
	l.last[target] = now
	return nil
}

// Reset clears the recorded scan time for target.
func (l *ScanLimiter) Reset(target string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, target)
}

// LastScan returns the time of the most recent allowed scan for target, and
// whether a record exists.
func (l *ScanLimiter) LastScan(target string) (time.Time, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.last[target]
	return t, ok
}
