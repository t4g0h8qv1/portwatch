// Package output provides formatters for writing scan results to various
// destinations such as stdout, files, or structured streams.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Format represents the output format for scan results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Result holds a single scan result for output purposes.
type Result struct {
	Host      string    `json:"host"`
	OpenPorts []int     `json:"open_ports"`
	NewPorts  []int     `json:"new_ports,omitempty"`
	GonePorts []int     `json:"gone_ports,omitempty"`
	ScannedAt time.Time `json:"scanned_at"`
}

// Writer writes Results to an io.Writer in a given Format.
type Writer struct {
	w      io.Writer
	format Format
}

// NewWriter returns a Writer that writes to w using the given format.
func NewWriter(w io.Writer, format Format) *Writer {
	return &Writer{w: w, format: format}
}

// Write renders r to the underlying writer.
func (wr *Writer) Write(r Result) error {
	switch wr.format {
	case FormatJSON:
		return wr.writeJSON(r)
	default:
		return wr.writeText(r)
	}
}

func (wr *Writer) writeJSON(r Result) error {
	enc := json.NewEncoder(wr.w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func (wr *Writer) writeText(r Result) error {
	tw := tabwriter.NewWriter(wr.w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Host:\t%s\n", r.Host)
	fmt.Fprintf(tw, "Scanned at:\t%s\n", r.ScannedAt.Format(time.RFC3339))
	fmt.Fprintf(tw, "Open ports:\t%v\n", r.OpenPorts)
	if len(r.NewPorts) > 0 {
		fmt.Fprintf(tw, "New ports:\t%v\n", r.NewPorts)
	}
	if len(r.GonePorts) > 0 {
		fmt.Fprintf(tw, "Gone ports:\t%v\n", r.GonePorts)
	}
	return tw.Flush()
}
