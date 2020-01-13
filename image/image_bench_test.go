package image_test

import (
	"testing"

	image "github.com/jacekolszak/pixiq/image"
)

func BenchmarkSelection_SetColor(b *testing.B) {
	b.StopTimer()
	var (
		color     = image.RGBA(10, 20, 30, 40)
		img       = image.New(1920, 1080, &fakeAcceleratedImage{})
		selection = img.WholeImageSelection()
		height    = selection.Height()
		width     = selection.Width()
	)
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
	var (
		img       = image.New(1920, 1080, &fakeAcceleratedImage{})
		selection = img.WholeImageSelection()
		height    = selection.Height()
		width     = selection.Width()
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.Color(x, y)
			}
		}
	}
}
