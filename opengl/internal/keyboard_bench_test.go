package internal_test

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/opengl/internal"
)

// Should be 0 allocs/op
func BenchmarkKeyboardEvents(b *testing.B) {
	b.StopTimer()
	events := internal.NewKeyboardEvents(8)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		const numberOfEvents = 8
		for i := 0; i < numberOfEvents; i++ {
			events.OnKeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
		}
		for {
			_, ok := events.Poll()
			if !ok {
				break
			}
		}
	}
}
