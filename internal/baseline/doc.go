// Package baseline provides functionality for saving, loading, and comparing
// port scan baselines for portwatch.
//
// A baseline captures the set of expected open ports for a given host at a
// point in time. It can be persisted to disk as a JSON file and later loaded
// to compare against a fresh scan, enabling detection of newly opened or
// unexpectedly closed ports.
//
// Typical usage:
//
//	// Create and save a baseline from an initial scan
//	b := baseline.New("192.168.1.10", openPorts)
//	if err := b.Save(".portwatch/baseline.json"); err != nil {
//		log.Fatal(err)
//	}
//
//	// Later, load and compare against current scan results
//	b, err := baseline.Load(".portwatch/baseline.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	newPorts, missingPorts := b.Diff(currentPorts)
package baseline
