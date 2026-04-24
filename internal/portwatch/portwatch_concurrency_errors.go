package portwatch

import "errors"

// ErrConcurrencyTimeout is returned when no scan slot becomes available
// within the configured AcquireTimeout.
var ErrConcurrencyTimeout = errors.New("portwatch: timed out waiting for concurrency slot")

// ErrTargetAlreadyScanning is returned when an Acquire is attempted for a
// target that already holds an active scan slot.
var ErrTargetAlreadyScanning = errors.New("portwatch: target is already being scanned")
