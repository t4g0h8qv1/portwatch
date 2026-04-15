// Package suppress manages suppression rules for known or expected open ports.
//
// Operators can add time-bounded suppression entries for ports that are
// intentionally open, preventing portwatch from raising repeated alerts.
// Suppression entries are persisted to disk as JSON and automatically
// expire after a configurable TTL.
//
// Example usage:
//
//	list, err := suppress.Load("/var/lib/portwatch/suppress.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// silence port 8080 for 24 hours
//	list.Add(8080, "temporary dev server", 24*time.Hour)
//	// filter new ports before alerting
//	filtered := list.Filter(newPorts)
package suppress
