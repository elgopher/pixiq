package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
)

func TestRGBA(t *testing.T) {
	t.Run("should create colors using 4 components", func(t *testing.T) {
		color := image.RGBA(10, 20, 30, 40)
		assert.Equal(t, byte(10), color.R())
		assert.Equal(t, byte(20), color.G())
		assert.Equal(t, byte(30), color.B())
		assert.Equal(t, byte(40), color.A())
	})
}

func TestRGB(t *testing.T) {
	t.Run("should create colors using 3 components", func(t *testing.T) {
		color := image.RGB(10, 20, 30)
		assert.Equal(t, byte(10), color.R())
		assert.Equal(t, byte(20), color.G())
		assert.Equal(t, byte(30), color.B())
		assert.Equal(t, byte(255), color.A())
	})
}
func TestRGBAi(t *testing.T) {

	t.Run("RGBAi should clamp components to [0-255] range", func(t *testing.T) {
		tests := map[string]struct {
			given image.Color
			rgba  []byte
		}{
			"pixiq.RGBAi(-1, 20, 30, 40)": {
				given: image.RGBAi(-1, 20, 30, 40),
				rgba:  []byte{0, 20, 30, 40},
			},
			"pixiq.RGBAi(256, 20, 30, 40)": {
				given: image.RGBAi(256, 20, 30, 40),
				rgba:  []byte{255, 20, 30, 40},
			},
			"pixiq.RGBAi(10, -1, 30, 40)": {
				given: image.RGBAi(10, -1, 30, 40),
				rgba:  []byte{10, 0, 30, 40},
			},
			"pixiq.RGBAi(10, 256, 30, 40)": {
				given: image.RGBAi(10, 256, 30, 40),
				rgba:  []byte{10, 255, 30, 40},
			},
			"pixiq.RGBAi(10, 20, -1, 40)": {
				given: image.RGBAi(10, 20, -1, 40),
				rgba:  []byte{10, 20, 0, 40},
			},
			"pixiq.RGBAi(10, 20, 256, 40)": {
				given: image.RGBAi(10, 20, 256, 40),
				rgba:  []byte{10, 20, 255, 40},
			},
			"pixiq.RGBAi(10, 20, 30, -1)": {
				given: image.RGBAi(10, 20, 30, -1),
				rgba:  []byte{10, 20, 30, 0},
			},
			"pixiq.RGBAi(10, 20, 30, 256)": {
				given: image.RGBAi(10, 20, 30, 256),
				rgba:  []byte{10, 20, 30, 255},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, test.rgba[0], test.given.R())
				assert.Equal(t, test.rgba[1], test.given.G())
				assert.Equal(t, test.rgba[2], test.given.B())
				assert.Equal(t, test.rgba[3], test.given.A())
			})
		}

	})

	t.Run("should create color using 4 components given as integer numbers", func(t *testing.T) {
		color := image.RGBAi(10, 20, 30, 40)
		assert.Equal(t, byte(10), color.R())
		assert.Equal(t, byte(20), color.G())
		assert.Equal(t, byte(30), color.B())
		assert.Equal(t, byte(40), color.A())
	})

}

func TestColor_RGBAf(t *testing.T) {
	t.Run("should convert to floats", func(t *testing.T) {
		tests := map[string]struct {
			color         image.Color
			expectedRed   float32
			expectedGreen float32
			expectedBlue  float32
			expectedAlpha float32
		}{
			"0.0, 0.2, 0.4, 0.6": {
				color:         image.RGBA(0, 51, 102, 153),
				expectedRed:   0.0,
				expectedGreen: 0.2,
				expectedBlue:  0.4,
				expectedAlpha: 0.6,
			},
			"0.8, 0.898, 0.95, 1": {
				color:         image.RGBA(204, 229, 242, 255),
				expectedRed:   0.8,
				expectedGreen: 0.898,
				expectedBlue:  0.95,
				expectedAlpha: 1.0,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				r, g, b, a := test.color.RGBAf()
				const delta = 0.00196 // smaller than 0.5/255
				assert.InDelta(t, test.expectedRed, r, delta)
				assert.InDelta(t, test.expectedGreen, g, delta)
				assert.InDelta(t, test.expectedBlue, b, delta)
				assert.InDelta(t, test.expectedAlpha, a, delta)
			})
		}
	})
}
