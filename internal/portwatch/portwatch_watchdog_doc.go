// Package portwatch provides the core monitoring loop and supporting
// managers for portwatch.
//
// # Watchdog
//
// WatchdogManager tracks whether monitored targets have been heard from
// within a configurable silence window. Each time a scan completes
// successfully for a target, the caller should invoke Ping to reset the
// timer. If a target has not been pinged within MaxSilence, IsExpired
// returns true and the target appears in the Expired list.
//
// This is useful for detecting targets that have silently stopped being
// scanned — for example, due to a misconfiguration or a crashed goroutine —
// without surfacing a scan error.
//
// Usage:
//
//	wdm, err := NewWatchdogManager(DefaultWatchdogConfig())
//	if err != nil { ... }
//	// after each successful scan:
//	_ = wdm.Ping(target)
//	// periodically:
//	for _, t := range wdm.Expired() { ... }
package portwatch
