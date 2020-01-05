package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, screens *pixiq.Screens) {
		window := gl.Windows().Open(320, 180)
		screens.Loop(window, func(frame *pixiq.Frame) {
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
		window.Close()
	})
}
