// Package gl provides Go abstractions for interacting with OpenGL in a safer way.
//
// It may be used with following versions and subsets of OpenGL:
// 	* OpenGL 3.3 and never
// 	* OpenGL ES 3.0 and never
package gl

import "C"
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
	// Returns the location of an attribute variable
	GetAttribLocation(program uint32, name *uint8) int32

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

// Camel-cased GL constants
const (
	arrayBuffer              = 0x8892
	staticDraw               = 0x88E4
	float                    = 0x1406
	floatVec2                = 0x8B50
	floatVec3                = 0x8B51
	floatVec4                = 0x8B52
	vertexShader             = 0x8B31
	fragmentShader           = 0x8B30
	compileStatus            = 0x8B81
	ffalse                   = 0
	infoLogLength            = 0x8B84
	linkStatus               = 0x8B82
	activeUniforms           = 0x8B86
	activeUniformMaxLength   = 0x8B87
	activeAttributeMaxLength = 0x8B8A
	activeAttributes         = 0x8B89
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
func (c *Context) NewFloatVertexBuffer(size int) *FloatVertexBuffer {
	if size < 0 {
		panic("negative size")
	}
	var id uint32
	c.api.GenBuffers(1, &id)
	c.api.BindBuffer(arrayBuffer, id)
	c.api.BufferData(arrayBuffer, size*4, Ptr(nil), staticDraw) // FIXME: Parametrize usage
	vb := &FloatVertexBuffer{
		id:   id,
		size: size,
		api:  c.api,
	}
	c.vertexBufferIDs[vb] = id
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

func (c *Context) compileVertexShader(src string) (uint32, error) {
	return c.compileShader(vertexShader, src)
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

func (p *program) use() {
	p.api.UseProgram(p.id)
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

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}
