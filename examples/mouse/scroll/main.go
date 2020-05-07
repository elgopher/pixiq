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
		window, err := openGL.OpenWindow(100, 40, glfw.Title("Move mouse wheel in all possible directions"), glfw.Zoom(20))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// TODO Hide cursor
		x := 50
		y := 20
		// Create mouse instance for window.
		mouseState := mouse.New(window)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			// Poll mouse events
			mouseState.Update()
			scroll := mouseState.Scroll()
			if scroll.X() < 0 {
				x++
			} else if scroll.X() > 0 {
				x--
			}
			if scroll.Y() < 0 {
				y++
			} else if scroll.Y() > 0 {
				y--
			}
			screen.SetColor(x, y, colornames.Azure)
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
	})
}
