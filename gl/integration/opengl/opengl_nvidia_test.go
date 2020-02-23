// +build nvidia

package opengl_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/gl"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/opengl"
)

func TestContext_Error_Nvidia(t *testing.T) {
	t.Run("should return out-of-memory error for too big vertex buffer", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.NewContext(openGL)
		petabyte := 1024 * 1024 * 1024 * 1024 * 1024
		context.NewFloatVertexBuffer(petabyte)
		// when
		err := context.Error()
		// then
		require.Error(t, err)
		assert.True(t, gl.IsOutOfMemory(err))
	})
}
