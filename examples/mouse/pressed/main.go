package main

import (
	"log"

	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/image"
	"github.com/elgopher/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 40, glfw.Title("Use left and right mouse buttons to draw"), glfw.Zoom(20))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// Create mouse instance for window.
		mouseState := mouse.New(window)
		for {
			screen := window.Screen()
			// Poll mouse events
			mouseState.Update()
			// Get cursor position
			pos := mouseState.Position()
			// Pressed returns true if given key is currently pressed.
			if mouseState.Pressed(mouse.Left) {
				// pos.X() and pos.Y() returns position in pixel dimensions
				drawSquare(screen, pos.X(), pos.Y(), colornames.White)
			}
			if mouseState.Pressed(mouse.Right) {
				drawSquare(screen, pos.X(), pos.Y(), colornames.Black)
			}
			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}

func drawSquare(screen image.Selection, x int, y int, color image.Color) {
	for xOff := -1; xOff <= 1; xOff++ {
		for yOff := -1; yOff <= 1; yOff++ {
			screen.SetColor(x+xOff, y+yOff, color)
		}
	}
}
