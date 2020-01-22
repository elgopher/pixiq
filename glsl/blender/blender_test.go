package blender_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glsl/blender"
	"github.com/jacekolszak/pixiq/glsl/shader"
)

func TestCompileImageBlender(t *testing.T) {
	t.Run("should compile ImageBlender", func(t *testing.T) {
		compiler := &fragmentCompilerStub{}
		// when
		imageBlender, err := blender.CompileImageBlender(compiler)
		// then
		require.NoError(t, err)
		assert.NotNil(t, imageBlender)
	})
	t.Run("should panic when compiler is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = blender.CompileImageBlender(nil)
		})
	})
}

type fragmentCompilerStub struct {
}

func (f fragmentCompilerStub) DrawTriangles() shader.GLProgram {
	panic("implement me")
}
