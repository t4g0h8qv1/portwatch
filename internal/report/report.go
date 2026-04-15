// Package report provides formatted output for port scan results.
package report

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the data for a single scan report.
type Report struct {
	Target    string
	ScannedAt time.Time
	NewPorts  []int
	GonePorts []int
	OpenPorts []int
}

// Writer writes reports to an output destination.
type Writer struct {
	out    io.Writer
	format Format
}

// NewWriter creates a new Writer with the given format.
// If out is nil, os.Stdout is used.
func NewWriter(out io.Writer, format Format) *Writer {
	if out == nil {
		out = os.Stdout
	}
	return &Writer{out: out, format: format}
}

// Write renders the report to the writer's output.
func (w *Writer) Write(r *Report) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(r)
	default:
		return w.writeText(r)
	}
}

func (w *Writer) writeText(r *Report) error {
	fmt.Fprintf(w.out, "=== Port Watch Report ===\n")
	fmt.Fprintf(w.out, "Target:     %s\n", r.Target)
	fmt.Fprintf(w.out, "Scanned At: %s\n", r.ScannedAt.Format(time.RFC3339))
	fmt.Fprintf(w.out, "Open Ports: %s\n", joinInts(r.OpenPorts))
	if len(r.NewPorts) > 0 {
		fmt.Fprintf(w.out, "NEW Ports:  %s\n", joinInts(r.NewPorts))
	}
	if len(r.GonePorts) > 0 {
		fmt.Fprintf(w.out, "Gone Ports: %s\n", joinInts(r.GonePorts))
	}
	return nil
}

func (w *Writer) writeJSON(r *Report) error {
	fmt.Fprintf(w.out, `{"target":%q,"scanned_at":%q,"open_ports":%s,"new_ports":%s,"gone_ports":%s}\n`,
		r.Target,
		r.ScannedAt.Format(time.RFC3339),
		jsonInts(r.OpenPorts),
		jsonInts(r.NewPorts),
		jsonInts(r.GonePorts),
	)
	return nil
}

// FromAlert builds a Report from an alert.Result and metadata.
func FromAlert(target string, result alert.Result, open []int) *Report {
	return &Report{
		Target:    target,
		ScannedAt: time.Now(),
		NewPorts:  result.NewPorts,
		GonePorts: result.GonePorts,
		OpenPorts: open,
	}
}

func joinInts(vals []int) string {
	if len(vals) == 0 {
		return "(none)"
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ", ")
}

func jsonInts(vals []int) string {
	if len(vals) == 0 {
		return "[]"
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
