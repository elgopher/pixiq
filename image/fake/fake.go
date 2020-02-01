package fake

import (
	"errors"

	"github.com/jacekolszak/pixiq/image"
)

// NewAcceleratedImage returns a new instance of *AcceleratedImage which can be
// used in unit tests.
//
// It is a fake implementation of image.AcceleratedImage which stores
// pixel colors in RAM.
func NewAcceleratedImage(width, height int) (*AcceleratedImage, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	img := &AcceleratedImage{
		width:  width,
		height: height,
		Pixels: make([]image.Color, width*height),
	}
	return img, nil
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	// Hide the instance variable
	Pixels []image.Color
	width  int
	height int
}

// Upload send pixels to a container in RAM
func (i *AcceleratedImage) Upload(pixels []image.Color) {
	if len(pixels) != i.width*i.height {
		panic("pixels slice is not of length width*height")
	}
	// copy pixels to ensure that Upload method has been called
	copy(i.Pixels, pixels)
}

// Download fills output slice with image colors
func (i *AcceleratedImage) Download(output []image.Color) {
	if len(output) != i.width*i.height {
		panic("output slice is not of length width*height")
	}
	for j := 0; j < len(output); j++ {
		output[j] = i.Pixels[j]
	}
}
