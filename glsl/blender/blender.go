package blender

import (
	"github.com/jacekolszak/pixiq/glsl/shader"
	"github.com/jacekolszak/pixiq/image"
)

func CompileImageBlender(compiler shader.Compiler) (*ImageBlender, error) {
	if compiler == nil {
		panic("nil compiler")
	}
	return &ImageBlender{}, nil
}

type ImageBlender struct {
}

func (b ImageBlender) Source(selection image.Selection) ImageBlenderCall {
	return ImageBlenderCall{}
}

type ImageBlenderCall struct {
}

func (c ImageBlenderCall) Into(selection image.Selection) {
	selection.SetColor(0, 0, image.RGB(255, 255, 255))
}
