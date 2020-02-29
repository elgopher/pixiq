package glfw_test

import (
	"testing"
)

// Must be at most 2 allocs/op
func BenchmarkMainThreadLoop_Execute(b *testing.B) {
	f := func() {}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mainThreadLoop.Execute(f)
	}
}
