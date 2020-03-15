package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/tools/clear"
	"github.com/jacekolszak/pixiq/tools/glclear"
)

type clearTool interface {
	SetColor(image.Color)
	Clear(image.Selection)
}

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(10, 10, glfw.Zoom(30))
		if err != nil {
			panic(err)
		}
		context := openGL.Context()
		tools := []clearTool{
			glclear.New(context.NewClearCommand()), // GPU one
			clear.New(),                            // CPU one
		}
		tools[0].SetColor(colornames.Cornflowerblue)
		tools[1].SetColor(colornames.Hotpink)
		currentTool := 0
		keys := keyboard.New(window)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			var (
				leftEye  = screen.Selection(2, 2).WithSize(2, 2)
				rightEye = screen.Selection(6, 2).WithSize(2, 2)
				nose     = screen.Selection(4, 5).WithSize(2, 2)
				mouth    = screen.Selection(2, 8).WithSize(6, 1)
			)
			tool := tools[currentTool]
			tool.Clear(leftEye)
			tool.Clear(rightEye)
			tool.Clear(nose)
			tool.Clear(mouth)
			keys.Update()
			if keys.JustReleased(keyboard.Space) {
				currentTool++
				currentTool = currentTool % len(tools)
			}

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})

	})
}
