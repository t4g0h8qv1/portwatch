package portwatch_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwatch"
)

func TestDefaultConfig_Defaults(t *testing.T) {
	c := portwatch.DefaultConfig()
	if c.Timeout != 500*time.Millisecond {
		t.Errorf("expected 500ms timeout, got %v", c.Timeout)
	}
	if c.BaselinePath == "" {
		t.Error("expected non-empty BaselinePath")
	}
	if c.AlertOnGone {
		t.Error("expected AlertOnGone to default to false")
	}
}

func TestValidate_MissingTarget(t *testing.T) {
	c := portwatch.DefaultConfig()
	c.Ports = []int{80}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing target")
	}
}

func TestValidate_MissingPorts(t *testing.T) {
	c := portwatch.DefaultConfig()
	c.Target = "localhost"
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing ports")
	}
}

func TestValidate_InvalidTimeout(t *testing.T) {
	c := portwatch.DefaultConfig()
	c.Target = "localhost"
	c.Ports = []int{80}
	c.Timeout = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestValidate_Valid(t *testing.T) {
	c := portwatch.DefaultConfig()
	c.Target = "localhost"
	c.Ports = []int{80, 443}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
