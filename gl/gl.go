// Package gl provides Go abstractions for interacting with OpenGL in a safer
// and easier way.
//
// It may be used with following versions and subsets of OpenGL:
// 	* OpenGL 3.3 and never
// 	* OpenGL ES 3.0 and never
package gl

import (
	"unsafe"

	"github.com/jacekolszak/pixiq/image"
)

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
	b.api.BufferSubData(arrayBuffer, offset*4, len(data)*4, b.api.Ptr(data))
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
	b.api.GetBufferSubData(arrayBuffer, offset*4, size*4, b.api.Ptr(output))
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

// ID returns VertexArray identifier (aka name)
func (a *VertexArray) ID() uint32 {
	return a.id
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}
