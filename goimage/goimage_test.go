package goimage_test

import (
	stdimage "image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/goimage"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
)

func TestFromSelection(t *testing.T) {

	t.Run("should create standard Go image of Selection size clamped by image size", func(t *testing.T) {
		img := image.New(fake.NewAcceleratedImage(2, 2))
		tests := map[string]struct {
			source         image.Selection
			opts           []goimage.Option // not used - TODO TEST FOR ZOOM
			expectedWidth  int
			expectedHeight int
		}{
			"0,0": {
				source:         img.Selection(0, 0).WithSize(0, 0),
				expectedWidth:  0,
				expectedHeight: 0,
			},
			"1,2": {
				source:         img.Selection(0, 0).WithSize(1, 2),
				expectedWidth:  1,
				expectedHeight: 2,
			},
			"2,1": {
				source:         img.Selection(0, 0).WithSize(2, 1),
				expectedWidth:  2,
				expectedHeight: 1,
			},
			"3,2": {
				source:         img.Selection(0, 0).WithSize(3, 2),
				expectedWidth:  2,
				expectedHeight: 2,
			},
			"2,3": {
				source:         img.Selection(0, 0).WithSize(2, 3),
				expectedWidth:  2,
				expectedHeight: 2,
			},
			"-1,1": {
				source:         img.Selection(0, 0).WithSize(-1, 1),
				expectedWidth:  0,
				expectedHeight: 1,
			},
			"selection xy negative": {
				source:         img.Selection(-1, -1).WithSize(2, 2),
				expectedWidth:  1,
				expectedHeight: 1,
			},
			"zoom 2": {
				source:         img.Selection(0, 0).WithSize(1, 2),
				opts:           []goimage.Option{goimage.Zoom(2)},
				expectedWidth:  2,
				expectedHeight: 4,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				stdImg := goimage.FromSelection(test.source, test.opts...)
				// then
				expectedBounds := stdimage.Rectangle{
					Max: stdimage.Point{
						X: test.expectedWidth,
						Y: test.expectedHeight,
					},
				}
				assert.Equal(t, expectedBounds, stdImg.Bounds())
			})
		}
	})

	t.Run("should create standard Go image with copied pixels", func(t *testing.T) {
		img := image.New(fake.NewAcceleratedImage(2, 2))
		selection := img.WholeImageSelection()
		pix00 := color.RGBA{R: 0, G: 10, B: 20, A: 30}
		selection.SetColor(0, 0, image.RGBA(0, 10, 20, 30))
		pix10 := color.RGBA{R: 100, G: 110, B: 120, A: 130}
		selection.SetColor(1, 0, image.RGBA(100, 110, 120, 130))
		pix01 := color.RGBA{R: 150, G: 160, B: 170, A: 180}
		selection.SetColor(0, 1, image.RGBA(150, 160, 170, 180))
		pix11 := color.RGBA{R: 200, G: 210, B: 220, A: 230}
		selection.SetColor(1, 1, image.RGBA(200, 210, 220, 230))
		tests := map[string]struct {
			source         image.Selection
			opts           []goimage.Option
			expectedPixels [][]color.RGBA
		}{
			"top-left": {
				source: img.Selection(0, 0).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{
						pix00,
					},
				},
			},
			"bottom-right": {
				source: img.Selection(1, 1).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{
						pix11,
					},
				},
			},
			"top": {
				source: img.Selection(0, 0).WithSize(2, 1),
				expectedPixels: [][]color.RGBA{
					{
						pix00, pix10,
					},
				},
			},
			"left": {
				source: img.Selection(0, 0).WithSize(1, 2),
				expectedPixels: [][]color.RGBA{
					{
						pix00,
					},
					{
						pix01,
					},
				},
			},
			"whole": {
				source: img.WholeImageSelection(),
				expectedPixels: [][]color.RGBA{
					{
						pix00, pix10,
					},
					{
						pix01, pix11,
					},
				},
			},
			"zoom top-left": {
				source: img.Selection(0, 0).WithSize(1, 1),
				opts:   []goimage.Option{goimage.Zoom(2)},
				expectedPixels: [][]color.RGBA{
					{
						pix00, pix00,
					},
					{
						pix00, pix00,
					},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				stdImg := goimage.FromSelection(test.source, test.opts...)
				assertPixels(t, stdImg, test.expectedPixels)
			})
		}
	})

}

func TestFillWithSelection(t *testing.T) {
	t.Run("should fill standard Go image with pixels from Selection", func(t *testing.T) {
		img := image.New(fake.NewAcceleratedImage(2, 2))
		selection := img.WholeImageSelection()
		pix00 := color.RGBA{R: 0, G: 10, B: 20, A: 30}
		selection.SetColor(0, 0, image.RGBA(0, 10, 20, 30))
		pix10 := color.RGBA{R: 100, G: 110, B: 120, A: 130}
		selection.SetColor(1, 0, image.RGBA(100, 110, 120, 130))
		pix01 := color.RGBA{R: 150, G: 160, B: 170, A: 180}
		selection.SetColor(0, 1, image.RGBA(150, 160, 170, 180))
		pix11 := color.RGBA{R: 200, G: 210, B: 220, A: 230}
		selection.SetColor(1, 1, image.RGBA(200, 210, 220, 230))
		tests := map[string]struct {
			target         draw.Image
			source         image.Selection
			opts           []goimage.Option
			expectedPixels [][]color.RGBA
		}{
			"top-left": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1)),
				source: img.Selection(0, 0).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{
						pix00,
					},
				},
			},
			"bottom-right": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1)),
				source: img.Selection(1, 1).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{
						pix11,
					},
				},
			},
			"top": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 2, 1)),
				source: img.Selection(0, 0).WithSize(2, 1),
				expectedPixels: [][]color.RGBA{
					{
						pix00, pix10,
					},
				},
			},
			"left": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 2)),
				source: img.Selection(0, 0).WithSize(1, 2),
				expectedPixels: [][]color.RGBA{
					{
						pix00,
					},
					{
						pix01,
					},
				},
			},
			"whole": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 2, 2)),
				source: img.WholeImageSelection(),
				expectedPixels: [][]color.RGBA{
					{
						pix00, pix10,
					},
					{
						pix01, pix11,
					},
				},
			},
			"zoom top-left": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 2, 2)),
				source: img.Selection(0, 0).WithSize(1, 1),
				opts:   []goimage.Option{goimage.Zoom(2)},
				expectedPixels: [][]color.RGBA{
					{
						pix00, pix00,
					},
					{
						pix00, pix00,
					},
				},
			},
			"image smaller than selection": {
				target: stdimage.NewRGBA(stdimage.Rect(0, 0, 1, 1)),
				source: img.WholeImageSelection(),
				expectedPixels: [][]color.RGBA{
					{
						pix00,
					},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				targetImage := test.target
				// when
				goimage.FillWithSelection(targetImage, test.source, test.opts...)
				assertPixels(t, targetImage, test.expectedPixels)
			})
		}
	})

}

func assertPixels(t *testing.T, img stdimage.Image, pixels [][]color.RGBA) {
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[y]); x++ {
			pixel := img.At(x, y)
			expectedPixel := pixels[y][x]
			assert.Equal(t, expectedPixel, pixel)
		}
	}
}
