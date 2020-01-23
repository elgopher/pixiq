package blender

import (
	"github.com/jacekolszak/pixiq/glsl/program"
	"github.com/jacekolszak/pixiq/image"
)

type GL interface {
	DrawProgram() program.Draw
}

func CompileImageBlender(gl GL) (*ImageBlender, error) {
	if gl == nil {
		panic("nil drawer")
	}
	prog := program.New(gl.DrawProgram())
	prog.AddSelectionParameter("source")

	vertexShader := program.NewVertexShader()
	vertexShader.SetMain("color = sampleSelection(tex, position);")
	prog.SetVertexShader(vertexShader)

	fragmentShader := program.NewFragmentShader()
	fragmentShader.SetMain("gl_Position = vec4(vertexPosition, 0.0, 1.0); position = texturePosition;")
	prog.SetFragmentShader(fragmentShader)

	compiledProgram, _ := prog.Compille()

	vertexFormat := program.VertexFormat{}
	vertexFormat.AddFloat2("vertexPosition")
	vertexFormat.AddFloat2("texturePosition")
	compiledProgram.SetVertexFormat(vertexFormat)

	// TODO Here we can verify that blending really works on this platform
	// by executing blend. If it does not we can return error. This is much
	// better than panicking during real Blending
	return &ImageBlender{program: compiledProgram}, nil
}

type ImageBlender struct {
	program *program.CompiledProgram
}

func (b *ImageBlender) Blend(source, target image.Selection) {
	call := b.program.NewCall()
	buffer := program.VertexBuffer{}
	buffer.AddFloat2(1, 2)
	buffer.AddFloat2(1, 2)

	call.SetVertexBuffer(buffer)
	call.SetSelection("source", source)
	target.Modify(call)
}
