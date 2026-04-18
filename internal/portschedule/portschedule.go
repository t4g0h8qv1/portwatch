// Package portschedule manages per-host scan scheduling with configurable intervals.
package portschedule

import (
	"errors"
	"sync"
	"time"
)

// Entry holds scheduling metadata for a single host.
type Entry struct {
	Host     string
	Interval time.Duration
	LastRun  time.Time
	NextRun  time.Time
}

// Scheduler tracks when each host should next be scanned.
type Scheduler struct {
	mu      sync.Mutex
	entries map[string]*Entry
	default_ time.Duration
}

// New creates a Scheduler with the given default interval.
func New(defaultInterval time.Duration) (*Scheduler, error) {
	if defaultInterval <= 0 {
		return nil, errors.New("portschedule: interval must be positive")
	}
	return &Scheduler{
		entries:  make(map[string]*Entry),
		default_: defaultInterval,
	}, nil
}

// Register adds or updates a host with an optional override interval.
// Pass 0 to use the default interval.
func (s *Scheduler) Register(host string, interval time.Duration) error {
	if host == "" {
		return errors.New("portschedule: host must not be empty")
	}
	if interval <= 0 {
		interval = s.default_
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.entries[host] = &Entry{
		Host:     host,
		Interval: interval,
		LastRun:  time.Time{},
		NextRun:  now,
	}
	return nil
}

// Due returns all hosts whose NextRun is at or before now.
func (s *Scheduler) Due() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	var hosts []string
	for _, e := range s.entries {
		if !now.Before(e.NextRun) {
			hosts = append(hosts, e.Host)
		}
	}
	return hosts
}

// Advance marks a host as having just been scanned, updating LastRun and NextRun.
func (s *Scheduler) Advance(host string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[host]
	if !ok {
		return errors.New("portschedule: host not registered: " + host)
	}
	now := time.Now()
	e.LastRun = now
	e.NextRun = now.Add(e.Interval)
	return nil
}

// Get returns the Entry for a host, or false if not registered.
func (s *Scheduler) Get(host string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[host]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}
