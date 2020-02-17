package opengl

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/jacekolszak/pixiq/image"
)

// Command is a procedure drawing primitives (such as triangles) in the AcceleratedImage.
type Command interface {
	// Implementations must not retain renderer and selections.
	RunGL(renderer *Renderer, selections []image.AcceleratedImageSelection)
}

// Renderer is an API for drawing primitives
type Renderer struct {
	program           *Program
	runInOpenGLThread func(func())
	allImages         allImages
}

// BindTexture assigns image.AcceleratedImage to a given textureUnit and uniform attribute.
// The bounded texture can be sampled in a fragment shader.
func (r *Renderer) BindTexture(textureUnit int, uniformAttributeName string, image image.AcceleratedImage) {
	if textureUnit < 0 {
		panic("negative textureUnit")
	}
	assertValidAttributeName(uniformAttributeName)
	textureLocation, ok := r.program.uniformLocations[uniformAttributeName]
	if !ok {
		panic("not existing uniform attribute name")
	}
	img, ok := r.allImages[image]
	if !ok {
		panic("image has not been created in this OpenGL context")
	}
	r.runInOpenGLThread(func() {
		gl.Uniform1i(textureLocation, int32(textureUnit))
		gl.ActiveTexture(uint32(gl.TEXTURE0 + textureUnit))
		gl.BindTexture(gl.TEXTURE_2D, img.textureID)
	})
}

func (r *Renderer) SetFloat(uniformAttributeName string, value float32) {
	assertValidAttributeName(uniformAttributeName)
	location, ok := r.program.uniformLocations[uniformAttributeName]
	if !ok {
		panic("not existing uniform attribute name")
	}
	r.runInOpenGLThread(func() {
		gl.Uniform1f(location, value)
	})
}

func assertValidAttributeName(uniformAttributeName string) {
	trimmed := strings.TrimSpace(uniformAttributeName)
	if trimmed == "" {
		panic("empty uniformAttributeName")
	}
}

// Mode defines which primitives will be drawn.
//
// See https://www.khronos.org/opengl/wiki/Primitive
type Mode struct {
	glMode uint32
}

var (
	// Points draws points using GL_POINTS.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Point_primitives
	Points = Mode{glMode: gl.POINTS}
	// LineStrip draws lines using GL_LINE_STRIP.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Line_primitives
	LineStrip = Mode{glMode: gl.LINE_STRIP}
	// LineLoop draws lines using GL_LINE_LOOP.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Line_primitives
	LineLoop = Mode{glMode: gl.LINE_LOOP}
	// Lines draws lines using GL_LINES.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Line_primitives
	Lines = Mode{glMode: gl.LINES}
	// TriangleStrip draws triangles using GL_TRIANGLE_STRIP.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Triangle_primitives
	TriangleStrip = Mode{glMode: gl.TRIANGLE_STRIP}
	// TriangleFan draws triangles using GL_TRIANGLE_FAN.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Triangle_primitives
	TriangleFan = Mode{glMode: gl.TRIANGLE_FAN}
	// Triangles draws triangles using GL_TRIANGLES.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Triangle_primitives
	Triangles = Mode{glMode: gl.TRIANGLES}
)

// DrawArrays draws primitives (such as triangles) using vertices defined in VertexArray.
//
// Before primitive is drawn this method validates if
func (r *Renderer) DrawArrays(array *VertexArray, mode Mode, first, count int) {
	r.validateAttributeTypes(array)
	r.runInOpenGLThread(func() {
		gl.BindVertexArray(array.id)
		gl.DrawArrays(mode.glMode, int32(first), int32(count))
	})
}

func (r *Renderer) validateAttributeTypes(array *VertexArray) {
	for i := 0; i < len(array.layout); i++ {
		if attr, ok := r.program.attributes[int32(i)]; ok {
			vertexArrayType := array.layout[i]
			if attr.typ != vertexArrayType {
				err := fmt.Sprintf("shader attribute %s with location %d has type %v, which is different than %v in the vertex array", attr.name, i, attr.typ, vertexArrayType)
				panic(err)
			}
		}
	}
}

// Clear clears the selection with a given color.
func (r *Renderer) Clear(color image.Color) {
	r.runInOpenGLThread(func() {
		gl.ClearColor(color.RGBAf())
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

// AcceleratedCommand is an image.AcceleratedCommand implementation. It delegates
// the drawing to Command.
type AcceleratedCommand struct {
	command           Command
	program           *Program
	runInOpenGLThread func(func())
	allImages         allImages
}

// Run implements image.AcceleratedCommand#Run.
func (c *AcceleratedCommand) Run(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) {
	if c.command == nil {
		return
	}
	if output.Image == nil {
		panic("nil output Image")
	}
	img, ok := c.allImages[output.Image]
	if !ok {
		panic("output image created in a different OpenGL context than program")
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
	c.command.RunGL(renderer, selections)
}
