// Package cache implements a lightweight TTL-based in-memory cache for
// storing port-scan results per host.
//
// Typical usage:
//
//	c, err := cache.New(30 * time.Second)
//	if err != nil { ... }
//
//	c.Set("192.168.1.1", []int{22, 80, 443})
//
//	if entry, ok := c.Get("192.168.1.1"); ok {
//		fmt.Println(entry.Ports)
//	}
//
// Entries that have exceeded their TTL are considered expired and will not be
// returned by Get. Call Prune periodically to reclaim memory.
package cache
