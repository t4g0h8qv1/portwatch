// Package portwatch provides the core monitoring loop and supporting
// managers for portwatch.
//
// # Scan Window Manager
//
// ScanWindowManager restricts port scans to a configured daily time window
// on a per-target basis. When no window is registered for a target, scans
// are always permitted.
//
// Example:
//
//	m := portwatch.NewScanWindowManager()
//	_ = m.Set("prod-host", portwatch.WindowConfig{
//		Start: 9 * time.Hour,  // 09:00
//		End:   17 * time.Hour, // 17:00
//	})
//	if m.Allowed("prod-host") {
//		// perform scan
//	}
package portwatch
