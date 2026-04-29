// Package portwatch provides the core run loop and supporting managers
// for the portwatch monitoring tool.
//
// # DigestManager
//
// DigestManager computes a deterministic SHA-256 digest of the port set
// observed for each target host. On every scan the caller invokes Record;
// the returned changed flag is true whenever the port set differs from the
// previous observation, making it a lightweight change-detection mechanism
// that does not require storing full port lists across scans.
//
// Port order is normalised before hashing, so {80,443} and {443,80} produce
// the same digest.
//
// Usage:
//
//	dm := portwatch.NewDigestManager()
//	digest, changed := dm.Record(target, openPorts, time.Now())
//	if changed {
//	    // trigger alert or downstream processing
//	}
package portwatch
