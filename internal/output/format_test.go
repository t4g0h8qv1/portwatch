package output_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/output"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  output.Format
	}{
		{"text", output.FormatText},
		{"TEXT", output.FormatText},
		{"json", output.FormatJSON},
		{"JSON", output.FormatJSON},
		{"", output.FormatText},
		{"  text  ", output.FormatText},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := output.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestFormat_String(t *testing.T) {
	if output.FormatJSON.String() != "json" {
		t.Errorf("expected \"json\", got %q", output.FormatJSON.String())
	}
	var f output.Format
	if f.String() != "text" {
		t.Errorf("expected \"text\" for zero value, got %q", f.String())
	}
}
