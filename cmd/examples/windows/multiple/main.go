package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

var (
	red  = pixiq.RGBA(255, 0, 0, 255)
	blue = pixiq.RGBA(0, 0, 255, 255)
)

// This example shows how to open two windows at the same time.
func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
		windows := gl.Windows()
		redWindow := windows.Open(320, 180)
		blueWindow := windows.Open(250, 90)
		// Start the loop in the background, because Loop method blocks
		// the current goroutine.
		go loops.Loop(redWindow, fillWith(red))
		// Start another one.
		loops.Loop(blueWindow, fillWith(blue))
	})
}

// fillWith returns a function filling whole Screen with specific color.
func fillWith(color pixiq.Color) func(frame *pixiq.Frame) {
	return func(frame *pixiq.Frame) {
		screen := frame.Screen()
		for y := 0; y < screen.Height(); y++ {
			for x := 0; x < screen.Width(); x++ {
				screen.SetColor(x, y, color)
			}
		}
	}
}
