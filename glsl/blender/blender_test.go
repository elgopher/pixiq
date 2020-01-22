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

func (f *fragmentCompilerStub) DrawTriangles(vertices []float32, vertexShader shader.VertexShader, fragmentShader shader.FragmentShader) shader.Call {
	panic("implement me")
}

func (f *fragmentCompilerStub) CompileVertexShader(glsl string) (shader.VertexShader, error) {
	panic("implement me")
}

func (f *fragmentCompilerStub) CompileFragmentShader(glsl string) (shader.FragmentShader, error) {
	return &fragmentShaderStub{}, nil
}

type fragmentShaderStub struct {
}
