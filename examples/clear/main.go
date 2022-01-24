package main

import (
	"github.com/elgopher/pixiq/clear"
	"github.com/elgopher/pixiq/colornames"
	"github.com/elgopher/pixiq/glclear"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/image"
	"github.com/elgopher/pixiq/keyboard"
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
		tools := []clearTool{
			glclear.New(openGL.Context()), // GPU one
			clear.New(),                   // CPU one
		}
		tools[0].SetColor(colornames.Cornflowerblue)
		tools[1].SetColor(colornames.Hotpink)
		currentTool := 0
		keys := keyboard.New(window)
		screen := window.Screen()
		for {
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

			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}
