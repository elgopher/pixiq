package main

import (
	"log"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 40, glfw.Title("Use left and right mouse buttons"), glfw.Zoom(4))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// Create mouse instance for window.
		mouseState := mouse.New(window)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			// Poll mouse events
			mouseState.Update()
			pos := mouseState.Position()
			// Pressed returns true if given key is currently pressed.
			if mouseState.Pressed(mouse.Left) {
				screen.SetColor(pos.X(), pos.Y(), colornames.White)
			}
			if mouseState.Pressed(mouse.Right) {
				screen.SetColor(pos.X(), pos.Y(), colornames.Black)
			}
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
	})
}
