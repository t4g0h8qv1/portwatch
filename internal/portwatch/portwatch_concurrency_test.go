package portwatch

import (
	"sort"
	"sync"
	"testing"
	"time"
)

func TestNewScanConcurrencyManager_InvalidMax(t *testing.T) {
	_, err := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: 0, AcquireTimeout: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxConcurrent=0")
	}
}

func TestNewScanConcurrencyManager_InvalidTimeout(t *testing.T) {
	_, err := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: 2, AcquireTimeout: 0})
	if err == nil {
		t.Fatal("expected error for AcquireTimeout=0")
	}
}

func TestNewScanConcurrencyManager_Valid(t *testing.T) {
	m, err := NewScanConcurrencyManager(DefaultConcurrencyConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Slots(); got != 4 {
		t.Fatalf("expected 4 free slots, got %d", got)
	}
}

func TestAcquire_ConsumesSlot(t *testing.T) {
	m, _ := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: 2, AcquireTimeout: time.Second})
	if err := m.Acquire("host-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Slots(); got != 1 {
		t.Fatalf("expected 1 free slot after acquire, got %d", got)
	}
	m.Release("host-a")
	if got := m.Slots(); got != 2 {
		t.Fatalf("expected 2 free slots after release, got %d", got)
	}
}

func TestAcquire_TargetAlreadyScanning(t *testing.T) {
	m, _ := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: 4, AcquireTimeout: time.Second})
	if err := m.Acquire("host-a"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if err := m.Acquire("host-a"); err != ErrTargetAlreadyScanning {
		t.Fatalf("expected ErrTargetAlreadyScanning, got %v", err)
	}
}

func TestAcquire_Timeout(t *testing.T) {
	m, _ := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: 1, AcquireTimeout: 20 * time.Millisecond})
	if err := m.Acquire("host-a"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if err := m.Acquire("host-b"); err != ErrConcurrencyTimeout {
		t.Fatalf("expected ErrConcurrencyTimeout, got %v", err)
	}
}

func TestRelease_Noop(t *testing.T) {
	m, _ := NewScanConcurrencyManager(DefaultConcurrencyConfig())
	// releasing a target that was never acquired should not panic or block
	m.Release("ghost")
	if got := m.Slots(); got != 4 {
		t.Fatalf("expected 4 free slots, got %d", got)
	}
}

func TestActive_ReflectsHeldTargets(t *testing.T) {
	m, _ := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: 4, AcquireTimeout: time.Second})
	_ = m.Acquire("alpha")
	_ = m.Acquire("beta")

	active := m.Active()
	sort.Strings(active)
	if len(active) != 2 || active[0] != "alpha" || active[1] != "beta" {
		t.Fatalf("unexpected active set: %v", active)
	}

	m.Release("alpha")
	active = m.Active()
	if len(active) != 1 || active[0] != "beta" {
		t.Fatalf("expected only beta active, got %v", active)
	}
}

func TestAcquire_ConcurrentSafe(t *testing.T) {
	max := 8
	m, _ := NewScanConcurrencyManager(ConcurrencyConfig{MaxConcurrent: max, AcquireTimeout: time.Second})
	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		target := string(rune('a' + i))
		wg.Add(1)
		go func(tgt string) {
			defer wg.Done()
			if err := m.Acquire(tgt); err != nil {
				t.Errorf("acquire %s: %v", tgt, err)
				return
			}
			time.Sleep(5 * time.Millisecond)
			m.Release(tgt)
		}(target)
	}
	wg.Wait()
	if got := m.Slots(); got != max {
		t.Fatalf("expected all %d slots free after completion, got %d", max, got)
	}
}
