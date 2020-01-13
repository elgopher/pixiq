package image_test

import (
	"github.com/jacekolszak/pixiq/image"
)

type fakeAcceleratedImage struct {
	pixels []image.Color
	width  int
	height int
}

func (a *fakeAcceleratedImage) Upload(pixels []image.Color) {
	a.pixels = make([]image.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	for i, pixel := range pixels {
		a.pixels[i] = pixel
	}
}

func (a *fakeAcceleratedImage) Download(output []image.Color) {
	for i := 0; i < len(output); i++ {
		output[i] = a.pixels[i]
	}
}
