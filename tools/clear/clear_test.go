package clear_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
	"github.com/jacekolszak/pixiq/tools/clear"
)

func TestNew(t *testing.T) {
	t.Run("should create tool", func(t *testing.T) {
		tool := clear.New()
		assert.NotNil(t, tool)
	})
}

func TestClear(t *testing.T) {
	t.Run("should clear selection", func(t *testing.T) {
		tool := clear.New()
		color := image.RGBA(1, 2, 3, 4)
		img := image.New(1, 1, fake.NewAcceleratedImage(1, 1))
		selection := img.WholeImageSelection()
		tool.SetColor(color)
		// when
		tool.Clear(selection)
		// then
		assert.Equal(t, color, selection.Color(0, 0))
	})
}
