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
		window := windows.Open(320, 160, opengl.Title("Use WSAD and ESC to close window"))
		keys := keyboard.New(window)
		x := 160
		y := 80
		loops.Loop(window, func(frame *pixiq.Frame) {
			keys.Update()
			if keys.JustPressed(keyboard.A) {
				x -= 10
			}
			if keys.JustPressed(keyboard.D) {
				x += 10
			}
			if keys.JustPressed(keyboard.W) {
				y -= 10
			}
			if keys.JustPressed(keyboard.S) {
				y += 10
			}
			if keys.JustPressed(keyboard.Esc) {
				frame.StopLoopEventually()
			}
			screen := frame.Screen()
			screen.SetColor(x, y, colornames.White)
		})
	})
}
