// Package resolve provides hostname-to-IP resolution for portwatch targets.
//
// It wraps the standard net.DefaultResolver and adds an optional in-memory
// cache to avoid redundant DNS lookups during repeated scan cycles.
//
// # Cache behaviour
//
// Resolved addresses are stored with a configurable TTL. Once the TTL expires
// the next call triggers a fresh DNS lookup and refreshes the cache entry.
// A TTL of zero disables caching entirely.
//
// # Usage
//
//	r := resolve.New(5 * time.Minute)
//	ip, err := r.First(ctx, "example.com")
//
// If the host argument is already a valid IP address it is returned
// unchanged without performing a DNS lookup.
//
// # All addresses
//
// To retrieve every address associated with a hostname use All instead of
// First:
//
//	ips, err := r.All(ctx, "example.com")
//	for _, ip := range ips {
//		fmt.Println(ip)
//	}
package resolve
