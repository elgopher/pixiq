package gl

import (
	"fmt"

	"github.com/jacekolszak/pixiq/image"
)

// Context is an OpenGL context
type Context struct {
	api             API
	vertexBufferIDs vertexBufferIDs
	allImages       allImages
	capabilities    *Capabilities
}

// API returns API passed during Context construction. It may be used for directly
// accessing OpenGL.
func (c *Context) API() API {
	return c.api
}

// Capabilities returns parameter values reported by current OpenGL instance.
func (c *Context) Capabilities() *Capabilities {
	return c.capabilities
}

// Capabilities contains parameter values reported by current OpenGL instance.
type Capabilities struct {
	maxTextureSize int
}

// MaxTextureSize returns OpenGL's MAX_TEXTURE_SIZE
func (c Capabilities) MaxTextureSize() int {
	return c.maxTextureSize
}

type glError uint32

func (e glError) Error() string {
	return fmt.Sprintf("gl error: %d", uint32(e))
}

// IsOutOfMemory returns true if given error indicates that OpenGL driver reported
// out-of-memory.
//
// This error is not recoverable. Once you get it - you have to destroy the whole
// OpenGL context and start a new one.
func IsOutOfMemory(err error) bool {
	e, ok := err.(glError)
	if !ok {
		return false
	}
	return e == outOfMemory
}

// Error returns next error reported by OpenGL driver. For performance reasons should
// be used sporadically, at most once per frame.
//
// See http://docs.gl/gl3/glGetError
func (c *Context) Error() error {
	var code = c.api.GetError()
	if code == noError {
		return nil
	}
	return glError(code)
}

type Usage struct {
	glUsage uint32
}

var (
	// The data store contents will be modified once and used at most a few times.
	// The data store contents are modified by the application, and used as the source
	// for GL drawing and image specification commands.
	StreamDraw = Usage{glUsage: streamDraw}
	// The data store contents will be modified once and used many times.
	// The data store contents are modified by the application, and used as the source
	// for GL drawing and image specification commands.
	StaticDraw = Usage{glUsage: staticDraw}
	// The data store contents will be modified repeatedly and used many times.
	// The data store contents are modified by the application, and used as the source
	// for GL drawing and image specification commands.
	DynamicDraw = Usage{glUsage: dynamicDraw}
)

type UsageNature int

// NewFloatVertexBuffer creates an OpenGL's Vertex Buffer Object (VBO) containing only float32 numb)ers.
func (c *Context) NewFloatVertexBuffer(size int, usage Usage) *FloatVertexBuffer {
	if size < 0 {
		panic("negative size")
	}
	var id uint32
	c.api.GenBuffers(1, &id)
	c.api.BindBuffer(arrayBuffer, id)
	c.api.BufferData(arrayBuffer, size*4, c.api.Ptr(nil), usage.glUsage)
	vb := &FloatVertexBuffer{
		id:   id,
		size: size,
		api:  c.api,
	}
	c.vertexBufferIDs[vb] = id
	return vb
}

// NewVertexArray creates a new instance of VertexArray. All vertex attributes
// specified in layout will be enabled.
func (c *Context) NewVertexArray(layout VertexLayout) *VertexArray {
	if len(layout) == 0 {
		panic("empty layout")
	}
	var id uint32
	c.api.GenVertexArrays(1, &id)
	c.api.BindVertexArray(id)
	for i := 0; i < len(layout); i++ {
		c.api.EnableVertexAttribArray(uint32(i))
	}
	return &VertexArray{
		id:              id,
		layout:          layout,
		api:             c.api,
		vertexBufferIDs: c.vertexBufferIDs,
	}
}

// CompileFragmentShader compiles fragment shader source code written in GLSL.
func (c *Context) CompileFragmentShader(sourceCode string) (*FragmentShader, error) {
	shaderID, err := c.compileShader(fragmentShader, sourceCode)
	if err != nil {
		return nil, err
	}
	return &FragmentShader{id: shaderID}, nil
}

// FragmentShader is a part of an OpenGL program which transforms each fragment
// (pixel) color into another one
type FragmentShader struct {
	id uint32
}

// CompileVertexShader compiles vertex shader source code written in GLSL.
func (c *Context) CompileVertexShader(sourceCode string) (*VertexShader, error) {
	shaderID, err := c.compileShader(vertexShader, sourceCode)
	if err != nil {
		return nil, err
	}
	return &VertexShader{id: shaderID}, nil
}

// VertexShader is a part of an OpenGL program which applies transformations
// to drawn vertices.
type VertexShader struct {
	id uint32
}

func (c *Context) compileShader(xtype uint32, src string) (uint32, error) {
	if src == "" {
		src = " "
	}
	shaderID := c.api.CreateShader(xtype)
	srcXString, free := c.api.Strs(src)
	defer free()
	length := int32(len(src))
	c.api.ShaderSource(shaderID, 1, srcXString, &length)
	c.api.CompileShader(shaderID)
	var success int32
	c.api.GetShaderiv(shaderID, compileStatus, &success)
	if success == ffalse {
		var logLen int32
		c.api.GetShaderiv(shaderID, infoLogLength, &logLen)
		infoLog := make([]byte, logLen)
		if logLen > 0 {
			c.api.GetShaderInfoLog(shaderID, logLen, nil, &infoLog[0])
		}
		return 0, fmt.Errorf("glCompileShader failed: %s", string(infoLog))
	}
	return shaderID, nil
}

// LinkProgram links an OpenGL program from shaders. Created program can be used
// in image.Modify
func (c *Context) LinkProgram(vertexShader *VertexShader, fragmentShader *FragmentShader) (*Program, error) {
	if vertexShader == nil {
		panic("nil vertexShader")
	}
	if fragmentShader == nil {
		panic("nil fragmentShader")
	}
	var (
		program          *program
		err              error
		uniformLocations map[string]int32
		attributes       map[int32]attribute
	)
	program, err = c.linkProgram(vertexShader.id, fragmentShader.id)
	if err == nil {
		uniformLocations = program.activeUniformLocations()
		attributes = program.attributes()
	}
	if err != nil {
		return nil, err
	}
	return &Program{
		program:          program,
		api:              c.api,
		uniformLocations: uniformLocations,
		attributes:       attributes,
		allImages:        c.allImages,
	}, err
}

func (c *Context) linkProgram(shaderIDs ...uint32) (*program, error) {
	programID := c.api.CreateProgram()
	for _, shaderID := range shaderIDs {
		c.api.AttachShader(programID, shaderID)
	}
	c.api.LinkProgram(programID)
	var success int32
	c.api.GetProgramiv(programID, linkStatus, &success)
	if success == ffalse {
		var infoLogLen int32
		c.api.GetProgramiv(programID, infoLogLength, &infoLogLen)
		infoLog := make([]byte, infoLogLen)
		if infoLogLen > 0 {
			c.api.GetProgramInfoLog(programID, infoLogLen, nil, &infoLog[0])
		}
		return nil, fmt.Errorf("error linking program: %s", string(infoLog))
	}
	return &program{
		id:  programID,
		api: c.api,
	}, nil
}

type program struct {
	api API
	id  uint32
}

func (p *program) activeUniformLocations() map[string]int32 {
	locationsByName := map[string]int32{}
	var count, bufSize, length, nameMaxLength int32
	var xtype uint32
	p.api.GetProgramiv(p.id, activeUniformMaxLength, &nameMaxLength)
	name := make([]byte, nameMaxLength)
	p.api.GetProgramiv(p.id, activeUniforms, &count)
	for location := int32(0); location < count; location++ {
		p.api.GetActiveUniform(p.id, uint32(location), nameMaxLength, &bufSize, &length, &xtype, &name[0])
		goName := p.api.GoStr(&name[0])
		locationsByName[goName] = location
	}
	return locationsByName
}

type attribute struct {
	typ  Type
	name string
}

func (p *program) attributes() map[int32]attribute {
	var count, bufSize, length, nameMaxLength int32
	var xtype uint32
	p.api.GetProgramiv(p.id, activeAttributeMaxLength, &nameMaxLength)
	name := make([]byte, nameMaxLength)
	p.api.GetProgramiv(p.id, activeAttributes, &count)
	attributes := map[int32]attribute{}
	for i := int32(0); i < count; i++ {
		p.api.GetActiveAttrib(p.id, uint32(i), nameMaxLength, &bufSize, &length, &xtype, &name[0])
		location := p.api.GetAttribLocation(p.id, &name[0])
		attributes[location] = attribute{typ: valueOf(xtype),
			name: p.api.GoStr(&name[0])}
	}
	return attributes
}

// Program is shaders linked together
type Program struct {
	*program
	uniformLocations map[string]int32
	attributes       map[int32]attribute
	api              API
	allImages        allImages
}

// AcceleratedCommand returns a potentially cached instance of *AcceleratedCommand.
func (p *Program) AcceleratedCommand(command Command) *AcceleratedCommand {
	return &AcceleratedCommand{
		command:   command,
		api:       p.api,
		program:   p,
		allImages: p.allImages,
	}
}

// ID returns program identifier (aka name)
func (p *Program) ID() uint32 {
	return p.id
}

func (p *Program) use() {
	if p.program != nil {
		p.api.UseProgram(p.id)
	}
}

// NewClearCommand returns a command clearing all pixels in image.Selection
func (c *Context) NewClearCommand() *ClearCommand {
	nilProgram := &Program{
		program:          nil,
		uniformLocations: map[string]int32{},
		attributes:       map[int32]attribute{},
		api:              c.api,
		allImages:        c.allImages,
	}
	cmd := &ClearCommand{}
	cmd.AcceleratedCommand = nilProgram.AcceleratedCommand(cmd)
	return cmd
}

// ClearCommand clears the image.Selection using given color. By default color
// is transparent.
type ClearCommand struct {
	*AcceleratedCommand
	color image.Color
}

// SetColor sets color which will be used to clear all pixels in image.Selection
func (c *ClearCommand) SetColor(color image.Color) {
	c.color = color
}

// RunGL implements gl.Command
func (c *ClearCommand) RunGL(renderer *Renderer, _ []image.AcceleratedImageSelection) {
	renderer.Clear(c.color)
}
