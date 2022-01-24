package main

import (
	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glfw"
)

func main() {
	// Use glfw on PCs with Linux, Windows and MacOS.
	// This package can open windows and draw images on them.
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 16, glfw.Zoom(5))
		if err != nil {
			panic(err)
		}
		// Draw window contents (screen) in the loop.
		for {
			screen := window.Screen()
			screen.SetColor(40, 8, colornames.White)
			// Draw will draw the screen and make changes visible to the user
			window.Draw()
			// If window was closed by the user ShouldClose will return true
			if window.ShouldClose() {
				// Stop the loop at the end of the iteration
				break
			}
		}
	})
}
