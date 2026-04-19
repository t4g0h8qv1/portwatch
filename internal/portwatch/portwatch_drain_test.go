package portwatch

import (
	"testing"
	"time"
)

func TestNewDrainManager_InvalidTimeout(t *testing.T) {
	_, err := NewDrainManager(DrainConfig{Timeout: 0})
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestNewDrainManager_Valid(t *testing.T) {
	dm, err := NewDrainManager(DefaultDrainConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dm == nil {
		t.Fatal("expected non-nil DrainManager")
	}
}

func TestAcquireRelease_InFlight(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())

	dm.Acquire("host-a")
	dm.Acquire("host-a")
	dm.Acquire("host-b")

	if got := dm.InFlight(); got != 3 {
		t.Fatalf("expected 3 in-flight, got %d", got)
	}

	dm.Release("host-a")
	if got := dm.InFlight(); got != 2 {
		t.Fatalf("expected 2 in-flight after release, got %d", got)
	}
}

func TestRelease_BelowZeroIsNoop(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())
	dm.Release("ghost") // never acquired
	if dm.InFlight() != 0 {
		t.Fatal("expected 0 in-flight")
	}
}

func TestWait_CompletesWhenEmpty(t *testing.T) {
	dm, _ := NewDrainManager(DefaultDrainConfig())
	if err := dm.Wait(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWait_CompletesAfterRelease(t *testing.T) {
	dm, _ := NewDrainManager(DrainConfig{Timeout: 2 * time.Second, PollInterval: 10 * time.Millisecond})
	dm.Acquire("host-a")

	go func() {
		time.Sleep(50 * time.Millisecond)
		dm.Release("host-a")
	}()

	if err := dm.Wait(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWait_TimesOut(t *testing.T) {
	dm, _ := NewDrainManager(DrainConfig{Timeout: 50 * time.Millisecond, PollInterval: 5 * time.Millisecond})
	dm.Acquire("stuck-host")

	err := dm.Wait()
	if err != errDrainTimeout {
		t.Fatalf("expected errDrainTimeout, got %v", err)
	}
}
