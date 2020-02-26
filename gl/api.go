package gl

import "unsafe"

// API is a gateway for directly accessing OpenGL driver.
type API interface {
	// GenBuffers generates buffer object names
	GenBuffers(n int32, buffers *uint32)
	// BindBuffer binds a named buffer object
	BindBuffer(target uint32, buffer uint32)
	// BufferData creates and initializes a buffer object's data store
	BufferData(target uint32, size int, data unsafe.Pointer, usage uint32)
	// BufferSubData updates a subset of a buffer object's data store
	BufferSubData(target uint32, offset int, size int, data unsafe.Pointer)
	// GetBufferSubData returns a subset of a buffer object's data store
	GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer)
	// DeleteBuffers deletes named buffer objects
	DeleteBuffers(n int32, buffers *uint32)
	// GenVertexArrays generates vertex array object names
	GenVertexArrays(n int32, arrays *uint32)
	// DeleteVertexArrays deletes vertex array objects
	DeleteVertexArrays(n int32, arrays *uint32)
	// BindVertexArray binds a vertex array object
	BindVertexArray(array uint32)
	// VertexAttribPointer defines an array of generic vertex attribute data
	VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer)
	// EnableVertexAttribArray enables a generic vertex attribute array
	EnableVertexAttribArray(index uint32)
	// CreateShader creates a shader object
	CreateShader(xtype uint32) uint32
	// ShaderSource replaces the source code in a shader object
	ShaderSource(shader uint32, count int32, xstring **uint8, length *int32)
	// CompileShader compiles a shader object
	CompileShader(shader uint32)
	// GetShaderiv returns a parameter from a shader object
	GetShaderiv(shader uint32, pname uint32, params *int32)
	// GetShaderInfoLog returns the information log for a shader object
	GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8)
	// DeleteShader deletes a shader object
	DeleteShader(shader uint32)
	// AttachShader attaches a shader object to a program object
	AttachShader(program uint32, shader uint32)
	// LinkProgram links a program object
	LinkProgram(program uint32)
	// GetProgramiv returns a parameter from a program object
	GetProgramiv(program uint32, pname uint32, params *int32)
	// GetProgramInfoLog returns the information log for a program object
	GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8)
	// UseProgram installs a program object as part of current rendering state
	UseProgram(program uint32)
	// CreateProgram creates a program object
	CreateProgram() uint32
	// GetActiveUniform returns information about an active uniform variable for the specified program object
	GetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8)
	// GetActiveAttrib returns information about an active attribute variable for the specified program object
	GetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8)
	// GetAttribLocation returns the location of an attribute variable
	GetAttribLocation(program uint32, name *uint8) int32
	// Enable enables or disable server-side GL capabilities
	Enable(cap uint32)
	// BindFramebuffer binds a framebuffer to a framebuffer target
	BindFramebuffer(target uint32, framebuffer uint32)
	// Scissor defines the scissor box
	Scissor(x int32, y int32, width int32, height int32)
	// Viewport sets the viewport
	Viewport(x int32, y int32, width int32, height int32)
	// ClearColor specifies clear values for the color buffers
	ClearColor(red float32, green float32, blue float32, alpha float32)
	// Clear clears buffers to preset values
	Clear(mask uint32)
	// DrawArrays render primitives from array data
	DrawArrays(mode uint32, first int32, count int32)
	// Uniform1f specifies the value of a uniform variable for the current program object
	Uniform1f(location int32, v0 float32)
	// Uniform2f specifies the value of a uniform variable for the current program object
	Uniform2f(location int32, v0 float32, v1 float32)
	// Uniform3f specifies the value of a uniform variable for the current program object
	Uniform3f(location int32, v0 float32, v1 float32, v2 float32)
	// Uniform4f specifies the value of a uniform variable for the current program object
	Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32)
	// Uniform1i specifies the value of a uniform variable for the current program object
	Uniform1i(location int32, v0 int32)
	// Uniform2i specifies the value of a uniform variable for the current program object
	Uniform2i(location int32, v0 int32, v1 int32)
	// Uniform3i specifies the value of a uniform variable for the current program object
	Uniform3i(location int32, v0 int32, v1 int32, v2 int32)
	// Uniform4i specifies the value of a uniform variable for the current program object
	Uniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32)
	// UniformMatrix3fv specifies the value of a uniform variable for the current program object
	UniformMatrix3fv(location int32, count int32, transpose bool, value *float32)
	// UniformMatrix4fv specifies the value of a uniform variable for the current program object
	UniformMatrix4fv(location int32, count int32, transpose bool, value *float32)
	// ActiveTexture selects active texture unit
	ActiveTexture(texture uint32)
	// BindTexture binds a named texture to a texturing target
	BindTexture(target uint32, texture uint32)
	// GetIntegerv returns the value or values of the specified parameter
	GetIntegerv(pname uint32, data *int32)
	// GenTextures generates texture names
	GenTextures(n int32, textures *uint32)
	// TexImage2D specifies a two-dimensional texture image
	TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer)
	// TexParameteri sets texture parameter
	TexParameteri(target uint32, pname uint32, param int32)
	// GenFramebuffers generates framebuffer object names
	GenFramebuffers(n int32, framebuffers *uint32)
	// FramebufferTexture2D attaches a level of a texture object as a logical buffer to the currently bound framebuffer object
	FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32)
	// TexSubImage2D specifies a two-dimensional texture subimage
	TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer)
	// GetTexImage returns a texture image
	GetTexImage(target uint32, level int32, format uint32, xtype uint32, pixels unsafe.Pointer)
	// GetError returns error information
	GetError() uint32
	// ReadPixels reads a block of pixels from the frame buffer
	ReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer)
	// Ptr takes a slice or pointer (to a singular scalar value or the first
	// element of an array or slice) and returns its GL-compatible address.
	//
	// For example:
	//
	// 	var data []uint8
	// 	...
	// 	api.TexImage2D(..., api.Ptr(&data[0]))
	Ptr(data interface{}) unsafe.Pointer
	// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
	// Useful for functions such as glVertexAttribPointer that take pointer
	// parameters indicating an offset rather than an absolute memory address.
	PtrOffset(offset int) unsafe.Pointer
	// GoStr takes a null-terminated string returned by OpenGL and constructs a
	// corresponding Go string.
	GoStr(cstr *uint8) string
	// Strs takes a list of Go strings (with or without null-termination) and
	// returns their C counterpart.
	//
	// The returned free function must be called once you are done using the strings
	// in order to free the memory.
	//
	// If no strings are provided as a parameter this function will panic.
	Strs(strs ...string) (cstrs **uint8, free func())
}
