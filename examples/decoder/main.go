package main

import (
	"github.com/elgopher/pixiq/blend"
	"github.com/elgopher/pixiq/decoder"
	"github.com/elgopher/pixiq/glfw"
)

func main() {
	glfw.RunOrDie(func(gl *glfw.OpenGL) {
		window, err := gl.OpenWindow(312, 240, glfw.Zoom(3))
		if err != nil {
			panic(err)
		}

		sourceOverBlender := blend.NewSourceOver()

		imageDecoder := decoder.New(gl)
		img, err := imageDecoder.DecodeFile("docs/pixiq-primitives.gif")
		if err != nil {
			panic(err)
		}

		for {
			sourceOverBlender.BlendSourceToTarget(img.WholeImageSelection(), window.Screen())

			window.Draw()
			if window.ShouldClose() {
				break
			}
		}

	})
}
