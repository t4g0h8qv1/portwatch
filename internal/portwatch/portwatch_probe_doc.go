// Package portwatch provides the core scan-loop, alerting, and lifecycle
// management for portwatch.
//
// # Probe Manager
//
// ProbeManager tracks per-target liveness by counting consecutive scan
// failures. When a target exceeds MaxConsecutiveFailures it is marked dead
// and skipped by the scan loop until it recovers.
//
// A dead target recovers after RecoveryThreshold consecutive successful scans,
// preventing flapping from causing spurious alerts.
//
// Usage:
//
//	cfg := portwatch.DefaultProbeConfig()
//	pm, err := portwatch.NewProbeManager(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// after each scan:
//	if scanErr != nil {
//	    pm.RecordFailure(target)
//	} else {
//	    pm.RecordSuccess(target)
//	}
//	if pm.IsDead(target) {
//	    // skip expensive scan or alert operator
//	}
package portwatch
