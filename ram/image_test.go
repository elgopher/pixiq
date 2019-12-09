package ram_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/color"
	"github.com/jacekolszak/pixiq/ram"
)

var black = color.RGBA(0, 0, 0, 255)

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
}

func TestSelection_Line(t *testing.T) {
	t.Run("should return line", func(t *testing.T) {
		img := ram.NewImage(4, 3)
		selection := img.Selection(1, 0, 2, 2)
		_ = selection.Line(1)
	})
}

func TestLine_Set(t *testing.T) {

	t.Run("should set color", func(t *testing.T) {
		img := ram.NewImage(2, 2)
		selection := img.Selection(0, 0, 2, 2)
		line := selection.Line(1)
		// when
		line.Set(0, black)
		// then
		assert.Equal(t, black, line.Get(0))
	})
}
