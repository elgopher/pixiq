package main

import (
	"log"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 40, glfw.Title("Use left and right mouse buttons to draw"), glfw.Zoom(20))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// Create mouse instance for window.
		mouseState := mouse.New(window)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
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
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
	})
}

func drawSquare(screen image.Selection, x int, y int, color image.Color) {
	for xOff := -1; xOff <= 1; xOff++ {
		for yOff := -1; yOff <= 1; yOff++ {
			screen.SetColor(x+xOff, y+yOff, color)
		}
	}
}
