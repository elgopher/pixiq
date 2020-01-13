package loop_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
)

func BenchmarkScreenLoops_Loop(b *testing.B) {
	b.StopTimer()
	var (
		screen = &noopScreen{image: image.New(1, 1, &acceleratedImageStub{})}
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		frameNumber := 0
		loop.Run(screen, func(frame *loop.Frame) {
			frameNumber++
			if frameNumber > 10000 {
				frame.StopLoopEventually()
			}
		})
	}
}

type noopScreen struct {
	image *image.Image
}

func (d *noopScreen) Image() *image.Image {
	return d.image
}

func (d *noopScreen) Draw() {
}

func (d *noopScreen) SwapImages() {
}
