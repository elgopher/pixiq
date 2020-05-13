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
			opts           []goimage.Option
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
			"img width +1, 2": {
				source:         img.Selection(0, 0).WithSize(3, 2),
				expectedWidth:  img.Width() + 1,
				expectedHeight: 2,
			},
			"2, img height + 1": {
				source:         img.Selection(0, 0).WithSize(2, 3),
				expectedWidth:  2,
				expectedHeight: img.Height() + 1,
			},
			"-1,1": {
				source:         img.Selection(0, 0).WithSize(-1, 1),
				expectedWidth:  0,
				expectedHeight: 1,
			},
			"1,-1": {
				source:         img.Selection(0, 0).WithSize(1, -1),
				expectedWidth:  1,
				expectedHeight: 0,
			},
			"selection xy negative": {
				source:         img.Selection(-1, -1).WithSize(2, 2),
				expectedWidth:  2,
				expectedHeight: 2,
			},
			"zoom 2": {
				source:         img.Selection(0, 0).WithSize(1, 2),
				opts:           []goimage.Option{goimage.Zoom(2)},
				expectedWidth:  2,
				expectedHeight: 4,
			},
			"zoom 0": {
				source:         img.WholeImageSelection(),
				opts:           []goimage.Option{goimage.Zoom(0)},
				expectedWidth:  img.Width(),
				expectedHeight: img.Height(),
			},
			"zoom -1": {
				source:         img.WholeImageSelection(),
				opts:           []goimage.Option{goimage.Zoom(-1)},
				expectedWidth:  img.Width(),
				expectedHeight: img.Height(),
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
		transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}
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
			"negative x": {
				source: img.Selection(-1, 0).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{transparent},
				},
			},
			"negative y": {
				source: img.Selection(0, -1).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{transparent},
				},
			},
			"negative x, width 2": {
				source: img.Selection(-1, 0).WithSize(2, 1),
				expectedPixels: [][]color.RGBA{
					{transparent, pix00},
				},
			},
			"negative y, height 2": {
				source: img.Selection(0, -1).WithSize(1, 2),
				expectedPixels: [][]color.RGBA{
					{transparent},
					{pix00},
				},
			},
			"negative x,y": {
				source: img.Selection(-1, -2).WithSize(3, 4),
				expectedPixels: [][]color.RGBA{
					{transparent, transparent, transparent},
					{transparent, transparent, transparent},
					{transparent, pix00, pix10},
					{transparent, pix01, pix11},
				},
			},
			"width greater than image size": {
				source: img.Selection(0, 0).WithSize(3, 1),
				expectedPixels: [][]color.RGBA{
					{pix00, pix10, transparent},
				},
			},
			"height greater than image size": {
				source: img.Selection(0, 0).WithSize(1, 3),
				expectedPixels: [][]color.RGBA{
					{pix00},
					{pix01},
					{transparent},
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
		transparent := color.RGBA{R: 0, G: 0, B: 0, A: 0}
		pix00 := color.RGBA{R: 0, G: 10, B: 20, A: 30}
		selection.SetColor(0, 0, image.RGBA(0, 10, 20, 30))
		pix10 := color.RGBA{R: 100, G: 110, B: 120, A: 130}
		selection.SetColor(1, 0, image.RGBA(100, 110, 120, 130))
		pix01 := color.RGBA{R: 150, G: 160, B: 170, A: 180}
		selection.SetColor(0, 1, image.RGBA(150, 160, 170, 180))
		pix11 := color.RGBA{R: 200, G: 210, B: 220, A: 230}
		selection.SetColor(1, 1, image.RGBA(200, 210, 220, 230))
		red := color.RGBA{R: 255, G: 30, B: 20, A: 200}
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
			"negative x": {
				target: imageWithColor(red, 1, 1),
				source: img.Selection(-1, 0).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{transparent},
				},
			},
			"negative y": {
				target: imageWithColor(red, 1, 1),
				source: img.Selection(0, -1).WithSize(1, 1),
				expectedPixels: [][]color.RGBA{
					{transparent},
				},
			},
			"negative x, width 2": {
				target: imageWithColor(red, 2, 1),
				source: img.Selection(-1, 0).WithSize(2, 1),
				expectedPixels: [][]color.RGBA{
					{transparent, pix00},
				},
			},
			"negative y, height 2": {
				target: imageWithColor(red, 1, 2),
				source: img.Selection(0, -1).WithSize(1, 2),
				expectedPixels: [][]color.RGBA{
					{transparent},
					{pix00},
				},
			},
			"negative x,y": {
				target: imageWithColor(red, 3, 4),
				source: img.Selection(-1, -2).WithSize(3, 4),
				expectedPixels: [][]color.RGBA{
					{transparent, transparent, transparent},
					{transparent, transparent, transparent},
					{transparent, pix00, pix10},
					{transparent, pix01, pix11},
				},
			},
			"width greater than image size": {
				target: imageWithColor(red, 2, 1),
				source: img.Selection(0, 0).WithSize(3, 1),
				expectedPixels: [][]color.RGBA{
					{pix00, pix10, transparent},
				},
			},
			"height greater than image size": {
				target: imageWithColor(red, 1, 2),
				source: img.Selection(0, 0).WithSize(1, 3),
				expectedPixels: [][]color.RGBA{
					{pix00},
					{pix01},
					{transparent},
				},
			},
			"zoom 0": {
				target: imageWithColor(red, 1, 1),
				source: img.WholeImageSelection(),
				opts:   []goimage.Option{goimage.Zoom(0)},
				expectedPixels: [][]color.RGBA{
					{pix00},
				},
			},
			"zoom -1": {
				target: imageWithColor(red, 1, 1),
				source: img.WholeImageSelection(),
				opts:   []goimage.Option{goimage.Zoom(-1)},
				expectedPixels: [][]color.RGBA{
					{pix00},
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

func TestCopyToSelection(t *testing.T) {
	c00 := image.RGBA(0, 10, 20, 30)
	p00 := color.RGBA{R: 0, G: 10, B: 20, A: 30}
	c10 := image.RGBA(40, 50, 60, 70)
	p10 := color.RGBA{R: 40, G: 50, B: 60, A: 70}
	c01 := image.RGBA(80, 90, 100, 110)
	p01 := color.RGBA{R: 80, G: 90, B: 100, A: 110}
	c11 := image.RGBA(120, 130, 140, 150)
	p11 := color.RGBA{R: 120, G: 130, B: 140, A: 150}

	bounds := stdimage.Rect(0, 0, 2, 2)
	rgbaImage := stdimage.NewRGBA(bounds)
	rgbaImage.SetRGBA(0, 0, p00)
	rgbaImage.SetRGBA(1, 0, p10)
	rgbaImage.SetRGBA(0, 1, p01)
	rgbaImage.SetRGBA(1, 1, p11)

	tests := map[string]struct {
		targetImage     *image.Image
		targetSelection func(img *image.Image) image.Selection
		options         []goimage.Option
		expectedColors  [][]image.Color
	}{
		"0x0 image and selection": {
			targetImage: image.New(fake.NewAcceleratedImage(0, 0)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0)
			},
		},
		"1x1 image and 0x0 selection": {
			targetImage: image.New(fake.NewAcceleratedImage(1, 1)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0)
			},
			expectedColors: [][]image.Color{
				{image.Transparent},
			},
		},
		"top-left": {
			targetImage: image.New(fake.NewAcceleratedImage(1, 1)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(1, 1)
			},
			expectedColors: [][]image.Color{
				{c00},
			},
		},
		"top": {
			targetImage: image.New(fake.NewAcceleratedImage(2, 1)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(2, 1)
			},
			expectedColors: [][]image.Color{
				{c00, c10},
			},
		},
		"left": {
			targetImage: image.New(fake.NewAcceleratedImage(1, 2)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(1, 2)
			},
			expectedColors: [][]image.Color{
				{c00},
				{c01},
			},
		},
		"whole": {
			targetImage: image.New(fake.NewAcceleratedImage(2, 2)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(2, 2)
			},
			expectedColors: [][]image.Color{
				{c00, c10},
				{c01, c11},
			},
		},
		"negative selection start": {
			targetImage: image.New(fake.NewAcceleratedImage(1, 1)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(-1, -1).WithSize(2, 2)
			},
			expectedColors: [][]image.Color{
				{c11},
			},
		},
		"top-left, zoom 2": {
			targetImage: image.New(fake.NewAcceleratedImage(2, 2)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(2, 2)
			},
			options: []goimage.Option{goimage.Zoom(2)},
			expectedColors: [][]image.Color{
				{c00, c00},
				{c00, c00},
			},
		},
		"whole, zoom 2": {
			targetImage: image.New(fake.NewAcceleratedImage(4, 4)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(4, 4)
			},
			options: []goimage.Option{goimage.Zoom(2)},
			expectedColors: [][]image.Color{
				{c00, c00, c10, c10},
				{c00, c00, c10, c10},
				{c01, c01, c11, c11},
				{c01, c01, c11, c11},
			},
		},
		"shifted by 1/2 of the zoom": {
			targetImage: image.New(fake.NewAcceleratedImage(3, 3)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(1, 1).WithSize(2, 2)
			},
			options: []goimage.Option{goimage.Zoom(2)},
			expectedColors: [][]image.Color{
				{image.Transparent, image.Transparent, image.Transparent},
				{image.Transparent, c00, c00},
				{image.Transparent, c00, c00},
			},
		},
		"to small selection for zoom": {
			targetImage: image.New(fake.NewAcceleratedImage(2, 2)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(0, 0).WithSize(2, 2)
			},
			options: []goimage.Option{goimage.Zoom(2)},
			expectedColors: [][]image.Color{
				{c00, c00},
				{c00, c00},
			},
		},
		"to small selection for zoom and selection shifted": {
			targetImage: image.New(fake.NewAcceleratedImage(2, 2)),
			targetSelection: func(img *image.Image) image.Selection {
				return img.Selection(1, 1).WithSize(2, 2)
			},
			options: []goimage.Option{goimage.Zoom(2)},
			expectedColors: [][]image.Color{
				{image.Transparent, image.Transparent},
				{image.Transparent, c00},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			selection := test.targetSelection(test.targetImage)
			// when
			goimage.CopyToSelection(rgbaImage, selection, test.options...)
			// then
			assertColors(t, test.targetImage, test.expectedColors)
		})
	}

}

func assertColors(t *testing.T, img *image.Image, pixels [][]image.Color) {
	selection := img.WholeImageSelection()
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[y]); x++ {
			pixel := selection.Color(x, y)
			expectedPixel := pixels[y][x]
			assert.Equal(t, expectedPixel, pixel)
		}
	}
}

func imageWithColor(color color.RGBA, width, height int) *stdimage.RGBA {
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color)
		}
	}
	return img
}
