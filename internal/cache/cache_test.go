package cache_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cache"
)

func TestNew_InvalidTTL(t *testing.T) {
	_, err := cache.New(0)
	if err == nil {
		t.Fatal("expected error for zero TTL, got nil")
	}
	_, err = cache.New(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative TTL, got nil")
	}
}

func TestSet_And_Get(t *testing.T) {
	c, err := cache.New(5 * time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ports := []int{22, 80, 443}
	c.Set("localhost", ports)

	entry, ok := c.Get("localhost")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if len(entry.Ports) != len(ports) {
		t.Fatalf("ports length: got %d, want %d", len(entry.Ports), len(ports))
	}
}

func TestGet_Miss(t *testing.T) {
	c, _ := cache.New(time.Second)
	_, ok := c.Get("nothere")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestGet_Expired(t *testing.T) {
	c, _ := cache.New(10 * time.Millisecond)
	c.Set("host", []int{8080})
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("host")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestInvalidate(t *testing.T) {
	c, _ := cache.New(time.Minute)
	c.Set("host", []int{22})
	c.Invalidate("host")
	_, ok := c.Get("host")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestPrune_RemovesExpired(t *testing.T) {
	c, _ := cache.New(10 * time.Millisecond)
	c.Set("a", []int{22})
	c.Set("b", []int{80})
	time.Sleep(20 * time.Millisecond)
	c.Set("c", []int{443}) // fresh entry

	removed := c.Prune()
	if removed != 2 {
		t.Fatalf("Prune removed %d entries, want 2", removed)
	}
	_, ok := c.Get("c")
	if !ok {
		t.Fatal("fresh entry should survive prune")
	}
}

func TestEntry_IsExpired(t *testing.T) {
	now := time.Now()
	expired := cache.Entry{ExpiresAt: now.Add(-time.Second)}
	if !expired.IsExpired() {
		t.Fatal("expected entry to be expired")
	}
	fresh := cache.Entry{ExpiresAt: now.Add(time.Minute)}
	if fresh.IsExpired() {
		t.Fatal("expected entry to be fresh")
	}
}
