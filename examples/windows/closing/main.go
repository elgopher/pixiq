package main

import (
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

// This example shows how to properly close the window.
func main() {
	opengl.Run(func(gl *opengl.OpenGL) {
		window := gl.Open(320, 180)
		// clean resources when function ends
		defer window.Close()
		loop.Run(window, func(frame *loop.Frame) {
			// If window was closed by the user ShouldClose will return true
			if window.ShouldClose() {
				// Stop the loop at the end of the iteration
				frame.StopLoopEventually()
			}
		})
	})
}
