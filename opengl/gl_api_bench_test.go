package opengl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/opengl"
)

// Must be at most 1 allocs/op
func BenchmarkContext_Clear(b *testing.B) {
	openGL, err := opengl.New(mainThreadLoop)
	require.NoError(b, err)
	defer openGL.Destroy()
	gl := openGL.ContextAPI()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gl.Clear(0x4000)
	}
	mainThreadLoop.Execute(func() {}) // wait until all commands are processed
}
