package blend_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
	"github.com/jacekolszak/pixiq/tools/blend"
)

// 3ms
func BenchmarkSource_BlendSourceToTarget(b *testing.B) {
	var (
		tool   = blend.NewSource()
		width  = 1920
		height = 1080
		source = newImageSelection(width, height)
		target = newImageSelection(width, height)
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.BlendSourceToTarget(source, target)
	}
}

// 12ms
func BenchmarkSourceOver_BlendSourceToTarget(b *testing.B) {
	var (
		tool   = blend.NewSourceOver()
		width  = 1920
		height = 1080
		source = newImageSelection(width, height)
		target = newImageSelection(width, height)
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.BlendSourceToTarget(source, target)
	}
}

// 34ms
func BenchmarkTool_BlendSourceToTarget(b *testing.B) {
	var (
		tool   = blend.New(blenderStub{})
		width  = 1920
		height = 1080
		source = newImageSelection(width, height)
		target = newImageSelection(width, height)
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.BlendSourceToTarget(source, target)
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
