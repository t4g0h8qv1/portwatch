// Package portdiff computes and formats the difference between two sets of
// open ports, typically from consecutive scans of the same host.
//
// Usage:
//
//	result := portdiff.Compute(previousPorts, currentPorts)
//	if result.HasChanges() {
//		portdiff.Format(os.Stdout, result)
//	}
package portdiff
