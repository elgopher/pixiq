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
		window := windows.Open(320, 160)
		keys := keyboard.New(window)
		x := 160
		y := 80
		loops.Loop(window, func(frame *pixiq.Frame) {
			keys.Update()
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
