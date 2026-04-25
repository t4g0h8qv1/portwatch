// Package portwatch provides the core monitoring loop and supporting
// managers for portwatch.
//
// # Trend Tracking
//
// ScanTrendManager records per-port open/closed observations across
// successive scans and derives a trend direction (rising, falling, or
// stable) based on a configurable sliding window.
//
// Usage:
//
//	cfg := portwatch.DefaultTrendConfig()
//	m, err := portwatch.NewScanTrendManager(cfg)
//	if err != nil { ... }
//
//	// after each scan:
//	for _, port := range openPorts {
//	    m.Observe(target, port, true)
//	}
//
//	// query:
//	trends := m.Trends(target)
//	portwatch.WriteTrendTable(os.Stdout, trends)
//	fmt.Println(portwatch.TrendSummary(trends))
package portwatch
