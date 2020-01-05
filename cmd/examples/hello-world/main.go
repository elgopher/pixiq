package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

var red = pixiq.RGBA(255, 0, 0, 255)

func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, screens *pixiq.Screens) {
		window := gl.Windows().Open(320, 16)
		screens.Loop(window, func(frame *pixiq.Frame) {
			screen := frame.Screen()
			screen.SetColor(4, 4, red)
			screen.SetColor(5, 4, red)
			screen.SetColor(6, 4, red)
			screen.SetColor(7, 5, red)
			screen.SetColor(7, 6, red)
			screen.SetColor(7, 7, red)
			screen.SetColor(6, 8, red)
			screen.SetColor(5, 8, red)
			screen.SetColor(4, 5, red)
			screen.SetColor(4, 6, red)
			screen.SetColor(4, 7, red)
			screen.SetColor(4, 8, red)
			screen.SetColor(4, 9, red)
			screen.SetColor(4, 10, red)
		})
	})
}
