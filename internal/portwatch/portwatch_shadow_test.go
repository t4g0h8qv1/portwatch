package portwatch

import (
	"testing"
	"time"
)

func TestNewShadowTracker_InvalidObs(t *testing.T) {
	_, err := NewShadowTracker(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for minObs=0")
	}
}

func TestNewShadowTracker_InvalidAge(t *testing.T) {
	_, err := NewShadowTracker(2, 0)
	if err == nil {
		t.Fatal("expected error for maxAge=0")
	}
}

func TestNewShadowTracker_Valid(t *testing.T) {
	st, err := NewShadowTracker(2, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestObserve_NotConfirmedBeforeThreshold(t *testing.T) {
	st, _ := NewShadowTracker(3, time.Minute)
	for i := 0; i < 2; i++ {
		if st.Observe("host", 8080) {
			t.Fatalf("should not confirm on observation %d", i+1)
		}
	}
}

func TestObserve_ConfirmedAtThreshold(t *testing.T) {
	st, _ := NewShadowTracker(3, time.Minute)
	for i := 0; i < 2; i++ {
		st.Observe("host", 8080)
	}
	if !st.Observe("host", 8080) {
		t.Fatal("expected confirmation on 3rd observation")
	}
}

func TestObserve_ConfirmedPortRemovedFromPending(t *testing.T) {
	st, _ := NewShadowTracker(2, time.Minute)
	st.Observe("host", 9090)
	st.Observe("host", 9090) // confirms
	pending := st.Pending("host")
	for _, e := range pending {
		if e.Port == 9090 {
			t.Fatal("confirmed port should not appear in pending")
		}
	}
}

func TestPending_ReturnsUnconfirmed(t *testing.T) {
	st, _ := NewShadowTracker(3, time.Minute)
	st.Observe("host", 1234)
	pending := st.Pending("host")
	if len(pending) != 1 || pending[0].Port != 1234 {
		t.Fatalf("expected port 1234 in pending, got %v", pending)
	}
}

func TestPending_PrunesExpired(t *testing.T) {
	now := time.Now()
	st, _ := NewShadowTracker(3, time.Minute)
	st.now = func() time.Time { return now }
	st.Observe("host", 5555)
	// advance time past maxAge
	st.now = func() time.Time { return now.Add(2 * time.Minute) }
	pending := st.Pending("host")
	if len(pending) != 0 {
		t.Fatalf("expected expired entry pruned, got %v", pending)
	}
}

func TestObserve_IndependentTargets(t *testing.T) {
	st, _ := NewShadowTracker(2, time.Minute)
	st.Observe("host-a", 80)
	if st.Observe("host-b", 80) {
		t.Fatal("observations on different targets should be independent")
	}
}
