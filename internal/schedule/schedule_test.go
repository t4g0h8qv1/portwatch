package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

func TestRun_ExecutesImmediately(t *testing.T) {
	var count int32
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	job := &schedule.Job{
		Interval: 1 * time.Hour, // long interval — only immediate run expected
		Task: func(_ context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		},
	}

	go job.Run(ctx)
	<-ctx.Done()

	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("expected 1 immediate execution, got %d", got)
	}
}

func TestRun_TicksOnInterval(t *testing.T) {
	var count int32
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	job := &schedule.Job{
		Interval: 100 * time.Millisecond,
		Task: func(_ context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		},
	}

	go job.Run(ctx)
	<-ctx.Done()

	// Expect: immediate + ~3 ticks within 350 ms
	if got := atomic.LoadInt32(&count); got < 3 {
		t.Errorf("expected at least 3 executions, got %d", got)
	}
}

func TestRun_CallsOnError(t *testing.T) {
	var errCount int32
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	job := &schedule.Job{
		Interval: 1 * time.Hour,
		Task: func(_ context.Context) error {
			return errors.New("scan failed")
		},
		OnError: func(_ error) {
			atomic.AddInt32(&errCount, 1)
		},
	}

	go job.Run(ctx)
	<-ctx.Done()

	if got := atomic.LoadInt32(&errCount); got != 1 {
		t.Errorf("expected 1 error callback, got %d", got)
	}
}

func TestParseDuration_Valid(t *testing.T) {
	d, err := schedule.ParseDuration("5m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 5*time.Minute {
		t.Errorf("expected 5m, got %v", d)
	}
}

func TestParseDuration_Invalid(t *testing.T) {
	if _, err := schedule.ParseDuration("notaduration"); err == nil {
		t.Error("expected error for invalid duration")
	}
	if _, err := schedule.ParseDuration("-1m"); err == nil {
		t.Error("expected error for negative duration")
	}
}
