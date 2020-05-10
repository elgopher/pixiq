// Package goimage provides functions for converting Pixiq image.Selection
// into standard Go image.Image and vice-versa
package goimage

import (
	stdimage "image"
	"image/color"
	"image/draw"

	"github.com/jacekolszak/pixiq/image"
)

// FromSelection creates a standard Go Image from Pixiq image.Selection
func FromSelection(source image.Selection, options ...Option) stdimage.Image {
	opts := buildOpts(source, options...)
	if opts.height == 0 {
		return stdimage.NewRGBA(stdimage.Rectangle{})
	}
	bounds := stdimage.Rect(0, 0, opts.width*opts.zoom, opts.height*opts.zoom)
	target := stdimage.NewRGBA(bounds)
	FillWithSelection(target, source, options...)
	return target
}

type opts struct {
	zoom          int
	width, height int
}

func buildOpts(source image.Selection, options ...Option) opts {
	opts := opts{
		zoom: 1,
	}
	for _, option := range options {
		opts = option(opts)
	}
	lines := source.Lines()
	opts.height = lines.Length()
	if opts.height != 0 {
		opts.width = len(lines.LineForRead(0))
	}
	return opts
}

func FillWithSelection(target draw.Image, source image.Selection, options ...Option) {
	opts := buildOpts(source, options...)
	lines := source.Lines()
	for y := 0; y < opts.height; y++ {
		for x := 0; x < opts.width; x++ {
			line := lines.LineForRead(y)
			c := line[x]
			rgba := color.RGBA{R: c.R(), G: c.G(), B: c.B(), A: c.A()}
			for zy := 0; zy < opts.zoom; zy++ {
				for zx := 0; zx < opts.zoom; zx++ {
					target.Set(x*opts.zoom+zx, y*opts.zoom+zy, rgba)
				}
			}
		}
	}
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
