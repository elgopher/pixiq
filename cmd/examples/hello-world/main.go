package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

var white = pixiq.RGBA(255, 255, 255, 255)

func main() {
	// Use OpenGL on PCs with Linux, Windows and MacOS. This package can open windows and draw images on them.
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
		windows := gl.Windows()
		window := windows.Open(320, 16)
		// Create a main loop for a screen. OpenGL's Window is a Screen (some day in the future Pixiq may support
		// different platforms such as mobile or browser, therefore we need a Screen abstraction).
		// Each iteration of the loop is a Frame.
		loops.Loop(window, func(frame *pixiq.Frame) {
			screen := frame.Screen()
			screen.SetColor(160, 8, white)
		})
	})
}
