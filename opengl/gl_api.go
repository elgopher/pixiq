package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"unsafe"
)

// GenBuffers generates buffer object names
func (g *OpenGL) GenBuffers(n int32, buffers *uint32) {
	g.runInOpenGLThread(func() {
		gl.GenBuffers(n, buffers)
	})
}

// BindBuffer binds a named buffer object
func (g *OpenGL) BindBuffer(target uint32, buffer uint32) {
	g.runInOpenGLThread(func() {
		gl.BindBuffer(target, buffer)
	})
}

// BufferData creates and initializes a buffer object's data store
func (g *OpenGL) BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	g.runInOpenGLThread(func() {
		gl.BufferData(target, size, data, usage)
	})
}

// BufferSubData updates a subset of a buffer object's data store
func (g *OpenGL) BufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	g.runInOpenGLThread(func() {
		gl.BufferSubData(target, offset, size, data)
	})
}

// GetBufferSubData returns a subset of a buffer object's data store
func (g *OpenGL) GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	g.runInOpenGLThread(func() {
		gl.GetBufferSubData(target, offset, size, data)
	})
}

// DeleteBuffers deletes named buffer objects
func (g *OpenGL) DeleteBuffers(n int32, buffers *uint32) {
	g.runInOpenGLThread(func() {
		gl.DeleteBuffers(n, buffers)
	})
}

// GenVertexArrays generates vertex array object names
func (g *OpenGL) GenVertexArrays(n int32, arrays *uint32) {
	g.runInOpenGLThread(func() {
		gl.GenVertexArrays(n, arrays)
	})
}

// DeleteVertexArrays deletes vertex array objects
func (g *OpenGL) DeleteVertexArrays(n int32, arrays *uint32) {
	g.runInOpenGLThread(func() {
		gl.DeleteVertexArrays(n, arrays)
	})
}

// BindVertexArray binds a vertex array object
func (g *OpenGL) BindVertexArray(array uint32) {
	g.runInOpenGLThread(func() {
		gl.BindVertexArray(array)
	})
}

// VertexAttribPointer defines an array of generic vertex attribute data
func (g *OpenGL) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	g.runInOpenGLThread(func() {
		gl.VertexAttribPointer(index, size, xtype, normalized, stride, pointer)
	})
}

// EnableVertexAttribArray enables a generic vertex attribute array
func (g *OpenGL) EnableVertexAttribArray(index uint32) {
	g.runInOpenGLThread(func() {
		gl.EnableVertexAttribArray(index)
	})
}

// CreateShader creates a shader object
func (g *OpenGL) CreateShader(xtype uint32) uint32 {
	var id uint32
	g.runInOpenGLThread(func() {
		id = gl.CreateShader(xtype)
	})
	return id
}

// ShaderSource replaces the source code in a shader object
func (g *OpenGL) ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	g.runInOpenGLThread(func() {
		gl.ShaderSource(shader, count, xstring, length)
	})
}

// CompileShader compiles a shader object
func (g *OpenGL) CompileShader(shader uint32) {
	g.runInOpenGLThread(func() {
		gl.CompileShader(shader)
	})
}

// GetShaderiv returns a parameter from a shader object
func (g *OpenGL) GetShaderiv(shader uint32, pname uint32, params *int32) {
	g.runInOpenGLThread(func() {
		gl.GetShaderiv(shader, pname, params)
	})
}

// GetShaderInfoLog returns the information log for a shader object
func (g *OpenGL) GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	g.runInOpenGLThread(func() {
		gl.GetShaderInfoLog(shader, bufSize, length, infoLog)
	})
}

// DeleteShader deletes a shader object
func (g *OpenGL) DeleteShader(shader uint32) {
	g.runInOpenGLThread(func() {
		gl.DeleteShader(shader)
	})
}

// AttachShader attaches a shader object to a program object

func (g *OpenGL) AttachShader(program uint32, shader uint32) {
	g.runInOpenGLThread(func() {
		gl.AttachShader(program, shader)
	})
}

// LinkProgram links a program object
func (g *OpenGL) LinkProgram(program uint32) {
	g.runInOpenGLThread(func() {
		gl.LinkProgram(program)
	})
}

// GetProgramiv returns a parameter from a program object
func (g *OpenGL) GetProgramiv(program uint32, pname uint32, params *int32) {
	g.runInOpenGLThread(func() {
		gl.GetProgramiv(program, pname, params)
	})
}

// GetProgramInfoLog returns the information log for a program object
func (g *OpenGL) GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {
	g.runInOpenGLThread(func() {
		gl.GetProgramInfoLog(program, bufSize, length, infoLog)
	})
}

// UseProgram installs a program object as part of current rendering state
func (g *OpenGL) UseProgram(program uint32) {
	g.runInOpenGLThread(func() {
		gl.UseProgram(program)
	})
}

// CreateProgram creates a program object
func (g *OpenGL) CreateProgram() uint32 {
	var program uint32
	g.runInOpenGLThread(func() {
		program = gl.CreateProgram()
	})
	return program
}

// GetActiveUniform returns information about an active uniform variable for the specified program object
func (g *OpenGL) GetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	g.runInOpenGLThread(func() {
		gl.GetActiveUniform(program, index, bufSize, length, size, xtype, name)
	})
}

// GetActiveAttrib returns information about an active attribute variable for the specified program object
func (g *OpenGL) GetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	g.runInOpenGLThread(func() {
		gl.GetActiveAttrib(program, index, bufSize, length, size, xtype, name)
	})
}

// Returns the location of an attribute variable
func (g *OpenGL) GetAttribLocation(program uint32, name *uint8) int32 {
	var loc int32
	g.runInOpenGLThread(func() {
		loc = gl.GetAttribLocation(program, name)
	})
	return loc
}

// GoStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func (g *OpenGL) GoStr(cstr *uint8) string {
	return gl.GoStr(cstr)
}

// Strs takes a list of Go strings (with or without null-termination) and
// returns their C counterpart.
//
// The returned free function must be called once you are done using the strings
// in order to free the memory.
//
// If no strings are provided as a parameter this function will panic.
func (g *OpenGL) Strs(strs ...string) (cstrs **uint8, free func()) {
	return gl.Strs(strs...)
}
