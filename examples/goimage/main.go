package main

import (
	"image"

	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/goimage"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		pixiqImage := openGL.NewImage(4, 4)
		selection := pixiqImage.WholeImageSelection()
		// Create a new standard Go image.Image from Pixiq Selection
		newImage := goimage.FromSelection(selection, goimage.Zoom(3))
		// Or fill existing standard Go image.Image
		existingImage := image.NewRGBA(image.Rect(0, 0, selection.Width(), selection.Height()))
		goimage.FillWithSelection(existingImage, selection)
		// Copy standard Go image.Image to Selection
		goimage.CopyToSelection(newImage, selection)
	})
}
