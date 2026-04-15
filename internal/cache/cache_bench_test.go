package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/cache"
)

func BenchmarkSet(b *testing.B) {
	c, _ := cache.New(time.Minute)
	ports := []int{22, 80, 443, 8080, 8443}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(fmt.Sprintf("host-%d", i%100), ports)
	}
}

func BenchmarkGet_Hit(b *testing.B) {
	c, _ := cache.New(time.Minute)
	c.Set("bench-host", []int{22, 80, 443})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("bench-host")
	}
}

func BenchmarkPrune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		c, _ := cache.New(10 * time.Millisecond)
		for j := 0; j < 500; j++ {
			c.Set(fmt.Sprintf("host-%d", j), []int{80})
		}
		time.Sleep(15 * time.Millisecond)
		b.StartTimer()
		c.Prune()
	}
}
