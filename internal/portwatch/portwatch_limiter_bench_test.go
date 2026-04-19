package portwatch

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkAllow_SingleTarget(b *testing.B) {
	l, _ := NewScanLimiter(time.Millisecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Allow("host1")
	}
}

func BenchmarkAllow_ManyTargets(b *testing.B) {
	l, _ := NewScanLimiter(time.Millisecond)
	targets := make([]string, 100)
	for i := range targets {
		targets[i] = fmt.Sprintf("host%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Allow(targets[i%len(targets)])
	}
}

func BenchmarkLastScan_Hit(b *testing.B) {
	l, _ := NewScanLimiter(time.Millisecond)
	_ = l.Allow("host1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = l.LastScan("host1")
	}
}
