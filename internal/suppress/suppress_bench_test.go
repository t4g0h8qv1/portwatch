package suppress_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/suppress"
)

// BenchmarkIsSuppressed measures lookup performance with a large list.
func BenchmarkIsSuppressed(b *testing.B) {
	l, _ := suppress.Load(b.TempDir() + "/s.json")
	for i := 0; i < 500; i++ {
		_ = l.Add(i+1, fmt.Sprintf("port %d", i+1), time.Hour)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.IsSuppressed(250)
	}
}

// BenchmarkFilter measures filtering performance over a typical port slice.
func BenchmarkFilter(b *testing.B) {
	l, _ := suppress.Load(b.TempDir() + "/s.json")
	for i := 0; i < 100; i++ {
		_ = l.Add(i+1, "bench", time.Hour)
	}

	input := make([]int, 1024)
	for i := range input {
		input[i] = i + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Filter(input)
	}
}
