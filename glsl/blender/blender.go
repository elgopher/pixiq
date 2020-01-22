package blender

import (
	"github.com/jacekolszak/pixiq/glsl/shader"
	"github.com/jacekolszak/pixiq/image"
)

type GL interface {
	DrawTriangles() shader.GLProgram
}

func CompileImageBlender(gl GL) (*ImageBlender, error) {
	if gl == nil {
		panic("nil drawer")
	}
	program := shader.NewProgram(gl.DrawTriangles())
	program.AddSelectionUniform("source")

	vertexShader := shader.NewVertexShader()
	vertexShader.SetMain("color = sampleSelection(tex, position);")
	program.SetVertexShader(vertexShader)

	fragmentShader := shader.NewFragmentShader()
	fragmentShader.SetMain("gl_Position = vec4(vertexPosition, 0.0, 1.0); position = texturePosition;")
	program.SetFragmentShader(fragmentShader)

	compiledProgram, _ := program.Compille()

	// TODO Here we can verify that blending really works on this platform
	// by executing blend. If it does not we can return error. This is much
	// better than panicking during real Blending
	return &ImageBlender{program: compiledProgram}, nil
}

type ImageBlender struct {
	program *shader.CompiledProgram
}

func (b *ImageBlender) Blend(source, target image.Selection) {
	call := b.program.New()
	call.SetVertices([]float32{})
	call.SetSelection("source", source)
	target.Modify(call)
}
