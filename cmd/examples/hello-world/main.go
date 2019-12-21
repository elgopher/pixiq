package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(images *pixiq.Images) {
		image := images.New(16, 16)
		selection := image.WholeImageSelection()
		red := pixiq.RGBA(255, 0, 0, 255)
		selection.SetColor(4, 4, red)
	})
}
