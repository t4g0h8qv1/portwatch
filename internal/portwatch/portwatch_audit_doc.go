// Package portwatch provides the core monitoring loop and supporting
// infrastructure for portwatch.
//
// # Audit Log
//
// AuditLog maintains a bounded, in-memory trail of scan events for all
// monitored targets. Each AuditEvent captures the target host, the ports
// observed, any ports that opened or closed since the previous scan, the
// wall-clock timestamp of the scan, and any error that occurred.
//
// Usage:
//
//	log, err := portwatch.NewAuditLog(200)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Record(portwatch.AuditEvent{
//		Target: "192.168.1.1",
//		Ports:  []int{22, 80},
//		Opened: []int{80},
//	})
//	events := log.All()
//	portwatch.WriteAuditTable(os.Stdout, events)
package portwatch
