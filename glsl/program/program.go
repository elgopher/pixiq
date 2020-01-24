package program

import (
	"github.com/jacekolszak/pixiq/image"
)

type VertexBuffer struct {
}

func (b *VertexBuffer) AddFloat(val float32) {

}

func (b *VertexBuffer) AddFloat2(val1 float32, val2 float32) {

}

func (b *VertexBuffer) AddFloat3(val1 float32, val2 float32) {

}

type Program struct {
	program    Draw
	parameters []Parameter
}

type Parameter struct {
	Name string
	Type ParameterType
}

type ParameterType int

func New(program Draw) *Program {
	return &Program{program: program}
}

func (p *Program) AddSelectionParameter(name string) {

}

func (p *Program) SetFragmentShader(shader *FragmentShader) {

}

func (p *Program) SetVertexShader(shader *VertexShader) {

}

func (p *Program) Compille() (*CompiledProgram, error) {
	p.program.SetVertexShader("...")
	p.program.SetFragmentShader("...")
	compiled, _ := p.program.Compile()
	parameterIndices := map[string]int{}
	for _, param := range p.parameters {
		parameterIndices[param.Name] = compiled.GetUniformLocation(param.Name)
	}

	return &CompiledProgram{compiled: compiled, parameterIndices: parameterIndices}, nil
}

func NewFragmentShader() *FragmentShader {
	return &FragmentShader{}
}

type FragmentShader struct {
}

func (s *FragmentShader) SetMain(sourceCode string) {

}

func NewVertexShader() *VertexShader {
	return &VertexShader{}
}

type VertexShader struct {
}

func (s *VertexShader) SetMain(sourceCode string) {

}

type CompiledProgram struct {
	compiled         CompiledDraw
	parameterIndices map[string]int
}

func (p *CompiledProgram) SetVertexFormat(format VertexFormat) {
}

func (p *CompiledProgram) NewCall(f func(call HighLevelCall)) image.AcceleratedCall {
	return nil
}

type VertexFormat struct {
	values []VertexValue
}

// 1 x float32
func (f VertexFormat) AddFloat(name string) {
}

// 2 x float32
func (f VertexFormat) AddFloat2(name string) {
}

// 3 x float32
func (f VertexFormat) AddFloat3(name string) {
}

// 4 x float32
func (f VertexFormat) AddFloat4(name string) {
}

func (f VertexFormat) AddByte(name string) {
}

type VertexValue struct {
	Index      int
	VertexType VertexType
	// Size specifies the number of components per generic vertex attribute.
	// Must be 1, 2, 3, 4. The initial value is 4.
	Size int
	// Stride specifies the byte offset between consecutive generic vertex attributes.
	// If stride is 0, the generic vertex attributes are understood to be tightly
	// packed in the array. The initial value is 0.
	Stride int
	// Specifies a pointer to the first generic vertex attribute in the array.
	// If a non-zero buffer is currently bound to the GL_ARRAY_BUFFER target,
	// pointer specifies an offset of into the array in the data store of that buffer.
	Offset int
}

type VertexType int // TODO FLOAT,

const (
	Float VertexType = iota
)

type HighLevelCall struct {
	VertexBuffer
	DrawCall
}

func (c *HighLevelCall) SetSelection(name string, selection image.Selection) {
	// save in some map or something
}
