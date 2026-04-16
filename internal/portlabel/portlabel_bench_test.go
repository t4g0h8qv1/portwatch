package portlabel_test

import (
	"testing"

	"github.com/example/portwatch/internal/portlabel"
)

func BenchmarkResolve_BuiltIn(b *testing.B) {
	l := portlabel.New(nil)
	for i := 0; i < b.N; i++ {
		l.Resolve(443)
	}
}

func BenchmarkResolve_Unknown(b *testing.B) {
	l := portlabel.New(nil)
	for i := 0; i < b.N; i++ {
		l.Resolve(9999)
	}
}

func BenchmarkLabelAll(b *testing.B) {
	l := portlabel.New(nil)
	ports := []int{22, 80, 443, 3306, 5432, 6379, 8080, 9999}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.LabelAll(ports)
	}
}
