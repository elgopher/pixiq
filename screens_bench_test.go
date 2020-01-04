package pixiq_test

import (
	"testing"

	"github.com/jacekolszak/pixiq"
)

func BenchmarkScreens_Loop(b *testing.B) {
	b.StopTimer()
	images := pixiq.NewImages(&fakeAcceleratedImages{})
	screens := pixiq.NewScreens(images)
	screen := &noopScreen{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		frameNumber := 0
		screens.Loop(screen, func(frame *pixiq.Frame) {
			frameNumber++
			if frameNumber > 10000 {
				frame.StopLoopEventually()
			}
		})
	}
}

type noopScreen struct {
}

func (d noopScreen) Draw(*pixiq.Image) {
}

func (d noopScreen) SwapImages() {
}

func (d noopScreen) Close() {
}

func (d noopScreen) Width() int {
	return 0
}

func (d noopScreen) Height() int {
	return 0
}
