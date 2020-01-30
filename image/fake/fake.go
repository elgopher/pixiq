package fake

import (
	"github.com/jacekolszak/pixiq/image"
)

// NewAcceleratedImage returns a new instance of *AcceleratedImage
func NewAcceleratedImage(imageWidth, imageHeight int) *AcceleratedImage {
	img := &AcceleratedImage{
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
		pixels:      make([]image.Color, imageWidth*imageHeight),
	}
	return img
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	pixels      []image.Color
	imageWidth  int
	imageHeight int
}

// Upload send pixels to a container in RAM
func (i *AcceleratedImage) Upload(pixels []image.Color) {
	i.pixels = make([]image.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	copy(i.pixels, pixels)
}

// Download fills output slice with image colors
func (i *AcceleratedImage) Download(output []image.Color) {
	for j := 0; j < len(output); j++ {
		output[j] = i.pixels[j]
	}
}
