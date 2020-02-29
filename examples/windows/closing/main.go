package main

import (
	"log"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/loop"
)

// This example shows how to properly close the window.
func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(320, 180)
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
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
