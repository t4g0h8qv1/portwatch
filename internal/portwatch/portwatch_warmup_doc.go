// Package portwatch provides the core portwatch engine.
//
// # Warmup
//
// WarmupManager tracks whether a scan target has completed its warmup period.
// During warmup, alerts and policy decisions may be suppressed to avoid noise
// from newly registered hosts whose baseline has not yet stabilised.
//
// A target is considered warm when either:
//   - It has accumulated at least MinScans successful scans, or
//   - MaxWait has elapsed since the first scan was recorded.
//
// Usage:
//
//	wm, err := portwatch.NewWarmupManager(portwatch.DefaultWarmupConfig())
//	if err != nil { ... }
//
//	wm.RecordScan("192.168.1.1")
//	if wm.IsWarm("192.168.1.1") {
//	    // proceed with alerting
//	}
package portwatch
