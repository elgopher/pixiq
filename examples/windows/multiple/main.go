package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

// This example shows how to open two windows at the same time.
//
// Please note that this functionality is experimental and may change in the
// near future. Such feature may be harmful for overall performance of Pixiq.
func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		redWindow, err := gl.OpenWindow(320, 180, opengl.Title("red"))
		if err != nil {
			panic(err)
		}
		blueWindow, err := gl.OpenWindow(250, 90, opengl.Title("blue"))
		if err != nil {
			panic(err)
		}
		// Start the loop in the background, because Loop method blocks
		// the current goroutine.
		go loop.Run(redWindow, fillWith(colornames.Red))
		// Start another one.
		loop.Run(blueWindow, fillWith(colornames.Blue))
	})
}

// fillWith returns a function filling whole Screen with specific color.
func fillWith(color image.Color) func(frame *loop.Frame) {
	return func(frame *loop.Frame) {
		screen := frame.Screen()
		for y := 0; y < screen.Height(); y++ {
			for x := 0; x < screen.Width(); x++ {
				screen.SetColor(x, y, color)
			}
		}
	}
}
