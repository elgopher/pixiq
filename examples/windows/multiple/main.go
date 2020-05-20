package main

import (
	"github.com/jacekolszak/pixiq/decoder"
	"github.com/jacekolszak/pixiq/glblend"
	"github.com/jacekolszak/pixiq/glfw"
	"log"
)

// This example shows how to open two windows at the same time.
//
// Please note that this functionality is EXPERIMENTAL and may change in the
// near future. Such feature may be harmful for overall performance of Pixiq.
//
// Please note also that at the moment most Pixiq primitives are not goroutine-safe.
// You can't safely use one *glfw.OpenGL instance in two different go-routines.
func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		redWindow, err := openGL.OpenWindow(11, 11, glfw.Zoom(20), glfw.Title("red"))
		if err != nil {
			log.Panicf("red OpenWindow failed: %v", err)
		}
		blueWindow, err := openGL.OpenWindow(11, 11, glfw.Zoom(20), glfw.Title("blue"))
		if err != nil {
			log.Panicf("blue OpenWindow failed: %v", err)
		}

		imageDecoder := decoder.New(openGL)

		redImg, err := imageDecoder.DecodeFile("examples/windows/multiple/red.png")
		if err != nil {
			log.Panicf("DecodeFile failed: %v", err)
		}

		blueImg, err := imageDecoder.DecodeFile("examples/windows/multiple/blue.png")
		if err != nil {
			log.Panicf("DecodeFile failed: %v", err)
		}

		blender, err := glblend.NewSourceOver(openGL.Context())
		if err != nil {
			log.Panicf("glblend.NewSourceOver failed: %v", err)
		}

		for {
			// Draw windows sequentially. Drawing it in parallel does not make much sense with
			// current glfw.OpenGL implementation because on the end all OpenGL calls are executed
			// in a main thread. glfw.OpenGL is also not goroutines-safe and you can't
			// run commands (such as those used by glblend package) atomically.
			blender.BlendSourceToTarget(redImg.WholeImageSelection(), redWindow.Screen())
			redWindow.Draw()
			blender.BlendSourceToTarget(blueImg.WholeImageSelection(), blueWindow.Screen())
			blueWindow.Draw()
			if redWindow.ShouldClose() || blueWindow.ShouldClose() {
				return
			}
		}

	})
}
