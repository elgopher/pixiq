package fake

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

// NewAcceleratedImage returns a new instance of *AcceleratedImage
func NewAcceleratedImage(width, height int) (*AcceleratedImage, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	img := &AcceleratedImage{
		imageWidth:  width,
		imageHeight: height,
		pixels:      make([]image.Color, width*height),
	}
	return img, nil
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
