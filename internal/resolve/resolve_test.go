package resolve_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/resolve"
)

func newFakeResolver(addrs []string, err error, ttl time.Duration) *resolve.Resolver {
	r := resolve.New(ttl)
	r.SetLookup(func(_ context.Context, _ string) ([]string, error) {
		return addrs, err
	})
	return r
}

func TestResolve_IPPassthrough(t *testing.T) {
	r := resolve.New(0)
	addrs, err := r.Resolve(context.Background(), "192.168.1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != 1 || addrs[0] != "192.168.1.1" {
		t.Fatalf("expected [192.168.1.1], got %v", addrs)
	}
}

func TestResolve_UsesLookup(t *testing.T) {
	r := newFakeResolver([]string{"10.0.0.1", "10.0.0.2"}, nil, 0)
	addrs, err := r.Resolve(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != 2 {
		t.Fatalf("expected 2 addrs, got %d", len(addrs))
	}
}

func TestResolve_LookupError(t *testing.T) {
	r := newFakeResolver(nil, errors.New("dns failure"), 0)
	_, err := r.Resolve(context.Background(), "bad.host")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolve_CachesResult(t *testing.T) {
	calls := 0
	r := resolve.New(time.Minute)
	r.SetLookup(func(_ context.Context, _ string) ([]string, error) {
		calls++
		return []string{"1.2.3.4"}, nil
	})
	for i := 0; i < 3; i++ {
		r.Resolve(context.Background(), "example.com") //nolint:errcheck
	}
	if calls != 1 {
		t.Fatalf("expected 1 lookup call, got %d", calls)
	}
}

func TestResolve_CacheExpires(t *testing.T) {
	calls := 0
	r := resolve.New(50 * time.Millisecond)
	r.SetLookup(func(_ context.Context, _ string) ([]string, error) {
		calls++
		return []string{"1.2.3.4"}, nil
	})
	r.Resolve(context.Background(), "example.com") //nolint:errcheck
	time.Sleep(100 * time.Millisecond)
	r.Resolve(context.Background(), "example.com") //nolint:errcheck
	if calls != 2 {
		t.Fatalf("expected 2 lookup calls after TTL expiry, got %d", calls)
	}
}

func TestResolve_InvalidateClears(t *testing.T) {
	calls := 0
	r := resolve.New(time.Minute)
	r.SetLookup(func(_ context.Context, _ string) ([]string, error) {
		calls++
		return []string{"1.2.3.4"}, nil
	})
	r.Resolve(context.Background(), "example.com") //nolint:errcheck
	r.Invalidate("example.com")
	r.Resolve(context.Background(), "example.com") //nolint:errcheck
	if calls != 2 {
		t.Fatalf("expected 2 lookup calls after invalidate, got %d", calls)
	}
}

func TestFirst_ReturnsFirstAddr(t *testing.T) {
	r := newFakeResolver([]string{"5.5.5.5", "6.6.6.6"}, nil, 0)
	addr, err := r.First(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "5.5.5.5" {
		t.Fatalf("expected 5.5.5.5, got %s", addr)
	}
}

func TestFirst_EmptyAddrs(t *testing.T) {
	r := newFakeResolver([]string{}, nil, 0)
	_, err := r.First(context.Background(), "empty.host")
	if err == nil {
		t.Fatal("expected error for empty addr list")
	}
}
