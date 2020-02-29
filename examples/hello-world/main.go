package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/loop"
)

func main() {
	// Use glfw on PCs with Linux, Windows and MacOS.
	// This package can open windows and draw images on them.
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 16, glfw.Zoom(5))
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
