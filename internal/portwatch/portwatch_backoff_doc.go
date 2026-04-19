// Package portwatch provides the core run loop and supporting types for
// portwatch.
//
// # Backoff
//
// BackoffManager implements per-target exponential backoff for scan errors.
// When a scan for a target fails repeatedly, the caller can use
// RecordFailure to obtain the recommended wait duration before the next
// attempt. Once the target recovers, RecordSuccess resets the counter so
// subsequent failures start from the initial interval again.
//
// Example:
//
//	bm := portwatch.NewBackoffManager(portwatch.DefaultBackoffConfig())
//
//	if err := scan(target); err != nil {
//		wait := bm.RecordFailure(target)
//		time.Sleep(wait)
//	} else {
//		bm.RecordSuccess(target)
//	}
package portwatch
