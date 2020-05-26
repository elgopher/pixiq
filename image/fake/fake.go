// Package fake provides a fake image.AcceleratedImage implementation
// which can be used in unit testing.
package fake

import (
	"github.com/jacekolszak/pixiq/image"
)

// NewAcceleratedImage returns a new instance of *AcceleratedImage which can be
// used in unit tests.
//
// It is a fake implementation of image.AcceleratedImage which stores
// pixel colors in RAM.
func NewAcceleratedImage(width, height int) *AcceleratedImage {
	if width < 0 {
		panic("negative width")
	}
	if height < 0 {
		panic("negative height")
	}
	return &AcceleratedImage{
		width:  width,
		height: height,
		pixels: make([]image.Color, width*height),
	}
}

// AcceleratedImage stores pixel data in RAM and uses CPU solely.
type AcceleratedImage struct {
	// Hide the instance variable
	pixels  []image.Color
	width   int
	height  int
	deleted bool
}

// Upload send pixels to a container in RAM
func (i *AcceleratedImage) Upload(pixels []image.Color) {
	if len(pixels) != i.width*i.height {
		panic("pixels slice is not of length width*height")
	}
	// copy pixels to ensure that Upload method has been called
	copy(i.pixels, pixels)
}

// Download fills output slice with image colors
func (i *AcceleratedImage) Download(output []image.Color) {
	if len(output) != i.width*i.height {
		panic("output slice is not of length width*height")
	}
	for j := 0; j < len(output); j++ {
		output[j] = i.pixels[j]
	}
}

// Width returns the number of pixels in a row.
func (i *AcceleratedImage) Width() int {
	return i.width
}

// Height returns the number of pixels in a column.
func (i *AcceleratedImage) Height() int {
	return i.height
}

// PixelsTable returns a copy of pixels in a form of 2D slice: first dimension
// is a row, second dimension is a pixel in the row.
func (i *AcceleratedImage) PixelsTable() [][]image.Color {
	table := make([][]image.Color, i.height)
	for row := 0; row < i.height; row++ {
		table[row] = make([]image.Color, i.width)
		for cell := 0; cell < i.width; cell++ {
			idx := row*i.width + cell
			table[row][cell] = i.pixels[idx]
		}
	}
	return table
}

// Delete marks AcceleratedImage as deleted
func (i *AcceleratedImage) Delete() {
	i.deleted = true
}

// Deleted returns true if AcceleratedImage.Delete() method was called
func (i *AcceleratedImage) Deleted() bool {
	return i.deleted
}
