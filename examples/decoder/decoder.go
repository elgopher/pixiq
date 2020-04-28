package main

import (
	"github.com/jacekolszak/pixiq/decoder"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/tools/blend"
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

		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			sourceOverBlender.BlendSourceToTarget(img.WholeImageSelection(), screen)

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})

	})
}
