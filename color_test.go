package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq"
)

func TestRGBA(t *testing.T) {
	t.Run("should create colors using 4 components", func(t *testing.T) {
		color := pixiq.RGBA(10, 20, 30, 40)
		assert.Equal(t, byte(10), color.R())
		assert.Equal(t, byte(20), color.G())
		assert.Equal(t, byte(30), color.B())
		assert.Equal(t, byte(40), color.A())
	})
}

func TestRGB(t *testing.T) {
	t.Run("should create colors using 3 components", func(t *testing.T) {
		color := pixiq.RGB(10, 20, 30)
		assert.Equal(t, byte(10), color.R())
		assert.Equal(t, byte(20), color.G())
		assert.Equal(t, byte(30), color.B())
		assert.Equal(t, byte(255), color.A())
	})
}
func TestRGBAi(t *testing.T) {

	t.Run("RGBAi should clamp components to [0-255] range", func(t *testing.T) {
		tests := map[string]struct {
			given pixiq.Color
			rgba  []byte
		}{
			"pixiq.RGBAi(-1, 20, 30, 40)": {
				given: pixiq.RGBAi(-1, 20, 30, 40),
				rgba:  []byte{0, 20, 30, 40},
			},
			"pixiq.RGBAi(256, 20, 30, 40)": {
				given: pixiq.RGBAi(256, 20, 30, 40),
				rgba:  []byte{255, 20, 30, 40},
			},
			"pixiq.RGBAi(10, -1, 30, 40)": {
				given: pixiq.RGBAi(10, -1, 30, 40),
				rgba:  []byte{10, 0, 30, 40},
			},
			"pixiq.RGBAi(10, 256, 30, 40)": {
				given: pixiq.RGBAi(10, 256, 30, 40),
				rgba:  []byte{10, 255, 30, 40},
			},
			"pixiq.RGBAi(10, 20, -1, 40)": {
				given: pixiq.RGBAi(10, 20, -1, 40),
				rgba:  []byte{10, 20, 0, 40},
			},
			"pixiq.RGBAi(10, 20, 256, 40)": {
				given: pixiq.RGBAi(10, 20, 256, 40),
				rgba:  []byte{10, 20, 255, 40},
			},
			"pixiq.RGBAi(10, 20, 30, -1)": {
				given: pixiq.RGBAi(10, 20, 30, -1),
				rgba:  []byte{10, 20, 30, 0},
			},
			"pixiq.RGBAi(10, 20, 30, 256)": {
				given: pixiq.RGBAi(10, 20, 30, 256),
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
		color := pixiq.RGBAi(10, 20, 30, 40)
		assert.Equal(t, byte(10), color.R())
		assert.Equal(t, byte(20), color.G())
		assert.Equal(t, byte(30), color.B())
		assert.Equal(t, byte(40), color.A())
	})

}
