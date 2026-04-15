// Package notify provides pluggable notification backends for portwatch alerts.
package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Notifier sends a notification message.
type Notifier interface {
	Notify(subject, body string) error
}

// LogNotifier writes notifications to a writer (default: os.Stderr).
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier writing to stderr.
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{Out: os.Stderr}
}

// Notify writes a timestamped notification to the configured writer.
func (n *LogNotifier) Notify(subject, body string) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(n.Out, "[%s] ALERT: %s\n%s\n", timestamp, subject, body)
	return err
}

// WebhookConfig holds configuration for a webhook notifier.
type WebhookConfig struct {
	URL     string
	Method  string // default: POST
	Headers map[string]string
}

// Multi dispatches a notification to multiple notifiers.
type Multi struct {
	notifiers []Notifier
}

// NewMulti creates a Multi notifier from the provided notifiers.
func NewMulti(notifiers ...Notifier) *Multi {
	return &Multi{notifiers: notifiers}
}

// Notify calls every registered notifier and collects errors.
func (m *Multi) Notify(subject, body string) error {
	var errs []string
	for _, n := range m.notifiers {
		if err := n.Notify(subject, body); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
