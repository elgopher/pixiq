package glfw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
)

func BenchmarkAcceleratedImage_Upload(b *testing.B) {
	openGL, err := glfw.NewOpenGL(mainThreadLoop)
	require.NoError(b, err)
	defer openGL.Destroy()
	if err != nil {
		panic(err)
	}
	img := openGL.Context().NewAcceleratedImage(1, 1)
	pixels := []image.Color{image.Transparent, image.Transparent}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.Upload(pixels)
	}
}
