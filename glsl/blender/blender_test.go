package blender_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glsl/blender"
	"github.com/jacekolszak/pixiq/glsl/program"
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

func (f fragmentCompilerStub) NewFloatVertexBuffer(program.BufferUsage) program.FloatVertexBuffer {
	panic("implement me")
}

func (f fragmentCompilerStub) DrawProgram() program.Draw {
	panic("implement me")
}
