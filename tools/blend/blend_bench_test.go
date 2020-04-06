package blend_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
	"github.com/jacekolszak/pixiq/tools/blend"
)

func BenchmarkSource_BlendSourceToTarget(b *testing.B) {
	var (
		tool   = blend.NewSource()
		width  = 1920
		height = 1080
		source = image.New(width, height, fake.NewAcceleratedImage(width, height)).
			WholeImageSelection()
		target = image.New(width, height, fake.NewAcceleratedImage(width, height)).
			WholeImageSelection()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.BlendSourceToTarget(source, target)
	}
}

func BenchmarkTool_BlendSourceToTarget(b *testing.B) {
	var (
		tool   = blend.New(blenderStub{})
		width  = 1920
		height = 1080
		source = image.New(width, height, fake.NewAcceleratedImage(width, height)).
			WholeImageSelection()
		target = image.New(width, height, fake.NewAcceleratedImage(width, height)).
			WholeImageSelection()
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
