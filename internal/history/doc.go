// Package history tracks and persists timestamped port scan results for
// portwatch. Each time a scan completes, a new Entry is appended to a
// JSON file on disk. This allows operators to review historical open-port
// data, spot trends, and correlate alerts with past scan snapshots.
//
// Typical usage:
//
//	err := history.Record("/var/lib/portwatch/history.json", "192.168.1.1", openPorts)
//
//	h, err := history.Load("/var/lib/portwatch/history.json")
//	if last, ok := h.Last(); ok {
//		fmt.Println("last scan at", last.Timestamp)
//	}
package history
