// Package portwatch provides the core monitoring loop for portwatch.
//
// # Snapshot Store
//
// SnapshotStore keeps the most recent scan result for every monitored
// target. Each call to Record replaces the previous entry for that
// target, so the store always reflects the latest known state.
//
// Typical usage:
//
//	store := portwatch.NewSnapshotStore()
//	store.Record("192.168.1.1", openPorts)
//	if entry, ok := store.Get("192.168.1.1"); ok {
//		fmt.Println(entry.Ports)
//	}
//
// WriteSnapshotTable renders a compact summary of all stored snapshots
// to any io.Writer, suitable for CLI status output.
package portwatch
