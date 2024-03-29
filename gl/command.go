package gl

import (
	"fmt"
	"strings"

	"github.com/elgopher/pixiq/image"
)

// Command is a procedure drawing primitives (such as triangles) in the AcceleratedImage.
type Command interface {
	// Implementations must not retain renderer and selections.
	RunGL(renderer *Renderer, selections []image.AcceleratedImageSelection)
}

// Renderer is an API for drawing primitives.
//
// It is using blending with formula:
//   R = S*sf + D*df
// where S is source color component, sf is a source factor, D is destination color
// component, df is a destination factor.
// By default sf is 1 and df is 0 (gl.SourceBlendFactors), which means that formula is
//   R = S
// BlendFactors can be changed by calling SetBlendFactors method.
type Renderer struct {
	program      *Program
	api          API
	allImages    allImages
	blendFactors BlendFactors
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
	r.api.Uniform1i(textureLocation, int32(textureUnit))
	r.api.ActiveTexture(uint32(texture0 + textureUnit))
	r.api.BindTexture(texture2D, img.textureID)
}

// SetFloat sets uniform attribute of type float
func (r *Renderer) SetFloat(uniformAttributeName string, value float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform1f(location, value)
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
	r.api.Uniform2f(location, v1, v2)
}

// SetVec3 sets uniform attribute of type vec3
func (r *Renderer) SetVec3(uniformAttributeName string, v1, v2, v3 float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform3f(location, v1, v2, v3)
}

// SetVec4 sets uniform attribute of type vec4
func (r *Renderer) SetVec4(uniformAttributeName string, v1, v2, v3, v4 float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform4f(location, v1, v2, v3, v4)
}

// SetInt sets uniform attribute of type int32
func (r *Renderer) SetInt(uniformAttributeName string, value int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform1i(location, value)
}

// SetIVec2 sets uniform attribute of type ivec2
func (r *Renderer) SetIVec2(uniformAttributeName string, v1, v2 int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform2i(location, v1, v2)
}

// SetIVec3 sets uniform attribute of type ivec3
func (r *Renderer) SetIVec3(uniformAttributeName string, v1, v2, v3 int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform3i(location, v1, v2, v3)
}

// SetIVec4 sets uniform attribute of type ivec4
func (r *Renderer) SetIVec4(uniformAttributeName string, v1, v2, v3, v4 int32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.Uniform4i(location, v1, v2, v3, v4)
}

// SetMat3 sets uniform attribute of type mat3
func (r *Renderer) SetMat3(uniformAttributeName string, value [9]float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.UniformMatrix3fv(location, 1, false, &value[0])
}

// SetMat4 sets uniform attribute of type mat4
func (r *Renderer) SetMat4(uniformAttributeName string, value [16]float32) {
	location := r.locationOrPanic(uniformAttributeName)
	r.api.UniformMatrix4fv(location, 1, false, &value[0])
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
	Points = Mode{glMode: points}
	// LineStrip draws lines using GL_LINE_STRIP.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Line_primitives
	LineStrip = Mode{glMode: lineStrip}
	// LineLoop draws lines using GL_LINE_LOOP.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Line_primitives
	LineLoop = Mode{glMode: lineLoop}
	// Lines draws lines using GL_LINES.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Line_primitives
	Lines = Mode{glMode: lines}
	// TriangleStrip draws triangles using GL_TRIANGLE_STRIP.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Triangle_primitives
	TriangleStrip = Mode{glMode: triangleStrip}
	// TriangleFan draws triangles using GL_TRIANGLE_FAN.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Triangle_primitives
	TriangleFan = Mode{glMode: triangleFan}
	// Triangles draws triangles using GL_TRIANGLES.
	//
	// See https://www.khronos.org/opengl/wiki/Primitive#Triangle_primitives
	Triangles = Mode{glMode: triangles}
)

// BlendFactor is source or destination factor
type BlendFactor uint32

const (
	// Zero is GL_ZERO. Multiplies all components by 0.
	Zero = BlendFactor(0)
	// One is GL_ONE. Multiplies all components by 1.
	One = BlendFactor(1)
	// OneMinusSrcAlpha is GL_ONE_MINUS_SRC_ALPHA. Multiplies all components by 1 minus
	// the source alpha value.
	OneMinusSrcAlpha = BlendFactor(0x0303)
	// SrcAlpha is GL_SRC_ALPHA. Multiplies all components by the source alpha value.
	SrcAlpha = BlendFactor(0x0302)
	// DstAlpha is GL_DST_ALPHA. Multiplies all components by the destination alpha value.
	DstAlpha = BlendFactor(0x0304)
	// OneMinusDstAlpha is GL_ONE_MINUS_DST_ALPHA. Multiplies all components by 1 minus
	// the destination alpha value.
	OneMinusDstAlpha = BlendFactor(0x0305)
)

// BlendFactors contains source and destination factors used by blending formula
// R = S*sf + D*df
type BlendFactors struct {
	SrcFactor, DstFactor BlendFactor
}

// SourceBlendFactors is default BlendFactors used by AcceleratedCommand.
var SourceBlendFactors = BlendFactors{SrcFactor: One, DstFactor: Zero}

// SetBlendFactors sets source and dest factors for blending formula:
// R = S*sf + D*df
func (r *Renderer) SetBlendFactors(factors BlendFactors) {
	r.blendFactors = factors
}

// DrawArrays draws primitives (such as triangles) using vertices defined in VertexArray.
//
// Before primitive is drawn this method validates if
func (r *Renderer) DrawArrays(array *VertexArray, mode Mode, first, count int) {
	r.validateAttributeTypes(array)
	r.api.BindVertexArray(array.id)
	r.api.BlendFunc(uint32(r.blendFactors.SrcFactor), uint32(r.blendFactors.DstFactor))
	r.api.DrawArrays(mode.glMode, int32(first), int32(count))
}

func (r *Renderer) validateAttributeTypes(array *VertexArray) {
	if len(array.layout) > len(r.program.attributes) {
		msg := fmt.Sprintf("vertex array has more enabled attributes (%d) than program (%d)", len(array.layout), len(r.program.attributes))
		panic(msg)
	}
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
	r.api.ClearColor(color.RGBAf())
	r.api.Clear(colorBufferBit)
}

// AcceleratedCommand is an image.AcceleratedCommand implementation. It delegates
// the drawing to Command.
type AcceleratedCommand struct {
	command   Command
	program   *Program
	api       API
	allImages allImages
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
		panic("output image created in a different OpenGL context than program or deleted")
	}
	loc := output.Location
	if shouldSkipCommandProcessing(loc, img) {
		return
	}
	x := int32(loc.X)
	y := int32(loc.Y)
	w := int32(loc.Width)
	h := int32(loc.Height)
	if x+w > int32(img.width) {
		w = int32(img.width) - x
	}
	if y+h > int32(img.height) {
		h = int32(img.height) - y
	}
	y = int32(img.height) - h - y

	c.program.use()
	c.api.Enable(scissorTest)
	c.api.Enable(blend)
	c.api.BindFramebuffer(framebuffer, img.frameBufferID)
	c.api.Scissor(x, y, w, h)
	c.api.Viewport(x, y, w, h)

	renderer := &Renderer{
		program:      c.program,
		api:          c.api,
		allImages:    c.allImages,
		blendFactors: SourceBlendFactors,
	}

	c.command.RunGL(renderer, selections)
}

func shouldSkipCommandProcessing(output image.AcceleratedImageLocation, img *AcceleratedImage) bool {
	if output.X >= img.width {
		return true
	}
	if output.Y >= img.height {
		return true
	}
	if output.Width == 0 {
		return true
	}
	if output.Height == 0 {
		return true
	}
	if output.X+output.Width <= 0 {
		// this not only skips unneeded processing but also fixes bug with Intel Iris GPU on MAC - see #115
		return true
	}
	if output.Y+output.Height <= 0 {
		// this not only skips unneeded processing but also fixes bug with Intel Iris GPU on MAC - see #115
		return true
	}
	return false
}
