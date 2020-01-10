package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
		windows := gl.Windows()
		window := windows.Open(80, 40, opengl.Title("Use WSAD and ESC to close window"), opengl.Zoom(4))
		// Create keyboard instance for window.
		keys := keyboard.New(window)
		x := 40
		y := 20
		loops.Loop(window, func(frame *pixiq.Frame) {
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
			if keys.Pressed(keyboard.Esc) {
				frame.StopLoopEventually()
			}
			screen := frame.Screen()
			screen.SetColor(x, y, colornames.White)
		})
	})
}
