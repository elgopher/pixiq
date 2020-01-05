package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
		window := gl.Windows().Open(320, 180)
		loops.Loop(window, func(frame *pixiq.Frame) {
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
		window.Close()
	})
}
