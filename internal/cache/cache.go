// Package cache provides a simple TTL-based in-memory cache for scan results.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached scan result with an expiry timestamp.
type Entry struct {
	Ports     []int
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe TTL store keyed by host string.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New returns a Cache with the given TTL.
func New(ttl time.Duration) (*Cache, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}, nil
}

// Set stores ports for the given host key, replacing any existing entry.
func (c *Cache) Set(host string, ports []int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.entries[host] = Entry{
		Ports:     ports,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get returns the cached entry for host and whether it was found and still valid.
func (c *Cache) Get(host string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[host]
	if !ok || e.IsExpired() {
		return Entry{}, false
	}
	return e, true
}

// Invalidate removes the entry for host if it exists.
func (c *Cache) Invalidate(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, host)
}

// Prune removes all expired entries from the cache.
func (c *Cache) Prune() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	removed := 0
	for k, e := range c.entries {
		if e.IsExpired() {
			delete(c.entries, k)
			removed++
		}
	}
	return removed
}
