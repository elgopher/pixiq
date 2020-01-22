package blender_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glsl/blender"
	"github.com/jacekolszak/pixiq/opengl"
)

// should be 0 alloc
func BenchmarkImageBlender_Blend(b *testing.B) {
	b.StopTimer()
	gl, err := opengl.New(mainThreadLoop)
	require.NoError(b, err)
	imageBlender, err := blender.CompileImageBlender(gl)
	require.NoError(b, err)
	source := gl.NewImage(0, 0).WholeImageSelection()
	target := gl.NewImage(0, 0).WholeImageSelection()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// when
		imageBlender.Blend(source, target)
	}
}
