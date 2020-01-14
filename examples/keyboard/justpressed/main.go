package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL) {
		window, err := gl.OpenWindow(80, 40, opengl.Title("Use WSAD and ESC to close window"), opengl.Zoom(4))
		if err != nil {
			panic(err)
		}
		// Create keyboard instance for window.
		keys := keyboard.New(window)
		x := 40
		y := 20
		loop.Run(window, func(frame *loop.Frame) {
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
			if keys.JustPressed(keyboard.Esc) {
				frame.StopLoopEventually()
			}
			screen := frame.Screen()
			screen.SetColor(x, y, colornames.White)
		})
	})
}
