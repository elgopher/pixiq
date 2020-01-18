package blender_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glsl/blender"
	"github.com/jacekolszak/pixiq/glsl/shader"
	"github.com/jacekolszak/pixiq/image"
)

func TestCompileImageBlender(t *testing.T) {
	t.Run("should compile ImageBlender", func(t *testing.T) {
		compiler := &compilerStub{}
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

type compilerStub struct {
}

func (c *compilerStub) Compile(fragmentShaderSource string) (shader.Program, error) {
	return &compiledShaderStub{programCallStub: &programCallStub{}}, nil
}

type compiledShaderStub struct {
	*programCallStub
}

func (c *compiledShaderStub) Call() shader.ProgramCall {
	return &programCallStub{}
}

type programCallStub struct {
}

func (p programCallStub) SetTexture(uniformName string, selection image.Selection) {

}

func (p programCallStub) Release() {
}
