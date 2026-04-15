// Package resolve provides hostname-to-IP resolution for portwatch targets.
//
// It wraps the standard net.DefaultResolver and adds an optional in-memory
// cache to avoid redundant DNS lookups during repeated scan cycles.
//
// Usage:
//
//	r := resolve.New(5 * time.Minute)
//	ip, err := r.First(ctx, "example.com")
//
// If the host argument is already a valid IP address it is returned
// unchanged without performing a DNS lookup.
package resolve
