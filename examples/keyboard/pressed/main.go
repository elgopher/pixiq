package main

import (
	"log"

	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/keyboard"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 40, glfw.Title("Use WSAD and ESC to close window"), glfw.Zoom(4))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// Create keyboard instance for window.
		keys := keyboard.New(window)
		x := 40
		y := 20
		for {
			// Poll keyboard events
			keys.Update()
			// Pressed returns true if given key is currently pressed.
			if keys.Pressed(keyboard.A) {
				x--
			}
			if keys.Pressed(keyboard.D) {
				x++
			}
			if keys.Pressed(keyboard.W) {
				y--
			}
			if keys.Pressed(keyboard.S) {
				y++
			}
			if keys.Pressed(keyboard.Esc) || window.ShouldClose() {
				break
			}
			window.Screen().SetColor(x, y, colornames.White)
			window.Draw()
		}
	})
}
