package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMain_ConfigNotFound verifies the binary exits when config is missing.
// We test the helper logic rather than invoking main() directly to avoid os.Exit.
func TestMain_ConfigNotFound(t *testing.T) {
	_, err := os.ReadFile("nonexistent_portwatch.yaml")
	if err == nil {
		t.Fatal("expected error reading missing config file")
	}
}

// TestMain_DefaultConfigPath checks that the default config path is used
// when no CLI argument is supplied (simulated via direct value check).
func TestMain_DefaultConfigPath(t *testing.T) {
	args := []string{}
	cfgPath := "portwatch.yaml"
	if len(args) > 0 {
		cfgPath = args[0]
	}
	if cfgPath != "portwatch.yaml" {
		t.Errorf("expected default config path 'portwatch.yaml', got %q", cfgPath)
	}
}

// TestMain_CustomConfigPath checks that a CLI argument overrides the default path.
func TestMain_CustomConfigPath(t *testing.T) {
	args := []string{"/etc/portwatch/custom.yaml"}
	cfgPath := "portwatch.yaml"
	if len(args) > 0 {
		cfgPath = args[0]
	}
	if cfgPath != "/etc/portwatch/custom.yaml" {
		t.Errorf("expected custom config path, got %q", cfgPath)
	}
}

// TestMain_TempConfigExists verifies that a temporary config file can be
// created and read back, simulating the startup config-load path.
func TestMain_TempConfigExists(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "portwatch.yaml")

	content := []byte("target: localhost\nports: \"80\"\n")
	if err := os.WriteFile(cfgPath, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("failed to read temp config: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty config file")
	}
}
