// Package portwatch provides the top-level orchestration loop for portwatch.
//
// Run performs a single scan-diff-notify cycle:
//
//  1. Scan the target host for open ports.
//  2. Load the persisted baseline (or create one on first run).
//  3. Compute a diff between the baseline and current scan.
//  4. Notify via the configured Notifier if changes are detected.
//  5. Persist the updated baseline for the next cycle.
package portwatch
