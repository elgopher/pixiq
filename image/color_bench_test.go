package image_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
)

func BenchmarkRGBAi(b *testing.B) {
	var (
		red   = 557
		green = -867
		blue  = 612
		alpha = -403
	)
	benchmarkRGBAi(b, red, green, blue, alpha)
}

func benchmarkRGBAi(b *testing.B, red, green, blue, alpha int) image.Color {
	var c image.Color
	for i := 0; i < b.N; i++ {
		c = image.RGBAi(red, green, blue, alpha)
	}
	return c
}
