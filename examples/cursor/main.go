package main

import (
	"log"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(600, 360, glfw.Title("Press left or right mouse button to change cursor look"))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// create image and draw simplified crosshair
		crosshair := crosshair(openGL)
		// create cursor from crosshair selection
		crosshairCursor := openGL.NewCursor(crosshair, glfw.CursorZoom(3))
		// create standard cursor
		ibeamCursor := openGL.NewStandardCursor(glfw.IBeam)

		mouseState := mouse.New(window)
		loop.Run(window, func(frame *loop.Frame) {
			mouseState.Update()
			if mouseState.JustPressed(mouse.Left) {
				window.SetCursor(crosshairCursor)
			}
			if mouseState.JustPressed(mouse.Right) {
				window.SetCursor(ibeamCursor)
			}
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
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
