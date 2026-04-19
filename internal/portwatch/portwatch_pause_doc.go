// Package portwatch provides the core scan-loop and supporting types for
// portwatch.
//
// # PauseManager
//
// PauseManager allows the scan loop to be temporarily suspended without
// stopping the process. A pause is set for a fixed duration and expires
// automatically; it can also be lifted early via Resume.
//
// Typical usage:
//
//	pm := portwatch.NewPauseManager()
//
//	// Pause for a 30-minute maintenance window.
//	if err := pm.Pause(30*time.Minute, "scheduled maintenance"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Inside the scan loop:
//	if pm.IsPaused() {
//		return // skip this tick
//	}
package portwatch
