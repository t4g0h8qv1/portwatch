package resolve_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/resolve"
)

func BenchmarkResolve_Cached(b *testing.B) {
	r := resolve.New(time.Hour)
	r.SetLookup(func(_ context.Context, _ string) ([]string, error) {
		return []string{"127.0.0.1"}, nil
	})
	// warm the cache
	r.Resolve(context.Background(), "bench.host") //nolint:errcheck

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Resolve(context.Background(), "bench.host") //nolint:errcheck
	}
}

func BenchmarkResolve_IPPassthrough(b *testing.B) {
	r := resolve.New(0)
	for i := 0; i < b.N; i++ {
		r.Resolve(context.Background(), "192.168.0.1") //nolint:errcheck
	}
}
