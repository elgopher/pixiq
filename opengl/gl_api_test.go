package opengl_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/opengl"
)

// Should be 0 allocs/op (but it's 2 allocs/op at the moment)
func BenchmarkContext_Clear(b *testing.B) {
	openGL, err := opengl.New(mainThreadLoop)
	if err != nil {
		panic(err)
	}
	defer openGL.Destroy()
	gl := openGL.Context().API()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gl.Clear(0x4000)
	}
}
