package glfw_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/glblend"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
)

var resolutions = map[string]struct {
	width, height int
}{
	"1920x1080": {
		width:  1920,
		height: 1080,
	},
	"32x32": {
		width:  32,
		height: 32,
	},
}

// 1920x1080 - 90us
// 32x32 	 - 31us
func BenchmarkSource_BlendSourceToTarget(b *testing.B) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()
	context := openGL.Context()
	tool, err := glblend.NewSource(context)
	if err != nil {
		panic(err)
	}
	for name, resolution := range resolutions {
		b.Run(name, func(b *testing.B) {
			source := newImageSelection(openGL, resolution.width, resolution.height)
			target := newImageSelection(openGL, resolution.width, resolution.height)
			source.Image().Upload()
			target.Image().Upload()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tool.BlendSourceToTarget(source, target)
				openGL.ContextAPI().Finish()
			}
		})
	}
}

// 1920x1080 - 82us
// 32x32     - 32us
func BenchmarkSourceOver_BlendSourceToTarget(b *testing.B) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()
	context := openGL.Context()
	tool, err := glblend.NewSourceOver(context)
	if err != nil {
		panic(err)
	}
	for name, resolution := range resolutions {
		b.Run(name, func(b *testing.B) {
			source := newImageSelection(openGL, resolution.width, resolution.height)
			target := newImageSelection(openGL, resolution.width, resolution.height)
			source.Image().Upload()
			target.Image().Upload()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tool.BlendSourceToTarget(source, target)
				openGL.ContextAPI().Finish()
			}
		})
	}
}

func newImageSelection(gl *glfw.OpenGL, width, height int) image.Selection {
	img := gl.NewImage(width, height)
	selection := img.WholeImageSelection()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			selection.SetColor(x, y, image.RGBA(byte(x), byte(y), byte(x), byte(y)))
		}
	}
	return selection
}
