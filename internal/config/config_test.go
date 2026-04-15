package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	p := writeTemp(t, `
target: localhost
ports: "1-1024"
baseline: baseline.json
alert:
  stdout: true
scan:
  timeout_ms: 300
  concurrency: 50
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Target != "localhost" {
		t.Errorf("target: got %q, want %q", cfg.Target, "localhost")
	}
	if cfg.Scan.TimeoutMs != 300 {
		t.Errorf("timeout_ms: got %d, want 300", cfg.Scan.TimeoutMs)
	}
	if cfg.Scan.Concurrency != 50 {
		t.Errorf("concurrency: got %d, want 50", cfg.Scan.Concurrency)
	}
}

func TestLoad_Defaults(t *testing.T) {
	p := writeTemp(t, `target: "192.168.1.1"
ports: "80,443"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scan.TimeoutMs != 500 {
		t.Errorf("default timeout_ms: got %d, want 500", cfg.Scan.TimeoutMs)
	}
	if cfg.Scan.Concurrency != 100 {
		t.Errorf("default concurrency: got %d, want 100", cfg.Scan.Concurrency)
	}
}

func TestLoad_MissingTarget(t *testing.T) {
	p := writeTemp(t, `ports: "1-100"
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing target, got nil")
	}
}

func TestLoad_MissingPorts(t *testing.T) {
	p := writeTemp(t, `target: localhost
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing ports, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTemp(t, `target: [unclosed`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
