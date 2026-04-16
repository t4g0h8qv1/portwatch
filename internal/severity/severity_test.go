package severity_test

import (
	"testing"

	"github.com/example/portwatch/internal/severity"
)

func TestParseLevel_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  severity.Level
	}{
		{"info", severity.Info},
		{"warning", severity.Warning},
		{"critical", severity.Critical},
	}
	for _, tc := range cases {
		got, err := severity.ParseLevel(tc.input)
		if err != nil {
			t.Fatalf("ParseLevel(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseLevel_Invalid(t *testing.T) {
	_, err := severity.ParseLevel("extreme")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestLevel_String(t *testing.T) {
	if severity.Critical.String() != "critical" {
		t.Errorf("unexpected string for Critical")
	}
	if severity.Info.String() != "info" {
		t.Errorf("unexpected string for Info")
	}
}

func TestClassify_Critical(t *testing.T) {
	c := severity.New([]int{22, 3306}, []int{8080})
	if got := c.Classify(22); got != severity.Critical {
		t.Errorf("expected Critical for port 22, got %v", got)
	}
}

func TestClassify_Warning(t *testing.T) {
	c := severity.New([]int{22}, []int{8080})
	if got := c.Classify(8080); got != severity.Warning {
		t.Errorf("expected Warning for port 8080, got %v", got)
	}
}

func TestClassify_Info(t *testing.T) {
	c := severity.New([]int{22}, []int{8080})
	if got := c.Classify(9999); got != severity.Info {
		t.Errorf("expected Info for unknown port, got %v", got)
	}
}

func TestClassifyAll(t *testing.T) {
	c := severity.New([]int{443}, []int{80})
	result := c.ClassifyAll([]int{443, 80, 9000})
	if result[443] != severity.Critical {
		t.Errorf("port 443 should be Critical")
	}
	if result[80] != severity.Warning {
		t.Errorf("port 80 should be Warning")
	}
	if result[9000] != severity.Info {
		t.Errorf("port 9000 should be Info")
	}
}
