// Package output provides types and helpers for rendering port-scan results
// to an io.Writer in either human-readable text or machine-readable JSON format.
//
// # Basic usage
//
//	w := output.NewWriter(os.Stdout, output.FormatText)
//	err := w.Write(output.Result{
//		Host:      "192.168.1.1",
//		OpenPorts: []int{22, 80},
//		NewPorts:  []int{80},
//		ScannedAt: time.Now(),
//	})
//
// Supported formats are FormatText (default, tabwriter-aligned) and FormatJSON
// (indented JSON suitable for piping to jq or log aggregators).
package output
