// Package schedule provides a simple periodic job runner for portwatch.
//
// A Job wraps a scan task and executes it immediately on start, then
// repeatedly at the configured Interval until the context is cancelled.
//
// Example usage:
//
//	job := &schedule.Job{
//		Interval: 5 * time.Minute,
//		Task: func(ctx context.Context) error {
//			// perform scan, compare baseline, send alerts
//			return nil
//		},
//		OnError: func(err error) {
//			log.Printf("scan error: %v", err)
//		},
//	}
//	job.Run(ctx)
package schedule
