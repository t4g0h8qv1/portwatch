package portsampler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portsampler"
)

func stubScan(ports []int, err error) portsampler.ScanFunc {
	return func(_ context.Context, _ string) ([]int, error) {
		return ports, err
	}
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := portsampler.New("localhost", 0, 0, stubScan(nil, nil))
	if !errors.Is(err, portsampler.ErrInvalidInterval) {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
}

func TestNew_NilScanFunc(t *testing.T) {
	_, err := portsampler.New("localhost", time.Second, 0, nil)
	if !errors.Is(err, portsampler.ErrNilScanFunc) {
		t.Fatalf("expected ErrNilScanFunc, got %v", err)
	}
}

func TestNew_InvalidJitter(t *testing.T) {
	_, err := portsampler.New("localhost", time.Second, -1, stubScan(nil, nil))
	if !errors.Is(err, portsampler.ErrInvalidJitter) {
		t.Fatalf("expected ErrInvalidJitter, got %v", err)
	}
}

func TestRun_ReceivesSample(t *testing.T) {
	ports := []int{80, 443}
	s, err := portsampler.New("localhost", 10*time.Millisecond, 0, stubScan(ports, nil))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	ch := s.Run(ctx)
	got := <-ch
	if got.Err != nil {
		t.Fatalf("unexpected error: %v", got.Err)
	}
	if len(got.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got.Ports))
	}
	if got.Host != "localhost" {
		t.Fatalf("expected host localhost, got %s", got.Host)
	}
}

func TestRun_PropagatesScanError(t *testing.T) {
	scanErr := errors.New("scan failed")
	s, _ := portsampler.New("host", 10*time.Millisecond, 0, stubScan(nil, scanErr))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	got := <-s.Run(ctx)
	if !errors.Is(got.Err, scanErr) {
		t.Fatalf("expected scan error, got %v", got.Err)
	}
}

func TestLast_NilBeforeFirstSample(t *testing.T) {
	s, _ := portsampler.New("host", time.Second, 0, stubScan(nil, nil))
	if s.Last() != nil {
		t.Fatal("expected nil before first sample")
	}
}

func TestLast_UpdatedAfterSample(t *testing.T) {
	s, _ := portsampler.New("host", 10*time.Millisecond, 0, stubScan([]int{22}, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	<-s.Run(ctx)
	if s.Last() == nil {
		t.Fatal("expected non-nil last sample")
	}
}
