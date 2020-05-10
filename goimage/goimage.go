// Package goimage provides functions for converting Pixiq image.Selection
// into standard Go image.Image and vice-versa
package goimage

import (
	stdimage "image"
	"image/color"

	"github.com/jacekolszak/pixiq/image"
)

// FromSelection creates a standard Go Image from Pixiq image.Selection
func FromSelection(source image.Selection, options ...Option) stdimage.Image {
	opts := opts{
		zoom: 1,
	}
	for _, option := range options {
		opts = option(opts)
	}
	lines := source.Lines()
	if lines.Length() == 0 {
		return stdimage.NewRGBA(stdimage.Rectangle{})
	}
	width := len(lines.LineForRead(0))
	height := lines.Length()
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, width*opts.zoom, height*opts.zoom))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			line := lines.LineForRead(y)
			c := line[x]
			rgba := color.RGBA{R: c.R(), G: c.G(), B: c.B(), A: c.A()}
			for zy := 0; zy < opts.zoom; zy++ {
				for zx := 0; zx < opts.zoom; zx++ {
					img.Set(x*opts.zoom+zx, y*opts.zoom+zy, rgba)
				}
			}
		}
	}
	return img
}

type opts struct {
	zoom int
}

func FillWithSelection(target stdimage.Image, source image.Selection, options ...Option) {
}

type Option func(opts) opts

func Zoom(zoom int) Option {
	return func(o opts) opts {
		o.zoom = zoom
		return o
	}
}

// CopyToSelection copies standard Go Image to Pixiq image.Selection
func CopyToSelection(source stdimage.Image, target image.Selection, options ...Option) {

}
