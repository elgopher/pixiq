package opengl_test

import (
	"testing"

	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func BenchmarkWindow_Draw(b *testing.B) {
	b.StopTimer()
	openGL := opengl.New(mainThreadLoop)
	win := openGL.Windows().Open(640, 360)
	images := pixiq.NewImages(openGL.AcceleratedImages())
	screenImage := images.New(640, 360)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		win.Draw(screenImage)
	}
}
