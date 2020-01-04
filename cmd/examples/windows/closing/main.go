package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.StartMainThreadLoop(func(loop *opengl.MainThreadLoop) {
		openGL := opengl.New(loop)
		windows := openGL.Windows()
		window := windows.Open(320, 180)
		images := pixiq.NewImages(openGL.AcceleratedImages())
		screens := pixiq.NewScreens(images)
		screens.Loop(window, func(frame *pixiq.Frame) {
			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
		window.Close()
	})
}
