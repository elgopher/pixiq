package opengl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/jacekolszak/pixiq/image"
)

type Command interface {
	// Implementations must not retain renderer and selections.
	RunGL(renderer *Renderer, selections []image.AcceleratedImageSelection) error
}

type Renderer struct {
	program           *Program
	runInOpenGLThread func(func())
	allImages         allImages
}

func (r *Renderer) BindTexture(textureUnit int, uniformAttributeName string, image image.AcceleratedImage) error {
	if textureUnit < 0 {
		return errors.New("negative textureUnit")
	}
	trimmed := strings.TrimSpace(uniformAttributeName)
	if trimmed == "" {
		return errors.New("empty uniformAttributeName")
	}
	textureLocation, err := r.program.uniformAttributeLocation(uniformAttributeName)
	if err != nil {
		return err
	}
	img, ok := r.allImages[image]
	if !ok {
		return errors.New("image has not been created in this OpenGL context")
	}
	r.runInOpenGLThread(func() {
		gl.Uniform1i(textureLocation, int32(textureUnit))
		gl.ActiveTexture(uint32(gl.TEXTURE0 + textureUnit))
		gl.BindTexture(gl.TEXTURE_2D, img.textureID)
	})
	return nil
}

type Mode struct {
	glMode uint32
}

var (
	Points        = Mode{glMode: gl.POINTS}
	LineStrip     = Mode{glMode: gl.LINE_STRIP}
	LineLoop      = Mode{glMode: gl.LINE_LOOP}
	Lines         = Mode{glMode: gl.LINES}
	TriangleStrip = Mode{glMode: gl.TRIANGLE_STRIP}
	TriangleFan   = Mode{glMode: gl.TRIANGLE_FAN}
	Triangles     = Mode{glMode: gl.TRIANGLES}
)

func (r *Renderer) DrawArrays(array *VertexArray, mode Mode, first, count int) error {
	if err := r.validateAttributeTypes(array); err != nil {
		return err
	}
	r.runInOpenGLThread(func() {
		gl.BindVertexArray(array.id)
		gl.DrawArrays(mode.glMode, int32(first), int32(count))
	})
	return nil
}

func (r *Renderer) validateAttributeTypes(array *VertexArray) error {
	attributesLen := len(array.layout)
	if len(r.program.attributes) < attributesLen {
		attributesLen = len(r.program.attributes)
	}
	for i := 0; i < attributesLen; i++ {
		attr := r.program.attributes[i]
		vertexArrayType := array.layout[i]
		if attr.typ != vertexArrayType {
			return fmt.Errorf("shader attribute %s with location %d has type %v, which is different than %v in the vertex array", attr.name, i, attr.typ, vertexArrayType)
		}
	}
	return nil
}

func (r *Renderer) Clear(color image.Color) {
	r.runInOpenGLThread(func() {
		gl.ClearColor(color.RGBAf())
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

// AcceleratedCommand is an image.AcceleratedCommand implementation.
type AcceleratedCommand struct {
	command           Command
	program           *Program
	runInOpenGLThread func(func())
	allImages         allImages
}

func (c *AcceleratedCommand) Run(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) error {
	if output.Image == nil {
		return errors.New("nil output Image")
	}
	img, ok := c.allImages[output.Image]
	if !ok {
		return errors.New("output image created in a different OpenGL context than program")
	}
	c.runInOpenGLThread(func() {
		c.program.use()
		gl.Enable(gl.SCISSOR_TEST)
		gl.BindFramebuffer(gl.FRAMEBUFFER, img.frameBufferID)
		loc := output.Location
		gl.Scissor(int32(loc.X), int32(loc.Y), int32(loc.Width), int32(loc.Height))
		gl.Viewport(int32(loc.X), int32(loc.Y), int32(loc.Width), int32(loc.Height))
	})
	renderer := &Renderer{
		program:           c.program,
		runInOpenGLThread: c.runInOpenGLThread,
		allImages:         c.allImages,
	}
	return c.command.RunGL(renderer, selections)
}
