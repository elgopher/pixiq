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
	textureLocation := r.locationOrPanic(uniformAttributeName)
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

// SetFloat sets uniform attribute of type float
func (r *Renderer) SetFloat(uniformAttributeName string, value float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform1f(location, value)
	})
}

func (r *Renderer) locationOrPanic(uniformAttributeName string) int32 {
	trimmed := strings.TrimSpace(uniformAttributeName)
	if trimmed == "" {
		panic("empty uniformAttributeName")
	}
	location, ok := r.program.uniformLocations[uniformAttributeName]
	if !ok {
		panic("not existing uniform attribute name: " + uniformAttributeName)
	}
	return location
}

// SetVec2 sets uniform attribute of type vec2
func (r *Renderer) SetVec2(uniformAttributeName string, v1, v2 float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform2f(location, v1, v2)
	})
}

// SetVec3 sets uniform attribute of type vec3
func (r *Renderer) SetVec3(uniformAttributeName string, v1, v2, v3 float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform3f(location, v1, v2, v3)
	})
}

// SetVec4 sets uniform attribute of type vec4
func (r *Renderer) SetVec4(uniformAttributeName string, v1, v2, v3, v4 float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform4f(location, v1, v2, v3, v4)
	})
}

// SetInt sets uniform attribute of type int32
func (r *Renderer) SetInt(uniformAttributeName string, value int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform1i(location, value)
	})
}

// SetIVec2 sets uniform attribute of type ivec2
func (r *Renderer) SetIVec2(uniformAttributeName string, v1, v2 int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform2i(location, v1, v2)
	})
}

// SetIVec3 sets uniform attribute of type ivec3
func (r *Renderer) SetIVec3(uniformAttributeName string, v1, v2, v3 int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform3i(location, v1, v2, v3)
	})
}

// SetIVec4 sets uniform attribute of type ivec4
func (r *Renderer) SetIVec4(uniformAttributeName string, v1, v2, v3, v4 int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.Uniform4i(location, v1, v2, v3, v4)
	})
}

// SetMat3 sets uniform attribute of type mat3
func (r *Renderer) SetMat3(uniformAttributeName string, value [9]float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.UniformMatrix3fv(location, 1, false, &value[0])
	})
}

// SetMat4 sets uniform attribute of type mat4
func (r *Renderer) SetMat4(uniformAttributeName string, value [16]float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.runInOpenGLThread(func() {
		gl.UniformMatrix4fv(location, 1, false, &value[0])
	})
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
		locHeight := loc.Height
		if locHeight > img.height {
			locHeight = img.height
		}
		x := int32(loc.X)
		y := int32(img.height - locHeight - loc.Y)
		w := int32(loc.Width)
		h := int32(img.height - loc.Y)
		gl.Scissor(x, y, w, h)
		gl.Viewport(x, y, w, h) // TODO Not well tested
	})
	renderer := &Renderer{
		program:           c.program,
		runInOpenGLThread: c.runInOpenGLThread,
		allImages:         c.allImages,
	}
	c.command.RunGL(renderer, selections)
}
