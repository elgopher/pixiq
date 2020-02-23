// Package gl provides Go abstractions for interacting with OpenGL in a safer
// and easier way.
//
// It may be used with following versions and subsets of OpenGL:
// 	* OpenGL 3.3 and never
// 	* OpenGL ES 3.0 and never
package gl

import (
	"fmt"
	"github.com/jacekolszak/pixiq/image"
	"reflect"
	"unsafe"
)

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

// NewContext returns an OpenGL's Context for given API.
func NewContext(api API) *Context {
	if api == nil {
		panic("nil api")
	}
	return &Context{
		api:             api,
		vertexBufferIDs: vertexBufferIDs{},
		allImages:       allImages{},
		capabilities:    gatherCapabilities(api),
	}
}

func gatherCapabilities(api API) *Capabilities {
	var maxTextureSizeVal int32
	api.GetIntegerv(maxTextureSize, &maxTextureSizeVal)
	return &Capabilities{
		maxTextureSize: int(maxTextureSizeVal),
	}
}

// VertexBuffer contains data about vertices.
type VertexBuffer interface {
	// ID returns OpenGL identifier/name.
	ID() uint32
}

// vertexBufferIDs contains all vertex buffer identifiers in OpenGL context
type vertexBufferIDs map[VertexBuffer]uint32
type allImages map[image.AcceleratedImage]*AcceleratedImage

// FloatVertexBuffer is a struct representing OpenGL's Vertex Buffer Object (VBO) containing only float32 numbers.
type FloatVertexBuffer struct {
	id      uint32
	deleted bool
	size    int
	api     API
}

// Size is the number of float values defined during creation time.
func (b *FloatVertexBuffer) Size() int {
	return b.size
}

// ID returns OpenGL identifier/name.
func (b *FloatVertexBuffer) ID() uint32 {
	return b.id
}

// Upload sends data to the vertex buffer. All slice data will be inserted starting at a given offset position.
//
// Panics when vertex buffer is too small to hold the data or offset is negative.
func (b *FloatVertexBuffer) Upload(offset int, data []float32) {
	if offset < 0 {
		panic("negative offset")
	}
	if b.size < len(data)+offset {
		panic("FloatVertexBuffer is to small to store data")
	}
	b.api.BindBuffer(arrayBuffer, b.id)
	b.api.BufferSubData(arrayBuffer, offset*4, len(data)*4, Ptr(data))
}

// Delete should be called whenever you don't plan to use vertex buffer anymore. Vertex Buffer is external resource
// (like file for example) and must be deleted manually
func (b *FloatVertexBuffer) Delete() {
	b.api.DeleteBuffers(1, &b.id)
	b.deleted = true
}

// Download gets data starting at a given offset in VRAM and put them into slice.
// Whole output slice will be filled with data, unless output slice is bigger then
// the vertex buffer.
func (b *FloatVertexBuffer) Download(offset int, output []float32) {
	if b.deleted {
		panic("deleted buffer")
	}
	if offset < 0 {
		panic("negative offset")
	}
	if len(output) == 0 {
		return
	}
	size := len(output)
	if size+offset > b.size {
		size = b.size - offset
	}
	b.api.BindBuffer(arrayBuffer, b.id)
	b.api.GetBufferSubData(arrayBuffer, offset*4, size*4, Ptr(output))
}

// VertexLayout defines data types of VertexArray locations.
type VertexLayout []Type

// Type is a kind of OpenGL's attribute.
type Type struct {
	components int32
	xtype      uint32
	name       string
}

func valueOf(xtype uint32) Type {
	switch xtype {
	case float:
		return Float
	case floatVec2:
		return Vec2
	case floatVec3:
		return Vec3
	case floatVec4:
		return Vec4
	}
	panic("not supported type")
}

func (t Type) String() string {
	return t.name
}

var (
	// Float is single-precision floating point number.
	// Equivalent of Go's float32.
	Float = Type{components: 1, xtype: float, name: "Float"}
	// Vec2 is a vector of two single-precision floating point numbers.
	// Equivalent of Go's [2]float32.
	Vec2 = Type{components: 2, xtype: float, name: "Vec2"}
	// Vec3 is a vector of three single-precision floating point numbers.
	// Equivalent of Go's [3]float32.
	Vec3 = Type{components: 3, xtype: float, name: "Vec3"}
	// Vec4 is a vector of four single-precision floating point numbers.
	// Equivalent of Go's [4]float32.
	Vec4 = Type{components: 4, xtype: float, name: "Vec4"}
)

// VertexArray is a thin abstraction for OpenGL's Vertex Array Object.
//
// https://www.khronos.org/opengl/wiki/Vertex_Specification#Vertex_Array_Object
type VertexArray struct {
	id              uint32
	layout          VertexLayout
	vertexBufferIDs vertexBufferIDs
	api             API
}

// Delete should be called whenever you don't plan to use VertexArray anymore.
// VertexArray is an external resource (like file for example) and must be deleted manually.
func (a *VertexArray) Delete() {
	a.api.DeleteVertexArrays(1, &a.id)
}

// VertexBufferPointer is a slice of VertexBuffer
type VertexBufferPointer struct {
	Buffer VertexBuffer
	Offset int
	Stride int
}

// Set sets a location of VertexArray pointing to VertexBuffer slice.
func (a *VertexArray) Set(location int, pointer VertexBufferPointer) {
	if pointer.Offset < 0 {
		panic("negative pointer offset")
	}
	if pointer.Stride < 0 {
		panic("negative pointer stride")
	}
	if pointer.Buffer == nil {
		panic("nil pointer buffer")
	}
	if location < 0 {
		panic("negative location")
	}
	if location >= len(a.layout) {
		panic("location out-of-bounds")
	}
	bufferID, ok := a.vertexBufferIDs[pointer.Buffer]
	if !ok {
		panic("vertex buffer has not been created in this context")
	}
	a.api.BindVertexArray(a.id)
	a.api.BindBuffer(arrayBuffer, bufferID)
	typ := a.layout[location]
	components := typ.components
	a.api.VertexAttribPointer(
		uint32(location),
		components,
		typ.xtype,
		false,
		int32(pointer.Stride*4),
		PtrOffset(pointer.Offset*4),
	)
}

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
//
// For example:
//
// 	var data []uint8
// 	...
// 	api.TexImage2D(texture2D, ..., unsignedByte, gl.Ptr(&data[0]))
func Ptr(data interface{}) unsafe.Pointer {
	if data == nil {
		return unsafe.Pointer(nil)
	}
	var addr unsafe.Pointer
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr:
		e := v.Elem()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			addr = unsafe.Pointer(e.UnsafeAddr())
		default:
			panic(fmt.Errorf("unsupported pointer to type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", e.Kind()))
		}
	case reflect.Uintptr:
		addr = unsafe.Pointer(v.Pointer())
	case reflect.Slice:
		addr = unsafe.Pointer(v.Index(0).UnsafeAddr())
	default:
		panic(fmt.Errorf("unsupported type %s; must be a slice or pointer to a singular scalar value or the first element of an array or slice", v.Type()))
	}
	return addr
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}
