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

// NewSource creates a new blending tool which replaces target selection with source
// colors. It is like coping of source selection colors into target.
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
	vertexBuffer := context.NewFloatVertexBuffer(16) // FIXME The buffer should not be not static
	vertexArray := makeVertexArray(context, vertexBuffer)
	command := program.AcceleratedCommand(&blendCommand{vertexBuffer: vertexBuffer, vertexArray: vertexArray})
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

type blendCommand struct {
	vertexBuffer *gl.FloatVertexBuffer
	vertexArray  *gl.VertexArray
}

func (c *blendCommand) RunGL(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
	source := selections[0]
	renderer.BindTexture(0, "tex", source.Image)
	var (
		imageWidth  = float32(source.Image.Width())
		left        = float32(source.Location.X) / imageWidth
		right       = float32(source.Location.X+source.Location.Width) / imageWidth
		imageHeight = float32(source.Image.Height())
		top         = (imageHeight - float32(source.Location.Y)) / imageHeight
		bottom      = (imageHeight - float32(source.Location.Y) - float32(source.Location.Height)) / imageHeight
	)
	// xy -> st
	vertices := []float32{
		-1, 1, left, top,
		1, 1, right, top,
		1, -1, right, bottom,
		-1, -1, left, bottom,
	}
	c.vertexBuffer.Upload(0, vertices)
	renderer.DrawArrays(c.vertexArray, gl.TriangleFan, 0, 4)
}

// Source is a blending tool which replaces target selection with source
// colors. It is like coping of source selection colors into target.
type Source struct {
	command *gl.AcceleratedCommand
}

// BlendSourceToTarget blends source into target selection.
// Only position of the target Selection is used and the source is not clamped by
// the target size.
func (s *Source) BlendSourceToTarget(source image.Selection, target image.Selection) {
	target = target.WithSize(source.Width(), source.Height())
	target.Modify(s.command, source) // is it fast? is it better to use the whole texture as target? (and update xy in the vertextbuffer)
}
