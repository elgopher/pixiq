package blend_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/blend"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
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

// 1920x1080 - 3ms
// 32x32     - 2us
func BenchmarkSource_BlendSourceToTarget(b *testing.B) {
	tool := blend.NewSource()
	for name, resolution := range resolutions {
		b.Run(name, func(b *testing.B) {
			source := newImageSelection(resolution.width, resolution.height)
			target := newImageSelection(resolution.width, resolution.height)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tool.BlendSourceToTarget(source, target)
			}
		})
	}

}

// 1920x1080 - 12ms
// 32x32     - 6us
func BenchmarkSourceOver_BlendSourceToTarget(b *testing.B) {
	tool := blend.NewSourceOver()
	for name, resolution := range resolutions {
		b.Run(name, func(b *testing.B) {
			source := newImageSelection(resolution.width, resolution.height)
			target := newImageSelection(resolution.width, resolution.height)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tool.BlendSourceToTarget(source, target)
			}
		})
	}
}

// 1920x1080 - 34ms
// 32x32     - 17us
func BenchmarkTool_BlendSourceToTarget(b *testing.B) {
	tool := blend.New(blenderStub{})
	for name, resolution := range resolutions {
		b.Run(name, func(b *testing.B) {
			source := newImageSelection(resolution.width, resolution.height)
			target := newImageSelection(resolution.width, resolution.height)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tool.BlendSourceToTarget(source, target)
			}
		})
	}
}

type blenderStub struct {
}

func (b blenderStub) BlendSourceToTargetColor(source, target image.Color) image.Color {
	return source
}

func newImageSelection(width, height int) image.Selection {
	img := image.New(fake.NewAcceleratedImage(width, height))
	selection := img.WholeImageSelection()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			selection.SetColor(x, y, image.RGBA(byte(x), byte(y), byte(x), byte(y)))
		}
	}
	return selection
}
