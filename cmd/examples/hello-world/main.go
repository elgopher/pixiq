package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(acceleratedImages pixiq.AcceleratedImages, systemWindows pixiq.SystemWindows) {
		images := pixiq.NewImages(acceleratedImages)
		windows := pixiq.NewWindows(images, systemWindows)
		windows.New(16, 16).Loop(func(frame *pixiq.Frame) {
			selection := frame.Image().WholeImageSelection()
			red := pixiq.RGBA(255, 0, 0, 255)
			selection.SetColor(4, 4, red)
		})
	})
}
