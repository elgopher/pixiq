package shader

import "github.com/jacekolszak/pixiq/image"

type GLProgram interface {
	SetVertexShader(glsl string)
	SetFragmentShader(glsl string)
	Compile() (GLCompiledProgram, error)
}

type GLCompiledProgram interface {
	New() GLCall
}

type GLCall interface {
	// data is (x,y) -> (u,v), that is: vertexPosition -> texturePosition
	Set(data []float32)
	SetTexture(name string, img *image.Image)
	SetFloat(name string, val float32)
	SetInt(name string, val int)
	SetMatrix4(name string, val [16]float32)
}

type Program struct {
	program GLProgram
}

func NewProgram(program GLProgram) *Program {
	return &Program{program: program}
}

func (p *Program) AddSelectionUniform(name string) {

}

func (p *Program) SetFragmentShader(shader *FragmentShader) {

}

func (p *Program) SetVertexShader(shader *VertexShader) {

}

func (p *Program) Compille() (*CompiledProgram, error) {
	p.program.SetVertexShader("...")
	p.program.SetFragmentShader("...")
	compiled, _ := p.program.Compile()
	return &CompiledProgram{compiled: compiled}, nil
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
	compiled GLCompiledProgram
}

func (p *CompiledProgram) New() *Call {
	return &Call{call: p.compiled.New()}
}

type Call struct {
	call GLCall
}

func (c *Call) SetSelection(name string, selection image.Selection) {
	c.call.SetTexture("source", selection.Image())
	c.call.SetInt("source_x", selection.ImageX())
	c.call.SetInt("source_y", selection.ImageY())
	c.call.SetInt("source_width", selection.Width())
	c.call.SetInt("source_height", selection.Height())
}

func (c *Call) SetVertices(data []float32) {
	c.call.Set(data)
}
