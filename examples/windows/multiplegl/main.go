package main

import (
	"log"

	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/image"
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
	for {
		screen := win.Screen()
		screen.SetColor(0, 0, color)
		screen.SetColor(1, 0, color)
		win.Draw()
		if win.ShouldClose() {
			break
		}
	}
}
