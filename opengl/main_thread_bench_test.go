package opengl_test

import (
	"testing"
)

// Must be 0 allocs/op
func BenchmarkMainThreadLoop_Execute(b *testing.B) {
	f := func() {}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mainThreadLoop.Execute(f)
	}
}

// Must be 0 allocs/op and super fast (up to 100ns/op)
func BenchmarkMainThreadLoop_ExecuteAsync(b *testing.B) {
	f := func() {}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mainThreadLoop.ExecuteAsync(f)
	}
	mainThreadLoop.Execute(f) // wait until all commands are processed
}
