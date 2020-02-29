package main

import (
	"log"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
)

// This example shows how too use two separate OpenGL instances. Each contains its
// own windows and images. It can be used to run two totally separated programs.
//
// Please note that this functionality is experimental and may change in the
// near future. Such feature may be harmful for overall performance of Pixiq.
func main() {
	glfw.StartMainThreadLoop(func(mainThreadLoop *glfw.MainThreadLoop) {
		go startOpenGL(mainThreadLoop, "Lime", colornames.Lime)
		startOpenGL(mainThreadLoop, "Pink", colornames.Pink)
	})
}

func startOpenGL(mainThreadLoop *glfw.MainThreadLoop, title string, color image.Color) {
	gl, err := glfw.NewOpenGL(mainThreadLoop)
	if err != nil {
		log.Panicf("%s NewOpenGL failed: %v", title, err)
	}
	defer gl.Destroy()
	win, err := gl.OpenWindow(2, 1, glfw.Zoom(100), glfw.Title(title))
	if err != nil {
		log.Panicf("%s OpenWindow failed: %v", title, err)
	}
	defer win.Close()
	loop.Run(win, func(frame *loop.Frame) {
		screen := frame.Screen()
		screen.SetColor(0, 0, color)
		screen.SetColor(1, 0, color)
		if win.ShouldClose() {
			frame.StopLoopEventually()
		}
	})
}
