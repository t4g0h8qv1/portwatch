// Package portwatch provides the top-level orchestration loop for portwatch.
//
// Run performs one complete scan cycle:
//
//  1. Scan the target host for open ports using [scanner.OpenPorts].
//  2. Load (or create) a [baseline.Baseline] from disk.
//  3. Diff the current open ports against the baseline via [alert.Evaluate].
//  4. If changes are detected, invoke the configured [notify.Notifier] and
//     persist an updated baseline.
//  5. Append the scan result to the [history] log.
//
// Callers are responsible for scheduling repeated calls to Run; see the
// [schedule] package for a helper that drives periodic execution.
package portwatch
