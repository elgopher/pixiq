package glblend

import (
	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/image"
)

const vertexShaderSrc = `
#version 330 core
	
layout(location = 0) in vec2 xy;
layout(location = 1) in vec2 st;
out vec2 interpolatedST;

void main() {
	gl_Position = vec4(xy, 0.0, 1.0);
	interpolatedST = st;
}
`

const fragmentShaderSrc = `
#version 330 core

uniform sampler2D tex;
in vec2 interpolatedST;
out vec4 color;

void main() {
	color = texture(tex, interpolatedST);
}
`

func NewSource(context *gl.Context) (*Source, error) {
	vertexShader, err := context.CompileVertexShader(vertexShaderSrc)
	if err != nil {
		return nil, err
	}
	fragmentShader, err := context.CompileFragmentShader(fragmentShaderSrc)
	if err != nil {
		return nil, err
	}
	program, err := context.LinkProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}
	vertexBuffer := makeVertexBuffer(context)
	vertexArray := makeVertexArray(context, vertexBuffer)
	command := program.AcceleratedCommand(&blendCommand{vertexArray: vertexArray})
	return &Source{command: command}, nil
}

func makeVertexArray(context *gl.Context, buffer *gl.FloatVertexBuffer) *gl.VertexArray {
	array := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec2})
	xy := gl.VertexBufferPointer{Offset: 0, Stride: 4, Buffer: buffer}
	array.Set(0, xy)
	st := gl.VertexBufferPointer{Offset: 2, Stride: 4, Buffer: buffer}
	array.Set(1, st)
	return array
}

func makeVertexBuffer(gl *gl.Context) *gl.FloatVertexBuffer {
	// xy -> st
	vertices := []float32{
		-1, 1, 0, 1, // top-left
		1, 1, 1, 1, // top-right
		1, -1, 1, 0, // bottom-right
		-1, -1, 0, 0, // bottom-left
	}
	buffer := gl.NewFloatVertexBuffer(len(vertices))
	buffer.Upload(0, vertices)
	return buffer
}

type blendCommand struct {
	vertexArray *gl.VertexArray
}

func (c *blendCommand) RunGL(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
	renderer.BindTexture(0, "tex", selections[0].Image)
	renderer.DrawArrays(c.vertexArray, gl.TriangleFan, 0, 4)
}

type Source struct {
	command *gl.AcceleratedCommand
}

func (o *Source) BlendSourceToTarget(source image.Selection, target image.Selection) {
	target = target.WithSize(source.Width(), source.Height())
	target.Modify(o.command, source)
}
