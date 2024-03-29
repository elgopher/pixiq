package main

import (
	"github.com/elgopher/pixiq/blend"
	"github.com/elgopher/pixiq/clear"
	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glblend"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/image"
)

func main() {
	glfw.RunOrDie(func(gl *glfw.OpenGL) {
		window, err := gl.OpenWindow(37, 40, glfw.Zoom(10), glfw.Title("Blend window"))
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
		glSourceOverBlender, err := glblend.NewSourceOver(gl.Context())
		if err != nil {
			panic(err)
		}

		clearTool := clear.New()
		clearTool.SetColor(colornames.Aliceblue)

		for {
			screen := window.Screen()
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

			// similar source-over blending using video card
			glSourceOverBlender.BlendSourceToTarget(face, screen.Selection(20, 24))

			window.Draw()
			if window.ShouldClose() {
				break
			}
		}

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
