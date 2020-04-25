package fake_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
)

func TestNewAcceleratedImage(t *testing.T) {
	t.Run("should panic when width<0", func(t *testing.T) {
		assert.Panics(t, func() {
			fake.NewAcceleratedImage(-1, 1)
		})
	})
	t.Run("should panic when height<0", func(t *testing.T) {
		assert.Panics(t, func() {
			fake.NewAcceleratedImage(1, -1)
		})
	})
	t.Run("should create image", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
		}{
			"0x0": {width: 0, height: 0},
			"1x2": {width: 1, height: 2},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := fake.NewAcceleratedImage(test.width, test.height)
				assert.NotNil(t, img)
				assert.Equal(t, test.width, img.Width())
				assert.Equal(t, test.height, img.Height())
			})
		}
	})
}

func TestAcceleratedImage_Download(t *testing.T) {
	white := image.RGB(255, 255, 255)
	t.Run("should download not uploaded image", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
			given         []image.Color
			expected      []image.Color
		}{
			"0x0": {
				width: 0, height: 0,
				given:    []image.Color{},
				expected: []image.Color{},
			},
			"1x2": {
				width: 1, height: 2,
				given:    []image.Color{white, white},
				expected: []image.Color{image.Transparent, image.Transparent},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := fake.NewAcceleratedImage(test.width, test.height)
				output := make([]image.Color, len(test.expected))
				// when
				img.Download(output)
				// then
				assert.Equal(t, test.expected, output)
			})
		}
	})
	t.Run("should panic when output slice is not of width*height length", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
			inputLength   int
		}{
			"1x1, input length 0": {
				width: 1, height: 1,
				inputLength: 0,
			},
			"1x1, input length 2": {
				width: 1, height: 1,
				inputLength: 2,
			},
			"1x2, input length 1": {
				width: 1, height: 2,
				inputLength: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := fake.NewAcceleratedImage(test.width, test.height)
				output := make([]image.Color, test.inputLength)
				assert.Panics(t, func() {
					img.Download(output)
				})
			})
		}
	})
}

func TestAcceleratedImage_Upload(t *testing.T) {
	t.Run("should panic when input slice is not of width*height length", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
			inputLength   int
		}{
			"1x1, input length 0": {
				width: 1, height: 1,
				inputLength: 0,
			},
			"1x1, input length 2": {
				width: 1, height: 1,
				inputLength: 2,
			},
			"1x2, input length 1": {
				width: 1, height: 2,
				inputLength: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := fake.NewAcceleratedImage(test.width, test.height)
				input := make([]image.Color, test.inputLength)
				assert.Panics(t, func() {
					img.Upload(input)
				})
			})
		}
	})
	t.Run("should upload colors", func(t *testing.T) {
		var (
			color0 = image.RGB(0, 0, 0)
			color1 = image.RGB(1, 1, 1)
		)
		tests := map[string]struct {
			width, height int
			colors        []image.Color
		}{
			"0x0": {
				width: 0, height: 0,
				colors: []image.Color{},
			},
			"1x1": {
				width: 1, height: 1,
				colors: []image.Color{color0},
			},
			"1x2": {
				width: 1, height: 2,
				colors: []image.Color{color0, color1},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := fake.NewAcceleratedImage(test.width, test.height)
				// when
				img.Upload(test.colors)
				// then
				output := make([]image.Color, len(test.colors))
				img.Download(output)
				assert.Equal(t, test.colors, output)
			})
		}
	})
	t.Run("should copy colors", func(t *testing.T) {
		color0 := image.RGB(0, 0, 0)
		color1 := image.RGB(1, 1, 1)
		img := fake.NewAcceleratedImage(1, 1)
		input := []image.Color{color0}
		// when
		img.Upload(input)
		input[0] = color1
		// then
		output := make([]image.Color, 1)
		img.Download(output)
		assert.Equal(t, []image.Color{color0}, output)
	})
}

func TestAcceleratedImage_PixelsTable(t *testing.T) {
	t.Run("should return 2d slice", func(t *testing.T) {
		color0 := image.RGB(0, 0, 0)
		color1 := image.RGB(1, 1, 1)
		tests := map[string]struct {
			width, height int
			pixels        []image.Color
			expectedTable [][]image.Color
		}{
			"0x0": {
				width: 0, height: 0,
				pixels:        []image.Color{},
				expectedTable: [][]image.Color{},
			},
			"1x1": {
				width: 1, height: 1,
				pixels:        []image.Color{color0},
				expectedTable: [][]image.Color{{color0}},
			},
			"1x2": {
				width: 1, height: 2,
				pixels: []image.Color{color0, color1},
				expectedTable: [][]image.Color{
					{color0},
					{color1},
				},
			},
			"2x1": {
				width: 2, height: 1,
				pixels: []image.Color{color0, color1},
				expectedTable: [][]image.Color{
					{color0, color1},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				acceleratedImage := fake.NewAcceleratedImage(test.width, test.height)
				acceleratedImage.Upload(test.pixels)
				// when
				table := acceleratedImage.PixelsTable()
				// then
				assert.Equal(t, test.expectedTable, table)
			})
		}
	})
}
