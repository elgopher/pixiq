package internal_test

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/opengl/internal"
)

func BenchmarkKeyboardEvents(b *testing.B) {
	events := internal.KeyboardEvents{}
	for i := 0; i < b.N; i++ {
		events.OnKeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
		events.OnKeyCallback(nil, glfw.KeyB, 0, glfw.Release, 0)
		events.OnKeyCallback(nil, glfw.KeyC, 0, glfw.Press, 0)
		for {
			_, ok := events.Poll()
			if !ok {
				break
			}
		}
	}
}
