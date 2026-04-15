// Package resolve provides hostname-to-IP resolution with caching support.
package resolve

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Resolver resolves hostnames to IP addresses.
type Resolver struct {
	cache   map[string]entry
	ttl     time.Duration
	lookup  func(ctx context.Context, host string) ([]string, error)
}

type entry struct {
	addrs   []string
	expires time.Time
}

// New returns a Resolver with the given cache TTL.
// Pass 0 to disable caching.
func New(ttl time.Duration) *Resolver {
	return &Resolver{
		cache:  make(map[string]entry),
		ttl:    ttl,
		lookup: defaultLookup,
	}
}

// Resolve returns the IP addresses for the given host.
// If host is already an IP address it is returned as-is.
// Results are cached for the configured TTL.
func (r *Resolver) Resolve(ctx context.Context, host string) ([]string, error) {
	if net.ParseIP(host) != nil {
		return []string{host}, nil
	}

	if r.ttl > 0 {
		if e, ok := r.cache[host]; ok && time.Now().Before(e.expires) {
			return e.addrs, nil
		}
	}

	addrs, err := r.lookup(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolve %q: %w", host, err)
	}

	if r.ttl > 0 {
		r.cache[host] = entry{addrs: addrs, expires: time.Now().Add(r.ttl)}
	}

	return addrs, nil
}

// First returns the first resolved IP address for host.
func (r *Resolver) First(ctx context.Context, host string) (string, error) {
	addrs, err := r.Resolve(ctx, host)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("resolve %q: no addresses returned", host)
	}
	return addrs[0], nil
}

// Invalidate removes a cached entry for the given host.
func (r *Resolver) Invalidate(host string) {
	delete(r.cache, host)
}

func defaultLookup(ctx context.Context, host string) ([]string, error) {
	return net.DefaultResolver.LookupHost(ctx, host)
}
