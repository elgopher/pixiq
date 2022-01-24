package clear_test

import (
	"testing"

	"github.com/elgopher/pixiq/clear"
	"github.com/elgopher/pixiq/image"
	"github.com/elgopher/pixiq/image/fake"
)

func BenchmarkTool_Clear(b *testing.B) {
	var (
		width     = 1920
		height    = 1080
		img       = image.New(fake.NewAcceleratedImage(width, height))
		selection = img.WholeImageSelection()
		tool      = clear.New()
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.Clear(selection)
	}
}
