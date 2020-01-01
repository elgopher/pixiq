package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.StartMainThreadLoop(func(loop *opengl.MainThreadLoop) {
		gl := opengl.New(loop)
		images := pixiq.NewImages(gl.AcceleratedImages())
		windows := pixiq.NewWindows(images, gl.SystemWindows())
		windows.New(16, 16).Loop(func(frame *pixiq.Frame) {
			selection := frame.Image().WholeImageSelection()
			red := pixiq.RGBA(255, 0, 0, 255)
			selection.SetColor(4, 4, red)
			selection.SetColor(5, 4, red)
			selection.SetColor(6, 4, red)
			selection.SetColor(7, 5, red)
			selection.SetColor(7, 6, red)
			selection.SetColor(7, 7, red)
			selection.SetColor(6, 8, red)
			selection.SetColor(5, 8, red)
			selection.SetColor(4, 5, red)
			selection.SetColor(4, 6, red)
			selection.SetColor(4, 7, red)
			selection.SetColor(4, 8, red)
			selection.SetColor(4, 9, red)
			selection.SetColor(4, 10, red)
		})
	})
}
