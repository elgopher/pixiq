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
			// 	JustPressed is true if A was pressed between two last keys.Update() calls
			if keys.JustPressed(keyboard.A) {
				x -= 2
			}
			if keys.JustPressed(keyboard.D) {
				x += 2
			}
			if keys.JustPressed(keyboard.W) {
				y -= 2
			}
			if keys.JustPressed(keyboard.S) {
				y += 2
			}
			if keys.JustPressed(keyboard.Esc) || window.ShouldClose() {
				break
			}
			window.Screen().SetColor(x, y, colornames.White)
			window.Draw()
		}
	})
}
