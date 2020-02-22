// Package gl provides Go abstractions for interacting with OpenGL in a safer way.
//
// It may be used with following versions and subsets of OpenGL:
// 	* OpenGL 3.3 and never
// 	* OpenGL ES 2.0 and never (TODO or 3.0?)
package gl

import (
	"fmt"
	"reflect"
	"unsafe"
)

// API is a gateway for directly accessing OpenGL driver.
type API interface {
	// GenBuffers generates buffer object names
	GenBuffers(n int32, buffers *uint32)
	// BindBuffer binds a named buffer object
	BindBuffer(target uint32, buffer uint32)
	// BufferData creates and initializes a buffer object's data     store
	BufferData(target uint32, size int, data unsafe.Pointer, usage uint32)
	// BufferSubData updates a subset of a buffer object's data store
	BufferSubData(target uint32, offset int, size int, data unsafe.Pointer)
	// GetBufferSubData returns a subset of a buffer object's data store
	GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer)
	// DeleteBuffers deletes named buffer objects
	DeleteBuffers(n int32, buffers *uint32)
}

const (
	arrayBuffer = 0x8892
	staticDraw  = 0x88E4
)

// ContextOf returns an OpenGL's Context for given API.
func ContextOf(api API) *Context {
	if api == nil {
		panic("nil api")
	}
	return &Context{
		api:             api,
		vertexBufferIDs: vertexBufferIDs{},
	}
}

// Context is an OpenGL context
type Context struct {
	api             API
	vertexBufferIDs vertexBufferIDs
}

// VertexBuffer contains data about vertices.
type VertexBuffer interface {
	// ID returns OpenGL identifier/name.
	ID() uint32
}

// vertexBufferIDs contains all vertex buffer identifiers in OpenGL context
type vertexBufferIDs map[VertexBuffer]uint32

// NewFloatVertexBuffer creates an OpenGL's Vertex Buffer Object (VBO) containing only float32 numbers.
func (g *Context) NewFloatVertexBuffer(size int) *FloatVertexBuffer {
	if size < 0 {
		panic("negative size")
	}
	var id uint32
	g.api.GenBuffers(1, &id)
	g.api.BindBuffer(arrayBuffer, id)
	g.api.BufferData(arrayBuffer, size*4, Ptr(nil), staticDraw) // FIXME: Parametrize usage
	vb := &FloatVertexBuffer{
		id:   id,
		size: size,
		api:  g.api,
	}
	g.vertexBufferIDs[vb] = id
	return vb
}

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

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
//
// For example:
//
// 	var data []uint8
// 	...
// 	gl.TexImage2D(gl.TEXTURE_2D, ..., gl.UNSIGNED_BYTE, gl.Ptr(&data[0]))
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
