package portwatch_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/portwatch"
)

func BenchmarkRun_NoChanges(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "baseline.json")
	n := &fakeNotifier{}

	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{}, // empty port list — fast path
		BaselinePath: path,
		Notifier:     n,
	}

	// Seed baseline.
	if err := portwatch.Run(context.Background(), cfg); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := portwatch.Run(context.Background(), cfg); err != nil {
			b.Fatal(err)
		}
	}
}
