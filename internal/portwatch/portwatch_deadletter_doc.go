// Package portwatch provides the core scanning and monitoring loop for
// portwatch.
//
// # Dead-Letter Queue
//
// DeadLetterQueue captures scan events that could not be processed after all
// retry attempts have been exhausted. Entries are retained in insertion order
// up to a configurable maximum; when the queue is full the oldest entry is
// evicted automatically.
//
// Typical usage:
//
//	q, err := portwatch.NewDeadLetterQueue(100)
//	if err != nil { ... }
//
//	// On unrecoverable scan failure:
//	q.Push(target, scanErr, attempts)
//
//	// Inspect at any time:
//	portwatch.WriteDeadLetterTable(os.Stdout, q.All())
package portwatch
