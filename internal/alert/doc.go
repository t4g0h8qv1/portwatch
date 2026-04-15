// Package alert implements alerting for portwatch.
//
// It defines the Notifier interface and a built-in StdoutNotifier that
// formats alert messages to any io.Writer. The Evaluate function acts as
// the bridge between a baseline.Diff result and the chosen notifier:
//
//	// After scanning and computing a diff:
//	//   diff, _ := baseline.Diff(saved, current)
//	//   alert.Evaluate(alert.NewStdoutNotifier(), diff)
//
// Alert levels:
//
//	- INFO  : informational, no port changes.
//	- WARN  : previously open ports have disappeared.
//	- ERROR : new, unexpected ports have appeared.
//
// Additional Notifier implementations (e.g. email, webhook, Slack) can be
// added by satisfying the Notifier interface.
package alert
