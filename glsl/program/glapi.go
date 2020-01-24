package program

import "github.com/jacekolszak/pixiq/image"

type Buffers interface {
	NewFloatVertexBuffer(BufferUsage) FloatVertexBuffer
}

type BufferUsage int

const (
	StreamDraw BufferUsage = iota
	StreamRead
	StreamCopy
	StaticDraw
	StaticRead
	StaticCopy
	DynamicDraw
	DynamicRead
	DynamicCopy
)

// FloatVertexBuffer represents a buffer of floats held in a video card memory
type FloatVertexBuffer interface {
	Update(offset int, data []float32)
	Pointer(start, length, stride int) VertexBufferPointer
	Delete()
}

type VertexBufferPointer interface {
}

type Draw interface {
	SetVertexShader(glsl string)
	SetFragmentShader(glsl string)
	Compile() (CompiledDraw, error)
}

type CompiledDraw interface {
	GetVertexAttributeLocation(name string) int
	GetUniformLocation(name string) int
	NewVertexArrayObject() VertexArrayObject
	NewCall(func(call DrawCall)) image.AcceleratedCall
	Delete()
}

type VertexArrayObject interface {
	SetVertexAttribute(location int, pointer VertexBufferPointer)
	Delete()
}

type DrawCall interface {
	BindVertexArrayObject(VertexArrayObject)
	BindTexture0(img *image.Image)
	SetFloatUniform(location int, val float32)
	SetIntUniform(location int, val int)
	SetMatrix4Uniform(location int, val [16]float32)
	Draw(mode Mode, first, count int)
}

// Mode specifies what kind of primitives to render
type Mode int

const (
	Points Mode = iota
	LineStrip
	LineLoop
	Lines
	TriangleStrip
	TriangleFan
	Triangles
)
