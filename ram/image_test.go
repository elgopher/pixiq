package ram_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/ram"
)

func TestNewImage(t *testing.T) {
	img := ram.NewImage(3, 2)
	assert.NotNil(t, img)
	assert.Equal(t, 3, img.Width())
	assert.Equal(t, 2, img.Height())
}

func TestImage_Selection(t *testing.T) {
	t.Run("should create selection", func(t *testing.T) {
		img := ram.NewImage(4, 4)
		selection := img.Selection(1, 0, 2, 3)
		assert.Equal(t, 1, selection.X())
		assert.Equal(t, 0, selection.Y())
		assert.Equal(t, 2, selection.Width())
		assert.Equal(t, 3, selection.Height())
	})
	t.Run("should return line", func(t *testing.T) {
		img := ram.NewImage(4, 3)
		selection := img.Selection(1, 0, 2, 2)
		_ = selection.Line(1)
	})
}
