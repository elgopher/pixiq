package main

import (
	"log"

	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 40, glfw.Title("Move mouse wheel in all possible directions"), glfw.Zoom(20))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		x := 40
		y := 20
		// Create mouse instance for window.
		mouseState := mouse.New(window)
		for {
			screen := window.Screen()
			// Poll mouse events
			mouseState.Update()
			scroll := mouseState.Scroll()
			x += int(scroll.X())
			y += int(scroll.Y())
			screen.SetColor(x, y, colornames.Azure)
			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}
