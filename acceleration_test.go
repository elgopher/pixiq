package pixiq_test

import (
	"github.com/jacekolszak/pixiq"
)

type fakeAcceleratedImages struct {
	images []*fakeAcceleratedImage
}

func (i *fakeAcceleratedImages) New(width, height int) pixiq.AcceleratedImage {
	image := &fakeAcceleratedImage{
		width:  width,
		height: height,
	}
	i.images = append(i.images, image)
	return image
}

type fakeAcceleratedImage struct {
	pixels []pixiq.Color
	width  int
	height int
}

func (a *fakeAcceleratedImage) Upload(pixels []pixiq.Color) {
	a.pixels = make([]pixiq.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	for i, pixel := range pixels {
		a.pixels[i] = pixel
	}
}

func (a *fakeAcceleratedImage) Download(output []pixiq.Color) {
	for i := 0; i < len(output); i++ {
		output[i] = a.pixels[i]
	}
}
