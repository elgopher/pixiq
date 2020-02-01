package fake_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
)

func TestNewAcceleratedImage(t *testing.T) {
	t.Run("should return error when width<0", func(t *testing.T) {
		img, err := fake.NewAcceleratedImage(-1, 1)
		assert.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should return error when height<0", func(t *testing.T) {
		img, err := fake.NewAcceleratedImage(1, -1)
		assert.Error(t, err)
		assert.Nil(t, img)
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
				img, err := fake.NewAcceleratedImage(test.width, test.height)
				require.NoError(t, err)
				assert.NotNil(t, img)
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
				img, _ := fake.NewAcceleratedImage(test.width, test.height)
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
				img, _ := fake.NewAcceleratedImage(test.width, test.height)
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
				img, _ := fake.NewAcceleratedImage(test.width, test.height)
				input := make([]image.Color, test.inputLength)
				assert.Panics(t, func() {
					img.Upload(input)
				})
			})
		}
	})
}
