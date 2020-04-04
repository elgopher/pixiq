package blend_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
	"github.com/jacekolszak/pixiq/tools/blend"
)

func TestNew(t *testing.T) {
	t.Run("should panic when colorBlender is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			blend.New(nil)
		})
	})
	t.Run("should create tool", func(t *testing.T) {
		tool := blend.New(multiplyColors{})
		assert.NotNil(t, tool)
	})
}

func TestTool_BlendSourceToTarget(t *testing.T) {
	t.Run("should skip blending when source image has 0 size", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
		}{
			"0 height": {
				width: 1,
			},
			"0 width": {
				height: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				source := image.New(test.width, test.height, fake.NewAcceleratedImage(test.width, test.height)).WholeImageSelection()
				target := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).WholeImageSelection()
				originalColor := image.RGBA(1, 2, 3, 4)
				target.SetColor(0, 0, originalColor)
				tool := blend.New(multiplyColors{})
				// when
				tool.BlendSourceToTarget(source, target)
				// then
				assert.Equal(t, originalColor, target.Color(0, 0))
			})
		}
	})
	t.Run("should blend whole source when target is 0 size", func(t *testing.T) {
		source := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).WholeImageSelection()
		target := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).Selection(0, 0)
		source.SetColor(0, 0, image.RGBA(1, 2, 3, 4))
		target.SetColor(0, 0, image.RGBA(5, 6, 7, 8))
		tool := blend.New(multiplyColors{})
		// when
		tool.BlendSourceToTarget(source, target)
		// then
		// target == source
		assert.Equal(t, image.RGBA(5, 12, 21, 32), target.Color(0, 0))
	})
	t.Run("should blend selections", func(t *testing.T) {
		source := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).WholeImageSelection()
		target := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).WholeImageSelection()
		source.SetColor(0, 0, image.RGBA(1, 2, 3, 4))
		target.SetColor(0, 0, image.RGBA(5, 6, 7, 8))
		tool := blend.New(multiplyColors{})
		// when
		tool.BlendSourceToTarget(source, target)
		// then
		assert.Equal(t, image.RGBA(5, 12, 21, 32), target.Color(0, 0))
		// and (TODO this move to some other test)
		assert.Equal(t, image.RGBA(1, 2, 3, 4), source.Color(0, 0))
	})
}

type multiplyColors struct{}

func (c multiplyColors) BlendSourceToTargetColor(source, target image.Color) image.Color {
	return image.RGBA(
		source.R()*target.R(),
		source.G()*target.G(),
		source.B()*target.B(),
		source.A()*target.A())
}
