package suppress_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/suppress"
)

// TestPrune_ExpiredEntriesDroppedOnAdd verifies that expired entries are
// removed from the list when a new entry is added.
func TestPrune_ExpiredEntriesDroppedOnAdd(t *testing.T) {
	path := tempPath(t)
	l, _ := suppress.Load(path)

	// add two entries: one expired, one valid
	_ = l.Add(22, "expired", -time.Minute)
	_ = l.Add(80, "active", time.Hour)

	// reload and check that only the active entry remains on disk
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var raw struct {
		Entries []suppress.Entry `json:"entries"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, e := range raw.Entries {
		if e.Port == 22 {
			t.Error("expired entry for port 22 should have been pruned")
		}
	}
	if len(raw.Entries) != 1 {
		t.Errorf("expected 1 persisted entry, got %d", len(raw.Entries))
	}
}

// TestFilter_EmptyInput returns nil/empty without panicking.
func TestFilter_EmptyInput(t *testing.T) {
	l, _ := suppress.Load(tempPath(t))
	result := l.Filter(nil)
	if result != nil && len(result) != 0 {
		t.Errorf("expected nil or empty slice, got %v", result)
	}
}

// TestMultipleSuppressions checks that several ports can be suppressed at once.
func TestMultipleSuppressions(t *testing.T) {
	path := tempPath(t)
	l, _ := suppress.Load(path)

	ports := []int{21, 23, 25}
	for _, p := range ports {
		if err := l.Add(p, "legacy", time.Hour); err != nil {
			t.Fatalf("Add(%d): %v", p, err)
		}
	}

	input := []int{21, 22, 23, 24, 25}
	result := l.Filter(input)
	if len(result) != 2 {
		t.Errorf("expected 2 unsuppressed ports, got %d: %v", len(result), result)
	}
	for _, p := range result {
		if p == 21 || p == 23 || p == 25 {
			t.Errorf("suppressed port %d should not appear in result", p)
		}
	}
}
