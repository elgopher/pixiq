package pixiq_test

import (
	"testing"

	"github.com/jacekolszak/pixiq"
)

func BenchmarkSelection_SetColor(b *testing.B) {
	b.StopTimer()
	color := pixiq.RGBA(10, 20, 30, 40)
	images := pixiq.NewImages(&fakeAcceleratedImages{})
	image := images.New(1920, 1080)
	selection := image.WholeImageSelection()
	height := selection.Height()
	width := selection.Width()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.SetColor(x, y, color)
			}
		}
	}
}

func BenchmarkSelection_Color(b *testing.B) {
	b.StopTimer()
	images := pixiq.NewImages(&fakeAcceleratedImages{})
	image := images.New(1920, 1080)
	selection := image.WholeImageSelection()
	height := selection.Height()
	width := selection.Width()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.Color(x, y)
			}
		}
	}
}
