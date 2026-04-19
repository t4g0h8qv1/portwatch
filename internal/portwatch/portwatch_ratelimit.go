package portwatch

import (
	"errors"
	"sync"
	"time"
)

// ScanRateLimiter enforces a maximum number of scans per target within a
// rolling time window.
type ScanRateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxScans int
	records  map[string][]time.Time
}

// NewScanRateLimiter creates a ScanRateLimiter with the given window and max
// scan count. Both must be positive.
func NewScanRateLimiter(window time.Duration, maxScans int) (*ScanRateLimiter, error) {
	if window <= 0 {
		return nil, errors.New("portwatch: window must be positive")
	}
	if maxScans <= 0 {
		return nil, errors.New("portwatch: maxScans must be positive")
	}
	return &ScanRateLimiter{
		window:   window,
		maxScans: maxScans,
		records:  make(map[string][]time.Time),
	}, nil
}

// Allow returns true if the target may be scanned now. It records the attempt
// when permitted.
func (r *ScanRateLimiter) Allow(target string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.prune(target, now)
	if len(r.records[target]) >= r.maxScans {
		return false
	}
	r.records[target] = append(r.records[target], now)
	return true
}

// Remaining returns how many scans are still permitted for target in the
// current window.
func (r *ScanRateLimiter) Remaining(target string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prune(target, time.Now())
	rem := r.maxScans - len(r.records[target])
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the scan history for target.
func (r *ScanRateLimiter) Reset(target string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, target)
}

func (r *ScanRateLimiter) prune(target string, now time.Time) {
	cutoff := now.Add(-r.window)
	times := r.records[target]
	i := 0
	for i < len(times) && times[i].Before(cutoff) {
		i++
	}
	r.records[target] = times[i:]
}
