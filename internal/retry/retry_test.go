package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after eventual success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		return errTemp
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelledMidRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	cfg := retry.Config{MaxAttempts: 5, InitialDelay: 50 * time.Millisecond, MaxDelay: time.Second, Multiplier: 2}
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := retry.Do(ctx, cfg, func() error {
		calls++
		return errTemp
	})
	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}

func TestDo_InvalidMaxAttempts(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 0}
	err := retry.Do(context.Background(), cfg, func() error { return nil })
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}
