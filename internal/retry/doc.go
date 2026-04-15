// Package retry provides a context-aware exponential-backoff retry helper.
//
// Usage:
//
//	cfg := retry.DefaultConfig()
//	err := retry.Do(ctx, cfg, func() error {
//	    return doSomethingFallible()
//	})
//
// The helper will call the supplied function up to MaxAttempts times,
// sleeping an exponentially growing delay (capped at MaxDelay) between
// consecutive failures. If the context is cancelled the loop exits
// immediately and wraps the context error.
package retry
