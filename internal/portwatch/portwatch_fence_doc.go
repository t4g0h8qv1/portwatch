// Package portwatch provides the core run loop and supporting managers
// for the portwatch port-monitoring tool.
//
// # Scan Fence Manager
//
// ScanFenceManager allows operators to temporarily block scans against
// specific targets — for example, during scheduled maintenance windows
// or planned deployments — without stopping the entire portwatch process.
//
// Fences are time-bounded: each fence expires after the configured MaxAge.
// Expired fences are ignored by IsFenced and can be cleaned up explicitly
// with Prune.
//
// Basic usage:
//
//	cfg := portwatch.DefaultFenceConfig()
//	m, err := portwatch.NewScanFenceManager(cfg)
//	if err != nil { ... }
//
//	m.Fence("192.168.1.10", "planned maintenance")
//	if m.IsFenced(target) {
//	    // skip scan
//	}
//	m.Unfence("192.168.1.10")
package portwatch
