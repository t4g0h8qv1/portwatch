package timeout_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/timeout"
)

func BenchmarkGet_NoOverride(b *testing.B) {
	m, _ := timeout.New(2 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("host")
	}
}

func BenchmarkGet_WithOverride(b *testing.B) {
	m, _ := timeout.New(2 * time.Second)
	_ = m.Set("host", 5*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("host")
	}
}

func BenchmarkSet_ManyHosts(b *testing.B) {
	m, _ := timeout.New(2 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Set(fmt.Sprintf("host-%d", i), time.Second)
	}
}
