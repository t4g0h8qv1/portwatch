// Package report provides formatted report generation for portwatch scan results.
//
// It supports multiple output formats (text and JSON) and can be used to
// present scan results to the user in a human-readable or machine-parseable form.
//
// Example usage:
//
//	w := report.NewWriter(os.Stdout, report.FormatText)
//	r := report.FromAlert(target, alertResult, openPorts)
//	w.Write(r)
package report
