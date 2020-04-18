package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/tools/blend"
	"github.com/jacekolszak/pixiq/tools/clear"
	"github.com/jacekolszak/pixiq/tools/glblend"
)

func main() {
	glfw.RunOrDie(func(gl *glfw.OpenGL) {
		window, err := gl.OpenWindow(37, 40, glfw.Zoom(10))
		if err != nil {
			panic(err)
		}

		face := face(gl)

		sourceBlender := blend.NewSource()
		glSourceBlender, err := glblend.NewSource(gl.Context())
		if err != nil {
			panic(err)
		}
		sourceOverBlender := blend.NewSourceOver()

		clearTool := clear.New()
		clearTool.SetColor(colornames.Aliceblue)

		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			// first clear the screen with some opaque color
			clearTool.Clear(screen)

			// source blending overrides the target with source colors
			// fully transparent pixels (RGBA 0x00000000) are rendered as black on the screen.
			sourceBlender.BlendSourceToTarget(face, screen.Selection(10, 7))

			// similar source blending using video card
			glSourceBlender.BlendSourceToTarget(face, screen.Selection(20, 7))

			// source-over blending mixes source and target colors together taking
			// into account alpha channels of both. In places where source has
			// transparent pixels the original target colors are preserved.
			sourceOverBlender.BlendSourceToTarget(face, screen.Selection(10, 24))

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})

	})
}

func face(gl *glfw.OpenGL) image.Selection {
	var (
		img       = gl.NewImage(7, 9)
		selection = img.WholeImageSelection()
		color     = colornames.Violet
	)
	selection.SetColor(2, 2, color)
	selection.SetColor(4, 2, color)
	selection.SetColor(3, 4, color)
	selection.SetColor(2, 6, color)
	selection.SetColor(3, 6, color)
	selection.SetColor(4, 6, color)
	return selection
}
