package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/opengl"
)

// This example shows how too use two separate OpenGL instances. Each contains its
// own windows and images. It can be used to run two totally separated programs.
//
// Please note that this functionality is experimental and may change in the
// near future. Such feature may be harmful for overall performance of Pixiq.
func main() {
	opengl.StartMainThreadLoop(func(mainThreadLoop *opengl.MainThreadLoop) {
		go startOpenGL(mainThreadLoop, "Lime", colornames.Lime)
		startOpenGL(mainThreadLoop, "Pink", colornames.Pink)
	})
}

func startOpenGL(loop *opengl.MainThreadLoop, title string, color pixiq.Color) {
	openGL := opengl.New(loop)
	defer openGL.Destroy()
	windows := openGL.Windows()
	acceleratedImages := openGL.AcceleratedImages()
	images := pixiq.NewImages(acceleratedImages)
	win := windows.Open(2, 1, opengl.Zoom(100), opengl.Title(title))
	defer win.Close()
	screenLoops := pixiq.NewScreenLoops(images)
	screenLoops.Loop(win, func(frame *pixiq.Frame) {
		screen := frame.Screen()
		screen.SetColor(0, 0, color)
		screen.SetColor(1, 0, color)
		if win.ShouldClose() {
			frame.StopLoopEventually()
		}
	})
}
