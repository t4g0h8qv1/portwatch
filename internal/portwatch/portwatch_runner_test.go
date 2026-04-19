package portwatch_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwatch"
)

func TestNewRunner_InvalidInterval(t *testing.T) {
	_, err := portwatch.NewRunner(portwatch.RunnerConfig{
		Interval: 0,
		Out:      &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewRunner_NilWriter(t *testing.T) {
	_, err := portwatch.NewRunner(portwatch.RunnerConfig{
		Interval: time.Second,
		Out:      nil,
	})
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestNewRunner_Valid(t *testing.T) {
	r, err := portwatch.NewRunner(portwatch.RunnerConfig{
		Interval: time.Second,
		Out:      &bytes.Buffer{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestRunner_MaxScans(t *testing.T) {
	buf := &bytes.Buffer{}
	port := freePort(t)
	cfg := portwatch.DefaultConfig()
	cfg.Target = "127.0.0.1"
	cfg.Ports = []int{port}

	r, err := portwatch.NewRunner(portwatch.RunnerConfig{
		ScanConfig: cfg,
		Interval:   10 * time.Millisecond,
		MaxScans:   2,
		Out:        buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	result := r.Start(ctx)
	if result.ScansCompleted != 2 {
		t.Errorf("expected 2 scans, got %d", result.ScansCompleted)
	}
}

func TestRunner_ContextCancel(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := portwatch.DefaultConfig()
	cfg.Target = "127.0.0.1"
	cfg.Ports = []int{freePort(t)}

	r, err := portwatch.NewRunner(portwatch.RunnerConfig{
		ScanConfig: cfg,
		Interval:   50 * time.Millisecond,
		Out:        buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	result := r.Start(ctx)
	if result.ScansCompleted < 1 {
		t.Error("expected at least one scan before context cancelled")
	}
}
