package opengl

import (
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type context struct {
	runInOpenGLThread func(func()) // FIXME: to be removed
	glThread          *glThread
	mutex             sync.Mutex

	genBuffers       *genBuffers
	clear            *clear
	bindTexture      *bindTexture
	texSubImage2D    *texSubImage2D
	bindBuffer       *bindBuffer
	drawArrays       *drawArrays
	bindVertexArray  *bindVertexArray
	uniform1f        *uniform1f
	uniform2f        *uniform2f
	uniform3f        *uniform3f
	uniform4f        *uniform4f
	uniform1i        *uniform1i
	uniform2i        *uniform2i
	uniform3i        *uniform3i
	uniform4i        *uniform4i
	useProgram       *useProgram
	bindFramebuffer  *bindFramebuffer
	scissor          *scissor
	viewport         *viewport
	uniformMatrix3fv *uniformMatrix3fv
	uniformMatrix4fv *uniformMatrix4fv
	activeTexture    *activeTexture
}

func newContext(runInOpenGLThread func(func()), glThread *glThread) *context {
	return &context{
		runInOpenGLThread: runInOpenGLThread,
		glThread:          glThread,
		genBuffers:        &genBuffers{},
		clear:             &clear{},
		bindTexture:       &bindTexture{},
		texSubImage2D:     &texSubImage2D{},
		bindBuffer:        &bindBuffer{},
		drawArrays:        &drawArrays{},
		bindVertexArray:   &bindVertexArray{},
		uniform1f:         &uniform1f{},
		uniform2f:         &uniform2f{},
		uniform3f:         &uniform3f{},
		uniform4f:         &uniform4f{},
		uniform1i:         &uniform1i{},
		uniform2i:         &uniform2i{},
		uniform3i:         &uniform3i{},
		uniform4i:         &uniform4i{},
		useProgram:        &useProgram{},
		bindFramebuffer:   &bindFramebuffer{},
		scissor:           &scissor{},
		viewport:          &viewport{},
		uniformMatrix3fv:  &uniformMatrix3fv{},
		uniformMatrix4fv:  &uniformMatrix4fv{},
		activeTexture:     &activeTexture{},
	}
}

type genBuffers struct {
	n       int32
	buffers *uint32
}

func (g *genBuffers) run() {
	gl.GenBuffers(g.n, g.buffers)
}

// GenBuffers generates buffer object names
func (g *context) GenBuffers(n int32, buffers *uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.genBuffers.n = n
	g.genBuffers.buffers = buffers
	g.glThread.execute(g.genBuffers)
}

type bindBuffer struct {
	target uint32
	buffer uint32
}

func (b *bindBuffer) run() {
	gl.BindBuffer(b.target, b.buffer)
}

// BindBuffer binds a named buffer object
func (g *context) BindBuffer(target uint32, buffer uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.bindBuffer.target = target
	g.bindBuffer.buffer = buffer
	g.glThread.executeAsync(g.bindBuffer)
}

// BufferData creates and initializes a buffer object's data store
func (g *context) BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.BufferData(target, size, data, usage)
	})
}

// BufferSubData updates a subset of a buffer object's data store
func (g *context) BufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.BufferSubData(target, offset, size, data)
	})
}

// GetBufferSubData returns a subset of a buffer object's data store
func (g *context) GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetBufferSubData(target, offset, size, data)
	})
}

// DeleteBuffers deletes named buffer objects
func (g *context) DeleteBuffers(n int32, buffers *uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.DeleteBuffers(n, buffers)
	})
}

// GenVertexArrays generates vertex array object names
func (g *context) GenVertexArrays(n int32, arrays *uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GenVertexArrays(n, arrays)
	})
}

// DeleteVertexArrays deletes vertex array objects
func (g *context) DeleteVertexArrays(n int32, arrays *uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.DeleteVertexArrays(n, arrays)
	})
}

type bindVertexArray struct {
	array uint32
}

func (b *bindVertexArray) run() {
	gl.BindVertexArray(b.array)
}

// BindVertexArray binds a vertex array object
func (g *context) BindVertexArray(array uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.bindVertexArray.array = array
	g.glThread.executeAsync(g.bindVertexArray)
}

// VertexAttribPointer defines an array of generic vertex attribute data
func (g *context) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.VertexAttribPointer(index, size, xtype, normalized, stride, pointer)
	})
}

// EnableVertexAttribArray enables a generic vertex attribute array
func (g *context) EnableVertexAttribArray(index uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.EnableVertexAttribArray(index)
	})
}

// CreateShader creates a shader object
func (g *context) CreateShader(xtype uint32) uint32 {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	var id uint32
	g.runInOpenGLThread(func() {
		id = gl.CreateShader(xtype)
	})
	return id
}

// ShaderSource replaces the source code in a shader object
func (g *context) ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.ShaderSource(shader, count, xstring, length)
	})
}

// CompileShader compiles a shader object
func (g *context) CompileShader(shader uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.CompileShader(shader)
	})
}

// GetShaderiv returns a parameter from a shader object
func (g *context) GetShaderiv(shader uint32, pname uint32, params *int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetShaderiv(shader, pname, params)
	})
}

// GetShaderInfoLog returns the information log for a shader object
func (g *context) GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetShaderInfoLog(shader, bufSize, length, infoLog)
	})
}

// DeleteShader deletes a shader object
func (g *context) DeleteShader(shader uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.DeleteShader(shader)
	})
}

// AttachShader attaches a shader object to a program object
func (g *context) AttachShader(program uint32, shader uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.AttachShader(program, shader)
	})
}

// LinkProgram links a program object
func (g *context) LinkProgram(program uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.LinkProgram(program)
	})
}

// GetProgramiv returns a parameter from a program object
func (g *context) GetProgramiv(program uint32, pname uint32, params *int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetProgramiv(program, pname, params)
	})
}

// GetProgramInfoLog returns the information log for a program object
func (g *context) GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetProgramInfoLog(program, bufSize, length, infoLog)
	})
}

type useProgram struct {
	program uint32
}

func (u *useProgram) run() {
	gl.UseProgram(u.program)
}

// UseProgram installs a program object as part of current rendering state
func (g *context) UseProgram(program uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.useProgram.program = program
	g.glThread.executeAsync(g.useProgram)
}

// CreateProgram creates a program object
func (g *context) CreateProgram() uint32 {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	var program uint32
	g.runInOpenGLThread(func() {
		program = gl.CreateProgram()
	})
	return program
}

// GetActiveUniform returns information about an active uniform variable for the specified program object
func (g *context) GetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetActiveUniform(program, index, bufSize, length, size, xtype, name)
	})
}

// GetActiveAttrib returns information about an active attribute variable for the specified program object
func (g *context) GetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetActiveAttrib(program, index, bufSize, length, size, xtype, name)
	})
}

// GetAttribLocation returns the location of an attribute variable
func (g *context) GetAttribLocation(program uint32, name *uint8) int32 {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	var loc int32
	g.runInOpenGLThread(func() {
		loc = gl.GetAttribLocation(program, name)
	})
	return loc
}

// Enable enables or disable server-side GL capabilities
func (g *context) Enable(cap uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.Enable(cap)
	})
}

type bindFramebuffer struct {
	target      uint32
	framebuffer uint32
}

func (b *bindFramebuffer) run() {
	gl.BindFramebuffer(b.target, b.framebuffer)
}

// BindFramebuffer binds a framebuffer to a framebuffer target
func (g *context) BindFramebuffer(target uint32, framebuffer uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.bindFramebuffer.target = target
	g.bindFramebuffer.framebuffer = framebuffer
	g.glThread.executeAsync(g.bindFramebuffer)
}

type scissor struct {
	x      int32
	y      int32
	width  int32
	height int32
}

func (s *scissor) run() {
	gl.Scissor(s.x, s.y, s.width, s.height)
}

// Scissor defines the scissor box
func (g *context) Scissor(x int32, y int32, width int32, height int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.scissor.x = x
	g.scissor.y = y
	g.scissor.width = width
	g.scissor.height = height
	g.glThread.executeAsync(g.scissor)
}

type viewport struct {
	x      int32
	y      int32
	width  int32
	height int32
}

func (v *viewport) run() {
	gl.Viewport(v.x, v.y, v.width, v.height)
}

// Viewport sets the viewport
func (g *context) Viewport(x int32, y int32, width int32, height int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.viewport.x = x
	g.viewport.y = y
	g.viewport.width = width
	g.viewport.height = height
	g.glThread.executeAsync(g.viewport)
}

// ClearColor specifies clear values for the color buffers
func (g *context) ClearColor(red float32, green float32, blue float32, alpha float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.ClearColor(red, green, blue, alpha)
	})
}

type clear struct {
	mask uint32
}

func (c *clear) run() {
	gl.Clear(c.mask)
}

// Clear clears buffers to preset values
func (g *context) Clear(mask uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.clear.mask = mask
	g.glThread.executeAsync(g.clear)
}

type drawArrays struct {
	mode  uint32
	first int32
	count int32
}

func (d *drawArrays) run() {
	gl.DrawArrays(d.mode, d.first, d.count)
}

// DrawArrays render primitives from array data
func (g *context) DrawArrays(mode uint32, first int32, count int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.drawArrays.mode = mode
	g.drawArrays.first = first
	g.drawArrays.count = count
	g.glThread.executeAsync(g.drawArrays)
}

type uniform1f struct {
	location int32
	v0       float32
}

func (u *uniform1f) run() {
	gl.Uniform1f(u.location, u.v0)
}

// Uniform1f specifies the value of a uniform variable for the current program object
func (g *context) Uniform1f(location int32, v0 float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform1f.location = location
	g.uniform1f.v0 = v0
	g.glThread.executeAsync(g.uniform1f)
}

type uniform2f struct {
	location int32
	v0       float32
	v1       float32
}

func (u *uniform2f) run() {
	gl.Uniform2f(u.location, u.v0, u.v1)
}

// Uniform2f specifies the value of a uniform variable for the current program object
func (g *context) Uniform2f(location int32, v0 float32, v1 float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform2f.location = location
	g.uniform2f.v0 = v0
	g.uniform2f.v1 = v1
	g.glThread.executeAsync(g.uniform2f)
}

type uniform3f struct {
	location int32
	v0       float32
	v1       float32
	v2       float32
}

func (u *uniform3f) run() {
	gl.Uniform3f(u.location, u.v0, u.v1, u.v2)
}

// Uniform3f specifies the value of a uniform variable for the current program object
func (g *context) Uniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform3f.location = location
	g.uniform3f.v0 = v0
	g.uniform3f.v1 = v1
	g.uniform3f.v2 = v2
	g.glThread.executeAsync(g.uniform3f)
}

type uniform4f struct {
	location int32
	v0       float32
	v1       float32
	v2       float32
	v3       float32
}

func (u *uniform4f) run() {
	gl.Uniform4f(u.location, u.v0, u.v1, u.v2, u.v3)
}

// Uniform4f specifies the value of a uniform variable for the current program object
func (g *context) Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform4f.location = location
	g.uniform4f.v0 = v0
	g.uniform4f.v1 = v1
	g.uniform4f.v2 = v2
	g.uniform4f.v3 = v3
	g.glThread.executeAsync(g.uniform4f)
}

type uniform1i struct {
	location int32
	v0       int32
}

func (u *uniform1i) run() {
	gl.Uniform1i(u.location, u.v0)
}

// Uniform1i specifies the value of a uniform variable for the current program object
func (g *context) Uniform1i(location int32, v0 int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform1i.location = location
	g.uniform1i.v0 = v0
	g.glThread.executeAsync(g.uniform1i)
}

type uniform2i struct {
	location int32
	v0       int32
	v1       int32
}

func (u *uniform2i) run() {
	gl.Uniform2i(u.location, u.v0, u.v1)
}

// Uniform2i specifies the value of a uniform variable for the current program object
func (g *context) Uniform2i(location int32, v0 int32, v1 int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform2i.location = location
	g.uniform2i.v0 = v0
	g.uniform2i.v1 = v1
	g.glThread.executeAsync(g.uniform2i)
}

type uniform3i struct {
	location int32
	v0       int32
	v1       int32
	v2       int32
}

func (u *uniform3i) run() {
	gl.Uniform3i(u.location, u.v0, u.v1, u.v2)
}

// Uniform3i specifies the value of a uniform variable for the current program object
func (g *context) Uniform3i(location int32, v0 int32, v1 int32, v2 int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform3i.location = location
	g.uniform3i.v0 = v0
	g.uniform3i.v1 = v1
	g.uniform3i.v2 = v2
	g.glThread.executeAsync(g.uniform3i)
}

type uniform4i struct {
	location int32
	v0       int32
	v1       int32
	v2       int32
	v3       int32
}

func (u *uniform4i) run() {
	gl.Uniform4i(u.location, u.v0, u.v1, u.v2, u.v3)
}

// Uniform4i specifies the value of a uniform variable for the current program object
func (g *context) Uniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniform4i.location = location
	g.uniform4i.v0 = v0
	g.uniform4i.v1 = v1
	g.uniform4i.v2 = v2
	g.uniform4i.v3 = v3
	g.glThread.executeAsync(g.uniform4i)
}

type uniformMatrix3fv struct {
	location  int32
	count     int32
	transpose bool
	value     *float32
}

func (u *uniformMatrix3fv) run() {
	gl.UniformMatrix3fv(u.location, u.count, u.transpose, u.value)
}

// UniformMatrix3fv specifies the value of a uniform variable for the current program object
func (g *context) UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniformMatrix3fv.location = location
	g.uniformMatrix3fv.count = count
	g.uniformMatrix3fv.transpose = transpose
	g.uniformMatrix3fv.value = value
	g.glThread.execute(g.uniformMatrix3fv) // fixme: async maybe?
}

type uniformMatrix4fv struct {
	location  int32
	count     int32
	transpose bool
	value     *float32
}

func (u *uniformMatrix4fv) run() {
	gl.UniformMatrix4fv(u.location, u.count, u.transpose, u.value)
}

// UniformMatrix4fv specifies the value of a uniform variable for the current program object
func (g *context) UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.uniformMatrix4fv.location = location
	g.uniformMatrix4fv.count = count
	g.uniformMatrix4fv.transpose = transpose
	g.uniformMatrix4fv.value = value
	g.glThread.execute(g.uniformMatrix4fv) // fixme: async maybe?
}

type activeTexture struct {
	texture uint32
}

func (a *activeTexture) run() {
	gl.ActiveTexture(a.texture)
}

// ActiveTexture selects active texture unit
func (g *context) ActiveTexture(texture uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.activeTexture.texture = texture
	g.glThread.executeAsync(g.activeTexture)
}

type bindTexture struct {
	target  uint32
	texture uint32
}

func (b *bindTexture) run() {
	gl.BindTexture(b.target, b.texture)
}

// BindTexture binds a named texture to a texturing target
func (g *context) BindTexture(target uint32, texture uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.bindTexture.target = target
	g.bindTexture.texture = texture
	g.glThread.executeAsync(g.bindTexture)
}

// GetIntegerv returns the value or values of the specified parameter
func (g *context) GetIntegerv(pname uint32, data *int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetIntegerv(pname, data)
	})
}

// GenTextures generates texture names
func (g *context) GenTextures(n int32, textures *uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GenTextures(n, textures)
	})
}

// TexImage2D specifies a two-dimensional texture image
func (g *context) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.TexImage2D(target, level, internalformat, width, height, border, format, xtype, pixels)
	})
}

// TexParameteri sets texture parameter
func (g *context) TexParameteri(target uint32, pname uint32, param int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.TexParameteri(target, pname, param)
	})
}

// GenFramebuffers generates framebuffer object names
func (g *context) GenFramebuffers(n int32, framebuffers *uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GenFramebuffers(n, framebuffers)
	})
}

// FramebufferTexture2D attaches a level of a texture object as a logical buffer to the currently bound framebuffer object
func (g *context) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
	})
}

type texSubImage2D struct {
	target  uint32
	level   int32
	xoffset int32
	yoffset int32
	width   int32
	height  int32
	format  uint32
	xtype   uint32
	pixels  unsafe.Pointer
}

func (t *texSubImage2D) run() {
	gl.TexSubImage2D(t.target, t.level, t.xoffset, t.yoffset, t.width, t.height, t.format, t.xtype, t.pixels)
}

// TexSubImage2D specifies a two-dimensional texture subimage
func (g *context) TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.texSubImage2D.target = target
	g.texSubImage2D.level = level
	g.texSubImage2D.xoffset = xoffset
	g.texSubImage2D.yoffset = yoffset
	g.texSubImage2D.width = width
	g.texSubImage2D.height = height
	g.texSubImage2D.format = format
	g.texSubImage2D.xtype = xtype
	g.texSubImage2D.pixels = pixels
	g.glThread.execute(g.texSubImage2D)
}

// GetTexImage returns a texture image
func (g *context) GetTexImage(target uint32, level int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.GetTexImage(target, level, format, xtype, pixels)
	})
}

// GetError returns error information
func (g *context) GetError() uint32 {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	var code uint32
	g.runInOpenGLThread(func() {
		code = gl.GetError()
	})
	return code
}

// ReadPixels reads a block of pixels from the frame buffer
func (g *context) ReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.runInOpenGLThread(func() {
		gl.ReadPixels(x, y, width, height, format, xtype, pixels)
	})
}

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
//
// For example:
//
// 	var data []uint8
// 	...
// 	api.TexImage2D(..., api.Ptr(&data[0]))
func (g *context) Ptr(data interface{}) unsafe.Pointer {
	return gl.Ptr(data)
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func (g *context) PtrOffset(offset int) unsafe.Pointer {
	return gl.PtrOffset(offset)
}

// GoStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func (g *context) GoStr(cstr *uint8) string {
	return gl.GoStr(cstr)
}

// Strs takes a list of Go strings (with or without null-termination) and
// returns their C counterpart.
//
// The returned free function must be called once you are done using the strings
// in order to free the memory.
//
// If no strings are provided as a parameter this function will panic.
func (g *context) Strs(strs ...string) (cstrs **uint8, free func()) {
	return gl.Strs(strs...)
}
