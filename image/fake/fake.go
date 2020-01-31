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
		pixels: make([]image.Color, width*height),
	}
	return img, nil
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	pixels []image.Color
	width  int
	height int
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

// TODO Experimental (not tested, POC)
func (i *AcceleratedImage) PixelTable() [][]image.Color {
	table := make([][]image.Color, i.height)
	for y := 0; y < i.height; y++ {
		table[y] = make([]image.Color, i.width)
		for x := 0; x < i.width; x++ {
			table[y][x] = i.pixels[y*i.width+x]
		}
	}
	return table
}
