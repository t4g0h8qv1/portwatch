// Package portwatch provides the core scan-diff-notify loop used by the
// portwatch CLI. It ties together scanning, baseline comparison, and
// notification into a single Run call suitable for both one-shot and
// scheduled execution.
//
// Basic usage:
//
//	err := portwatch.Run(ctx, portwatch.Config{
//		Host:     "192.168.1.1",
//		Ports:    ports,
//		Baseline: "/var/lib/portwatch/baseline.json",
//	}, myNotifier)
package portwatch
