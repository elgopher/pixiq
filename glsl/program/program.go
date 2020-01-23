package program

import "github.com/jacekolszak/pixiq/image"

type Draw interface {
	SetVertexShader(glsl string)
	SetFragmentShader(glsl string)
	Compile() (CompiledDraw, error)
}

type CompiledDraw interface {
	GetVertexIndex(name string)
	GetParameterIndex(name string) int
	SetVertexFormat(values []VertexValue)
	New() DrawCall
}

type DrawCall interface {
	SetVertexBuffer(buffer VertexBuffer)
	SetTexture(index int, img *image.Image)
	SetFloat(index int, val float32)
	SetInt(index int, val int)
	SetMatrix4(index int, val [16]float32)
}

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
		parameterIndices[param.Name] = compiled.GetParameterIndex(param.Name)
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
	p.compiled.SetVertexFormat(format.values)
}

func (p *CompiledProgram) NewCall() *Call {
	return &Call{call: p.compiled.New(), parameterIndices: p.parameterIndices}
}

type Call struct {
	call             DrawCall
	parameterIndices map[string]int
	buffer           VertexBuffer
}

func (c *Call) SetSelection(name string, selection image.Selection) {
	c.call.SetTexture(c.parameterIndices[name], selection.Image())
	c.call.SetInt(c.parameterIndices[name+"_x"], selection.ImageX())
	c.call.SetInt(c.parameterIndices[name+"_y"], selection.ImageY())
	c.call.SetInt(c.parameterIndices[name+"_width"], selection.Width())
	c.call.SetInt(c.parameterIndices[name+"_height"], selection.Height())
}

func (c *Call) SetVertexBuffer(buffer VertexBuffer) {
	c.buffer = buffer
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
	Size       int
	VertexType VertexType
	Stride     int
	Offset     int
}

type VertexType int // TODO FLOAT,
