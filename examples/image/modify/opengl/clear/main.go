package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		window, err := gl.OpenWindow(10, 10, opengl.Zoom(20))
		if err != nil {
			panic(err)
		}
		context := gl.Context()
		clear := context.ClearCommand()
		clear.Color = &colornames.Indianred
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			var (
				leftEye  = screen.Selection(2, 2).WithSize(2, 2)
				rightEye = screen.Selection(6, 2).WithSize(2, 2)
				nose     = screen.Selection(4, 5).WithSize(2, 2)
				mouth    = screen.Selection(2, 8).WithSize(6, 1)
			)
			leftEye.Modify(clear)
			rightEye.Modify(clear)
			nose.Modify(clear)
			mouth.Modify(clear)
		})

	})
}
