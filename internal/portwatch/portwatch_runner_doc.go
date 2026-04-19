// Package portwatch provides the core scan-and-alert loop for portwatch.
//
// The Runner type adds a continuous, interval-based execution layer on top of
// the single-shot Run function.  Create a Runner with NewRunner, then call
// Start with a context to begin periodic scanning.  The loop fires immediately
// on the first call and then waits for each tick of the configured interval.
//
// MaxScans can be set to a positive integer to cap the total number of scans
// before Start returns; leaving it at zero means the runner continues until
// the context is cancelled.
//
// RunnerResult summarises the session once Start returns.  Use WriteRunnerResult
// or RunnerSummary to render the result for human consumption.
package portwatch
