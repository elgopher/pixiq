package pixiq_test

import (
	"testing"

	"github.com/jacekolszak/pixiq"
)

func BenchmarkWindows_Loop(b *testing.B) {
	images := pixiq.NewImages(&fakeAcceleratedImages{})
	windows := pixiq.NewWindows(images)
	win := &noopWindow{}
	for i := 0; i < b.N; i++ {
		frameNumber := 0
		windows.Loop(win, func(frame *pixiq.Frame) {
			frameNumber++
			if frameNumber > 10000 {
				frame.CloseWindowEventually()
			}
		})
	}
}

type noopWindow struct {
}

func (d noopWindow) Draw(*pixiq.Image) {
}

func (d noopWindow) SwapImages() {
}

func (d noopWindow) Close() {
}

func (d noopWindow) Width() int {
	return 0
}

func (d noopWindow) Height() int {
	return 0
}
