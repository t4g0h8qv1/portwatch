package portprofile_test

import (
	"testing"

	"github.com/example/portwatch/internal/portprofile"
)

func BenchmarkRegister(b *testing.B) {
	r := portprofile.New()
	ports := []int{80, 443, 8080, 8443}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Register("web", ports)
	}
}

func BenchmarkGet_Hit(b *testing.B) {
	r := portprofile.Default()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Get("web")
	}
}

func BenchmarkGet_Miss(b *testing.B) {
	r := portprofile.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Get("missing")
	}
}
