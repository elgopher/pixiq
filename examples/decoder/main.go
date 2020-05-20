package main

import (
	"github.com/jacekolszak/pixiq/blend"
	"github.com/jacekolszak/pixiq/decoder"
	"github.com/jacekolszak/pixiq/glfw"
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
