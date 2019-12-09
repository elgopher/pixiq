package color_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/color"
)

func TestRGBA(t *testing.T) {
	t.Run("should create color from 4 components", func(t *testing.T) {
		c := color.RGBA(0, 1, 2, 3)
		assert.Equal(t, c.R(), 0)
		assert.Equal(t, c.G(), 1)
		assert.Equal(t, c.B(), 2)
		assert.Equal(t, c.A(), 3)
	})
}
