package blender_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glsl/blender"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/opengl"
)

var mainThreadLoop *opengl.MainThreadLoop

func TestMain(m *testing.M) {
	var exit int
	opengl.StartMainThreadLoop(func(main *opengl.MainThreadLoop) {
		mainThreadLoop = main
		exit = m.Run()
	})
	os.Exit(exit)

}

func TestIntegration_With_OpenGL_Package(t *testing.T) {
	t.Run("should compile ImageBlender", func(t *testing.T) {
		gl, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer gl.Destroy()
		imageBlender, err := blender.CompileImageBlender(gl)
		// then
		require.NoError(t, err)
		assert.NotNil(t, imageBlender)
	})
	t.Run("should copy colors", func(t *testing.T) {
		gl, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer gl.Destroy()
		imageBlender, err := blender.CompileImageBlender(gl)
		require.NoError(t, err)
		var (
			source          = gl.NewImage(1, 1)
			target          = gl.NewImage(1, 1)
			white           = image.RGB(255, 255, 255)
			black           = image.RGB(0, 0, 0)
			sourceSelection = source.WholeImageSelection()
			targetSelection = target.WholeImageSelection()
		)
		sourceSelection.SetColor(0, 0, white)
		targetSelection.SetColor(0, 0, black)
		// when
		imageBlender.Blend(sourceSelection, targetSelection)
		// then
		assert.Equal(t, white, targetSelection.Color(0, 0))
	})
}
