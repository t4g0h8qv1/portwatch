// Package portwatch provides the core monitoring loop and supporting
// managers for portwatch.
//
// # Roster
//
// RosterManager maintains the canonical set of scan targets known to
// the portwatch runtime. Each target is tracked with its registration
// time, the last time it was successfully observed, and an active flag
// that allows targets to be soft-deleted without losing history.
//
// Typical usage:
//
//	rm := portwatch.NewRosterManager()
//	_ = rm.Register("192.168.1.1")
//	_ = rm.Touch("192.168.1.1", time.Now())
//
//	// Later, retire a target without discarding its record:
//	_ = rm.Deactivate("192.168.1.1")
//
// WriteRosterTable and RosterSummary provide human-readable output
// suitable for CLI display.
package portwatch
