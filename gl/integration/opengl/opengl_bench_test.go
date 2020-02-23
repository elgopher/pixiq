package opengl_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/opengl"
)

func BenchmarkAcceleratedImage_Upload(b *testing.B) {
	openGL, err := opengl.New(mainThreadLoop)
	if err != nil {
		panic(err)
	}
	img := openGL.Context().NewAcceleratedImage(1, 1)
	pixels := []image.Color{image.Transparent, image.Transparent}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.Upload(pixels)
	}
}
