// Package portwatch provides the core monitoring loop for portwatch.
//
// # Scan Retry Manager
//
// ScanRetryManager tracks per-target retry state for failed port scans.
// It implements truncated exponential backoff: each consecutive failure
// for a given target doubles the wait delay up to a configurable maximum.
//
// Usage:
//
//	cfg := portwatch.DefaultRetryConfig()
//	m, err := portwatch.NewScanRetryManager(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if m.ShouldRetry(target) {
//		delay := m.NextDelay(target)
//		time.Sleep(delay)
//		// ... retry scan ...
//	}
//
//	// On success:
//	m.Reset(target)
//
// ScanRetryManager is safe for concurrent use.
package portwatch
