// Package portwatch provides the core monitoring loop and supporting
// subsystems for portwatch.
//
// # Scan Quota
//
// ScanQuotaManager enforces a maximum number of scans per target within a
// configurable rolling time window. This prevents runaway scan loops from
// overwhelming a target host or consuming excessive local resources.
//
// Usage:
//
//	q, err := portwatch.NewScanQuotaManager(portwatch.ScanQuotaConfig{
//	    MaxScansPerHour: 30,
//	    Window:          time.Hour,
//	})
//	if err != nil { ... }
//
//	if q.Allow(target) {
//	    // perform scan
//	}
//
// WriteQuotaTable and QuotaSummary provide formatted output suitable for
// CLI display or log output.
package portwatch
