package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glclear"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/loop"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(2, 2, glfw.Zoom(100))
		if err != nil {
			panic(err)
		}
		context := openGL.Context()
		tool := glclear.New(context.NewClearCommand())
		tool.SetColor(colornames.Yellow)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			tool.Clear(screen.Selection(-1, -1).WithSize(1, 1))

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})

	})
}
