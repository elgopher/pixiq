package main

import (
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/goimage"
	"image"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		pixiqImage := openGL.NewImage(4, 4)
		selection := pixiqImage.WholeImageSelection()
		//
		newImage := goimage.FromSelection(selection, goimage.Zoom(3))
		//
		existingImage := image.NewRGBA(image.Rect(0, 0, selection.Width(), selection.Height()))
		goimage.CopyFromSelection(selection, existingImage)
		//
		goimage.CopyToSelection(newImage, selection)
	})
}
