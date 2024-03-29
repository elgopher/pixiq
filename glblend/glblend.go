// Package glblend provides blend tools using video card
// TODO Add Blender with fragment shader as a parameter. Add methods for setting uniforms
package glblend

import (
	"github.com/elgopher/pixiq/gl"
	"github.com/elgopher/pixiq/image"
)

// NewSource creates a new blending tool which replaces target selection with source
// colors. It is like coping of source selection colors into target.
func NewSource(context *gl.Context) (*Source, error) {
	command, err := newBlendCommand(context, gl.SourceBlendFactors)
	if err != nil {
		return nil, err
	}
	return &Source{command: command}, nil
}

// NewSourceOver creates a new blending tool which blends together source and target
// taking into account alpha channel of both. Source-over means that source will be
// painted on top of the target.
func NewSourceOver(context *gl.Context) (*SourceOver, error) {
	command, err := newBlendCommand(context, gl.BlendFactors{
		SrcFactor: gl.One,
		DstFactor: gl.OneMinusSrcAlpha,
	})
	if err != nil {
		return nil, err
	}
	return &SourceOver{source: &Source{command: command}}, nil
}

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
	// color is blended with buffer using formula: S * sf + D * df 
	color = texture(tex, interpolatedST);
}
`

func newBlendCommand(context *gl.Context, factors gl.BlendFactors) (*gl.AcceleratedCommand, error) {
	if context == nil {
		panic("nil context")
	}
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
	vertexBuffer := context.NewFloatVertexBuffer(16, gl.DynamicDraw)
	vertexArray := makeVertexArray(context, vertexBuffer)
	command := program.AcceleratedCommand(
		&blendCommand{
			vertexBuffer: vertexBuffer,
			vertexArray:  vertexArray,
			factors:      factors,
		})
	return command, nil
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
	factors      gl.BlendFactors
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
	renderer.SetBlendFactors(c.factors)
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
	source = clampSourceToTargetImage(source, target)
	target = target.WithSize(source.Width(), source.Height())
	// FIXME is it fast enough? or is it better to use the whole texture as a target and update xy in the vertextbuffer accordingly?
	target.Modify(s.command, source)
}

func clampSourceToTargetImage(source image.Selection, target image.Selection) image.Selection {
	width := source.Width()
	if width+target.ImageX() > target.Image().Width() {
		width = target.Image().Width() - target.ImageX()
	}
	height := source.Height()
	if height+target.ImageY() > target.Image().Height() {
		height = target.Image().Height() - target.ImageY()
	}
	return source.WithSize(width, height)
}

// SourceOver (aka Normal) is a blending tool which blends together source and target
// taking into account alpha channel of both. Source-over means that source will be
// painted on top of the target.
type SourceOver struct {
	source *Source
}

// BlendSourceToTarget blends source into target selection.
// Only position of the target Selection is used and the source is not clamped by
// the target size.
func (s *SourceOver) BlendSourceToTarget(source image.Selection, target image.Selection) {
	s.source.BlendSourceToTarget(source, target)
}
