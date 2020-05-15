package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glclear"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/loop"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(10, 10, glfw.Zoom(30))
		if err != nil {
			panic(err)
		}
		context := openGL.Context()
		tool := glclear.New(context.NewClearCommand())
		tool.SetColor(colornames.Yellow)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			tool.Clear(screen.Selection(-2, -2).WithSize(4, 4))
			tool.Clear(screen.Selection(8, 8).WithSize(4, 4))

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})

	})
}
