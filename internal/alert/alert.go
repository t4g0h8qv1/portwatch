// Package alert provides alerting mechanisms for unexpected open ports.
package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/baseline"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert represents a single alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	NewPorts  []int
	GonePorts []int
}

// Notifier sends alerts to a destination.
type Notifier interface {
	Notify(a Alert) error
}

// StdoutNotifier writes alerts to an io.Writer (default: os.Stdout).
type StdoutNotifier struct {
	Writer io.Writer
}

// NewStdoutNotifier creates a StdoutNotifier writing to stdout.
func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{Writer: os.Stdout}
}

// Notify formats and writes the alert to the configured writer.
func (s *StdoutNotifier) Notify(a Alert) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s %s\n", a.Level, a.Timestamp.Format(time.RFC3339), a.Message))
	if len(a.NewPorts) > 0 {
		sb.WriteString(fmt.Sprintf("  NEW ports : %v\n", a.NewPorts))
	}
	if len(a.GonePorts) > 0 {
		sb.WriteString(fmt.Sprintf("  GONE ports: %v\n", a.GonePorts))
	}
	_, err := fmt.Fprint(s.Writer, sb.String())
	return err
}

// Evaluate inspects a baseline.Diff and emits an alert via the notifier
// when unexpected changes are detected. Returns true if an alert was sent.
func Evaluate(n Notifier, diff baseline.Diff) (bool, error) {
	if len(diff.New) == 0 && len(diff.Gone) == 0 {
		return false, nil
	}

	lvl := LevelWarn
	if len(diff.New) > 0 {
		lvl = LevelError
	}

	a := Alert{
		Timestamp: time.Now(),
		Level:     lvl,
		Message:   "Port state change detected",
		NewPorts:  diff.New,
		GonePorts: diff.Gone,
	}

	if err := n.Notify(a); err != nil {
		return false, fmt.Errorf("alert: notify: %w", err)
	}
	return true, nil
}
