package portwatch

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"sync"
	"time"
)

// DigestEntry holds a computed digest for a target's port set at a point in time.
type DigestEntry struct {
	Target    string
	Digest    string
	ComputedAt time.Time
	PortCount int
}

// DigestManager tracks SHA-256 digests of observed port sets per target.
// A changed digest indicates the port set has changed since last observation.
type DigestManager struct {
	mu      sync.RWMutex
	entries map[string]DigestEntry
}

// NewDigestManager returns an initialised DigestManager.
func NewDigestManager() *DigestManager {
	return &DigestManager{
		entries: make(map[string]DigestEntry),
	}
}

// Record computes and stores the digest for target's current port set.
// Returns the digest string and whether the digest changed since last call.
func (d *DigestManager) Record(target string, ports []int, now time.Time) (digest string, changed bool) {
	digest = computeDigest(ports)

	d.mu.Lock()
	defer d.mu.Unlock()

	prev, exists := d.entries[target]
	changed = !exists || prev.Digest != digest

	d.entries[target] = DigestEntry{
		Target:    target,
		Digest:    digest,
		ComputedAt: now,
		PortCount: len(ports),
	}
	return digest, changed
}

// Get returns the most recent DigestEntry for target, or false if not found.
func (d *DigestManager) Get(target string) (DigestEntry, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	e, ok := d.entries[target]
	return e, ok
}

// Targets returns all tracked target names in sorted order.
func (d *DigestManager) Targets() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]string, 0, len(d.entries))
	for t := range d.entries {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// computeDigest returns a deterministic SHA-256 hex digest of a sorted port list.
func computeDigest(ports []int) string {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	h := sha256.New()
	for _, p := range sorted {
		fmt.Fprintf(h, "%d\n", p)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
