// Package ratelimit provides a simple token-bucket rate limiter to prevent
// flooding notifiers or scanners when portwatch is run on a tight schedule.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter is a thread-safe token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to max events per interval.
// For example, New(5, time.Minute) allows 5 events per minute.
func New(max int, interval time.Duration) (*Limiter, error) {
	if max <= 0 {
		return nil, fmt.Errorf("ratelimit: max must be positive, got %d", max)
	}
	if interval <= 0 {
		return nil, fmt.Errorf("ratelimit: interval must be positive, got %s", interval)
	}
	return &Limiter{
		tokens:   float64(max),
		max:      float64(max),
		rate:     float64(max) / interval.Seconds(),
		lastTick: time.Now(),
		clock:    time.Now,
	}, nil
}

// Allow reports whether an event may proceed. It refills tokens based on
// elapsed time since the last call before checking.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}
	return false
}

// Tokens returns the current number of available tokens (for inspection/testing).
func (l *Limiter) Tokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
