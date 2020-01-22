package blender

import (
	"github.com/jacekolszak/pixiq/glsl/shader"
	"github.com/jacekolszak/pixiq/image"
)

type GL interface {
	shader.VertexShaderCompiler
	shader.FragmentShaderCompiler
	DrawTriangles(vertices []float32, vertexShader shader.VertexShader, fragmentShader shader.FragmentShader) shader.Call
}

func CompileImageBlender(gl GL) (*ImageBlender, error) {
	if gl == nil {
		panic("nil drawer")
	}
	fragmentShader, _ := gl.CompileFragmentShader(
		"color=vec4(source.get(x,y))")
	// TODO Here we can verify that blending really works on this platform
	// by executing blend. If it does not we can return error. This is much
	// better than panicking during real Blending
	return &ImageBlender{shader: fragmentShader, gl: gl}, nil
}

type ImageBlender struct {
	shader shader.FragmentShader
	gl     GL
}

func (b *ImageBlender) Blend(source, target image.Selection) {
	call := b.gl.DrawTriangles(nil, nil, b.shader)
	call.SetSelection("source", source)
	target.Modify(call)
}
