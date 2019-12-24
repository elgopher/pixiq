package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (i *fakeAcceleratedImages) assertOneImageWithPixels(t *testing.T, expectedPixels []pixiq.Color) {
	require.Len(t, i.images, 1)
	assert.Equal(t, expectedPixels, i.images[0].pixels)
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

var openWindowMock = func(width, height int) pixiq.SystemWindow {
	return &systemWindowMock{}
}

type systemWindowMock struct {
	imagesDrawn []pixiq.AcceleratedImage
}

func (f *systemWindowMock) Draw(image pixiq.AcceleratedImage) {
	f.imagesDrawn = append(f.imagesDrawn, image)
}
