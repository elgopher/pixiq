// Package goimage provides functions for converting Pixiq image.Selection
// into standard Go image.Image and vice-versa
package goimage

import (
	stdimage "image"

	"github.com/jacekolszak/pixiq/image"
)

// FromSelection creates a standard Go Image from Pixiq image.Selection
func FromSelection(source image.Selection, options ...Option) stdimage.Image {
	return stdimage.NewRGBA(stdimage.Rect(0, 0, source.Width(), source.Height()))
}

func FillWithSelection(target stdimage.Image, source image.Selection) {
}

type Option func()

func Zoom(zoom int) Option {
	return func() {

	}
}

// CopyToSelection copies standard Go Image to Pixiq image.Selection
func CopyToSelection(source stdimage.Image, target image.Selection) {

}
