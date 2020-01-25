package blender

import (
	"github.com/jacekolszak/pixiq/glsl/program"
	"github.com/jacekolszak/pixiq/image"
)

type GL interface {
	DrawProgram() program.Draw
	program.Buffers
}

func CompileImageBlender(gl GL) (*ImageBlender, error) {
	if gl == nil {
		panic("nil drawer")
	}

	lowLevelProg := gl.DrawProgram()
	lowLevelProg.SetVertexShader(`
		in vec2 vertexPosition;
		in vec2 texturePosition;
		
		out vec2 position;
		
		void main() {
			gl_Position = vec4(vertexPosition, 0.0, 1.0);
			position = texturePosition;
		}
		`)
	lowLevelProg.SetFragmentShader(`
		in vec2 position;
		
		out vec4 color;
		
		uniform sampler2D tex;
		
		void main() {
			color = texture(tex, position);
		}
		`)
	compiled, _ := lowLevelProg.Compile()

	vertexPosition := compiled.GetVertexAttributeLocation("vertexPosition")
	texturePosition := compiled.GetVertexAttributeLocation("texturePosition")

	buffer := gl.NewFloatVertexBuffer(program.StaticDraw)
	buffer.Update(0, []float32{
		-1, -1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
		1, -1, 1, 1,
		1, 1, 1, 0,
		-1, -1, 0, 1,
		1, 1, 1, 0,
		-1, 1, 0, 0,
	})

	vao := compiled.NewVertexArrayObject()
	vao.SetVertexAttribute(vertexPosition, buffer.Pointer(0, 2, 4))
	vao.SetVertexAttribute(texturePosition, buffer.Pointer(2, 2, 4))

	// high level:
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
	return &ImageBlender{
			program:  compiledProgram,
			compiled: compiled,
			vao:      vao,
		},
		nil
}

type ImageBlender struct {
	program  *program.CompiledProgram
	compiled program.CompiledDraw
	vao      program.VertexArrayObject
}

func (b *ImageBlender) Blend(source, target image.Selection) {
	// high-level
	call := b.program.NewCall(func(call program.HighLevelCall) {
		call.AddFloat2(1, 2) // TODO To slabe, bo jak niby zrobic to statycznie?
		call.AddFloat2(1, 2)
		call.SetSelection("source", source)
		call.Draw(program.Triangles, 0, 3)
		call.Draw(program.Triangles, 3, 3)
	})
	target.Modify(call)

	// or low-level
	lowLevelCall := b.compiled.NewCall(func(call program.DrawCall) {
		// TODO niebezpieczne troche bo closure wciaz ma dostep do zmiennych w funkcji otaczajacej
		call.BindVertexArrayObject(b.vao)
		call.BindTexture0(source.Image())
		call.Draw(program.Triangles, 0, 3)
		call.Draw(program.Triangles, 3, 3)
	})
	target.Modify(lowLevelCall)

	vertexBuffer := struct{}{}

	target.Modify2(b.program, func(call image.ProgramCall) {
		call.SetSelection("source", source) // tu nie cieknie abstrakcja wiec ta funkckja moze dzialac z roznymi technologiami
		call.SetFloat("a", 1.0)
		call.Draw(vertexBuffer, program.Triangles)
	})
}
