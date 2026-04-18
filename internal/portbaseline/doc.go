// Package portbaseline provides a versioned, per-host port baseline store.
//
// Baselines record the set of expected open ports for each monitored host.
// Each update increments a version counter so callers can detect staleness.
// The store is persisted as JSON and can be reloaded across process restarts.
//
// Usage:
//
//	s, err := portbaseline.Load("/var/lib/portwatch/baselines.json")
//	s.Set("192.168.1.1", []int{22, 80, 443})
//	_ = s.Save()
package portbaseline
