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
	opts := buildOpts(options...)
	bounds := stdimage.Rect(0, 0, source.Width()*opts.zoom, source.Height()*opts.zoom)
	target := stdimage.NewRGBA(bounds)
	FillWithSelection(target, source, options...)
	return target
}

// Option is a conversion option
type Option func(opts) opts

// Zoom increases the image during conversion. Zoom <= 0 is treated as zoom 1.
func Zoom(zoom int) Option {
	return func(o opts) opts {
		if zoom > 0 {
			o.zoom = zoom
		} else {
			o.zoom = 1
		}
		return o
	}
}

type opts struct {
	zoom int
}

func buildOpts(options ...Option) opts {
	opts := opts{
		zoom: 1,
	}
	for _, option := range options {
		opts = option(opts)
	}
	return opts
}

// FillWithSelection fills existing standard Go draw.Image with pixels from image.Selection
func FillWithSelection(target draw.Image, source image.Selection, options ...Option) {
	opts := buildOpts(options...)
	for y := 0; y < source.Height()*opts.zoom; y++ {
		for x := 0; x < source.Width()*opts.zoom; x++ {
			c := source.Color(x/opts.zoom, y/opts.zoom)
			rgba := color.RGBA{R: c.R(), G: c.G(), B: c.B(), A: c.A()}
			target.Set(x, y, rgba)
		}
	}
}

// CopyToSelection copies standard Go Image to Pixiq image.Selection.
// The size of target Selection limits how much of the source Image is copied.
func CopyToSelection(source stdimage.Image, target image.Selection, options ...Option) {
	opts := buildOpts(options...)
	lines := target.Lines()
	if lines.Length() == 0 {
		return
	}
	width := len(target.Lines().LineForWrite(0))
	for y := 0; y < lines.Length(); y++ {
		line := lines.LineForWrite(y)
		for x := 0; x < width; x++ {
			pixel := source.At(x/opts.zoom+lines.XOffset(), y/opts.zoom+lines.YOffset())
			r, g, b, a := pixel.RGBA()
			line[x] = image.RGBA(byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}
}
