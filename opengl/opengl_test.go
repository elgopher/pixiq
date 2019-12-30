package opengl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func TestRun(t *testing.T) {
	t.Run("should execute callback with opengl implementations of AcceleratedImages and SystemWindows", func(t *testing.T) {
		var images pixiq.AcceleratedImages
		var windows pixiq.SystemWindows
		// when
		opengl.Run(func(gl *opengl.OpenGL) {
			images = gl.AcceleratedImages()
			windows = gl.SystemWindows()
		})
		// then
		assert.NotNil(t, images)
		assert.NotNil(t, windows)
	})
}
