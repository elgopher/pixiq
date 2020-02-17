// +build nvidia

package opengl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/opengl"
)

func TestOpenGL_Error_Nvidia(t *testing.T) {
	t.Run("should return out-of-memory error for too big vertex buffer", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		petabyte := 1024 * 1024 * 1024 * 1024 * 1024
		openGL.NewFloatVertexBuffer(petabyte)
		// when
		err := openGL.Error()
		// then
		require.Error(t, err)
		assert.True(t, opengl.IsOutOfMemory(err))
	})
}
