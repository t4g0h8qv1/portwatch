// Package portwatch provides the ScanConcurrencyManager, which limits the
// number of port scans that may execute simultaneously across all targets.
//
// # Overview
//
// When portwatch monitors many hosts the scheduler may attempt to trigger
// multiple scans at the same time. ScanConcurrencyManager uses a semaphore
// to cap the number of in-flight scans, preventing resource exhaustion on
// the host running portwatch.
//
// # Usage
//
//	cfg := portwatch.DefaultConcurrencyConfig()
//	cfg.MaxConcurrent = 6
//	mgr, err := portwatch.NewScanConcurrencyManager(cfg)
//	if err != nil { ... }
//
//	if err := mgr.Acquire(target); err != nil {
//		// handle ErrConcurrencyTimeout or ErrTargetAlreadyScanning
//	}
//	defer mgr.Release(target)
//	// perform scan …
package portwatch
