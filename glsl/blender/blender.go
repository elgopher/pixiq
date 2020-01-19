package blender

import (
	"github.com/jacekolszak/pixiq/glsl/shader"
	"github.com/jacekolszak/pixiq/image"
)

func CompileImageBlender(compiler shader.Compiler) (*ImageBlender, error) {
	if compiler == nil {
		panic("nil compiler")
	}
	// TODO Here we can verify that blending really work on this platform
	// by executing blend. If it does not we can return error. This is much
	// better than panicking during real Blending
	return &ImageBlender{}, nil
}

// TODO Refactor this a bit - it should stick to the global interface
// for all sorts of blending to support Go implicit interfaces:
//
//   type Blender interface {
//	   Blend(source, target image.Selection)
//   }
//
// Eventually blend function can be run using Blend(source,target)
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
