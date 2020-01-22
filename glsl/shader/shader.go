package shader

import "github.com/jacekolszak/pixiq/image"

type FragmentShaderCompiler interface {
	CompileFragmentShader(glsl string) (FragmentShader, error)
}
type VertexShaderCompiler interface {
	CompileVertexShader(glsl string) (VertexShader, error)
}
type FragmentShader interface {
}

type VertexShader interface {
}

type Call struct {
	underlyingCall UnderlyingCall
}

func NewCall(underlyingCall UnderlyingCall) Call {
	return Call{
		underlyingCall: underlyingCall,
	}
}

type UnderlyingCall interface {
	SetTexture(name string, img *image.Image)
	SetFloat(name string, val float32)
	SetInt(name string, val int)
	SetMatrix4(name string, val [16]float32)
}

func (c *Call) SetSelection(name string, selection image.Selection) {
	c.underlyingCall.SetTexture("source", selection.Image())
	c.underlyingCall.SetInt("source_x", selection.ImageX())
	c.underlyingCall.SetInt("source_y", selection.ImageY())
	c.underlyingCall.SetInt("source_width", selection.Width())
	c.underlyingCall.SetInt("source_height", selection.Height())
	// also add functions to sample the selection
}
