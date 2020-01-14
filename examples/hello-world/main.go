package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	// Use OpenGL on PCs with Linux, Windows and MacOS.
	// This package can open windows and draw images on them.
	opengl.Run(func(gl *opengl.OpenGL) {
		window, err := gl.OpenWindow(80, 16, opengl.Zoom(5))
		if err != nil {
			panic(err)
		}
		// Create a loop for a screen. OpenGL's Window is a Screen (some day
		// in the future Pixiq may support different platforms such as mobile
		// or browser, therefore we need a Screen abstraction).
		// Each iteration of the loop is a Frame.
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			screen.SetColor(40, 8, colornames.White)
		})
	})
}
