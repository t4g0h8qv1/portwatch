// Package hostmeta collects basic metadata about a scanned host.
package hostmeta

import (
	"fmt"
	"net"
	"time"
)

// Meta holds metadata resolved for a target host.
type Meta struct {
	Input     string    `json:"input"`
	Resolved  string    `json:"resolved"`
	Hostnames []string  `json:"hostnames"`
	ScannedAt time.Time `json:"scanned_at"`
}

// Collector resolves and stores host metadata.
type Collector struct {
	lookupAddr func(string) ([]string, error)
	lookupHost func(string) ([]string, error)
	now        func() time.Time
}

// New returns a Collector using real DNS lookups.
func New() *Collector {
	return &Collector{
		lookupAddr: net.LookupAddr,
		lookupHost: net.LookupHost,
		now:        time.Now,
	}
}

// Collect resolves the given target and returns a Meta.
func (c *Collector) Collect(target string) (Meta, error) {
	m := Meta{
		Input:     target,
		ScannedAt: c.now(),
	}

	// If already an IP, skip forward lookup.
	if ip := net.ParseIP(target); ip != nil {
		m.Resolved = ip.String()
	} else {
		addrs, err := c.lookupHost(target)
		if err != nil {
			return Meta{}, fmt.Errorf("hostmeta: lookup %q: %w", target, err)
		}
		if len(addrs) == 0 {
			return Meta{}, fmt.Errorf("hostmeta: no addresses for %q", target)
		}
		m.Resolved = addrs[0]
	}

	// Reverse lookup for hostnames.
	names, _ := c.lookupAddr(m.Resolved)
	m.Hostnames = names

	return m, nil
}
