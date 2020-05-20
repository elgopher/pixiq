package main

import (
	"log"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(600, 360, glfw.Title("Press left or right mouse button to change cursor look"), glfw.Zoom(0))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// create image and draw simplified crosshair
		crosshair := crosshair(openGL)
		// create cursor from crosshair selection
		crosshairCursor := openGL.NewCursor(crosshair, glfw.CursorZoom(6), glfw.Hotspot(1, 1))
		// create standard cursor
		ibeamCursor := openGL.NewStandardCursor(glfw.Hand)

		mouseState := mouse.New(window)

		for {
			mouseState.Update()
			if mouseState.JustPressed(mouse.Left) {
				window.SetCursor(crosshairCursor)
			}
			if mouseState.JustPressed(mouse.Right) {
				window.SetCursor(ibeamCursor)
			}

			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}

func crosshair(openGL *glfw.OpenGL) image.Selection {
	cursorImage := openGL.NewImage(3, 3)
	selection := cursorImage.WholeImageSelection()
	selection.SetColor(1, 0, colornames.Lime)
	selection.SetColor(0, 1, colornames.Lime)
	selection.SetColor(1, 1, colornames.Lime)
	selection.SetColor(2, 1, colornames.Lime)
	selection.SetColor(1, 2, colornames.Lime)
	return selection
}
