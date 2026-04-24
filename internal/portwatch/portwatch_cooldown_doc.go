// Package portwatch provides the core watch loop and supporting managers
// for portwatch.
//
// # ScanCooldownManager
//
// ScanCooldownManager enforces a minimum gap between successive scans of the
// same target host. This prevents rapid re-scanning when errors occur or when
// a caller drives the scan loop at a high frequency.
//
// Usage:
//
//	cfg := portwatch.DefaultScanCooldownConfig()
//	cfg.MinGap = 10 * time.Second
//
//	m, err := portwatch.NewScanCooldownManager(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if m.Ready(target) {
//		// perform scan
//		m.Observe(target)
//	}
package portwatch
