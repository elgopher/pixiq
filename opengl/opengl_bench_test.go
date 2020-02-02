package opengl_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/opengl"
)

func BenchmarkTexture_Upload(b *testing.B) {
	openGL, err := opengl.New(mainThreadLoop)
	if err != nil {
		panic(err)
	}
	texture, err := openGL.NewTexture(1, 1)
	if err != nil {
		panic(err)
	}
	pixels := []image.Color{image.Transparent, image.Transparent}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		texture.Upload(pixels)
	}
}
