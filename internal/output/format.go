package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a raw string (e.g. from a CLI flag or config file) into
// a Format constant. It is case-insensitive. An error is returned for unknown
// values so callers can surface helpful diagnostics.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case string(FormatText), "":
		return FormatText, nil
	case string(FormatJSON):
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("output: unknown format %q (want \"text\" or \"json\")", s)
	}
}

// String implements fmt.Stringer.
func (f Format) String() string {
	if f == "" {
		return string(FormatText)
	}
	return string(f)
}
