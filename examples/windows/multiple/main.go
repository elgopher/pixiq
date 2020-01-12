package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/opengl"
)

// This example shows how to open two windows at the same time.
//
// Please note that this functionality is experimental and may change in the
// near future. Such feature may be harmful for overall performance of Pixiq.
func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
		windows := gl.Windows()
		redWindow := windows.Open(320, 180, opengl.Title("red"))
		blueWindow := windows.Open(250, 90, opengl.Title("blue"))
		// Start the loop in the background, because Loop method blocks
		// the current goroutine.
		go loops.Loop(redWindow, fillWith(colornames.Red))
		// Start another one.
		loops.Loop(blueWindow, fillWith(colornames.Blue))
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
