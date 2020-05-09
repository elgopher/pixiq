// Package goimage provides functions for converting Pixiq's image.Image into standard Go image.Image and vice-versa
package goimage

import (
	"github.com/jacekolszak/pixiq/image"
	stdimage "image"
)

// FromSelection creates a standard Go Image from Pixiq's image.Selection
func FromSelection(source image.Selection, options ...Option) stdimage.Image {
	return stdimage.NewRGBA(stdimage.Rect(0, 0, source.Width(), source.Height()))
}

func CopyFromSelection(source image.Selection, target stdimage.Image) {
}

type Option func()

func Zoom(zoom int) Option {
	return func() {

	}
}

// CopyToSelection copies standard Go Image to Pixiq's image.Selection
func CopyToSelection(source stdimage.Image, target image.Selection) {

}
