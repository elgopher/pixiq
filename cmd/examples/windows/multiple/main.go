package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

var (
	red  = pixiq.RGBA(255, 0, 0, 255)
	blue = pixiq.RGBA(0, 0, 255, 255)
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, screens *pixiq.Screens) {
		windows := gl.Windows()
		redWindow := windows.Open(320, 180)
		blueWindow := windows.Open(250, 90)
		go screens.Loop(redWindow, fillWith(red))
		screens.Loop(blueWindow, fillWith(blue))
	})
}

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
