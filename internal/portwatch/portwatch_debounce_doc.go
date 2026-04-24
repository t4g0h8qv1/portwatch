// Package portwatch provides the core monitoring loop and supporting
// managers for portwatch.
//
// # Debounce Manager
//
// DebounceManager prevents alert storms by suppressing repeated alerts
// for the same target within a configurable time window.
//
// Usage:
//
//	dm, err := portwatch.NewDebounceManager(portwatch.DebounceConfig{
//		Window: 30 * time.Second,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if dm.Ready(target) {
//		// fire alert
//		dm.Observe(target)
//	}
package portwatch
