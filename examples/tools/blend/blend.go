package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/tools/blend"
)

type blender interface {
	BlendSourceToTarget(source, target image.Selection)
}

func main() {
	glfw.RunOrDie(func(gl *glfw.OpenGL) {
		window, err := gl.OpenWindow(60, 60, glfw.Zoom(5))
		if err != nil {
			panic(err)
		}
		tools := []blender{
			// TODO Add GPU
			blend.NewSource(),
			blend.NewSourceOver(),
		}
		currentTool := 0

		face := face(gl)

		keys := keyboard.New(window)
		position := position{keyboard: keys}

		loop.Run(window, func(frame *loop.Frame) {
			keys.Update()
			if keys.JustReleased(keyboard.Space) {
				currentTool++
				currentTool = currentTool % len(tools)
			}
			position.update()
			tool := tools[currentTool]
			// face will be blended into screen at a given position
			target := frame.Screen().Selection(position.x, position.y)
			tool.BlendSourceToTarget(face, target)

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})

	})
}

type position struct {
	x, y     int
	keyboard *keyboard.Keyboard
}

func (p *position) update() {
	if p.keyboard.Pressed(keyboard.Left) {
		p.x--
	}
	if p.keyboard.Pressed(keyboard.Right) {
		p.x++
	}
	if p.keyboard.Pressed(keyboard.Up) {
		p.y--
	}
	if p.keyboard.Pressed(keyboard.Down) {
		p.y++
	}
}

func face(gl *glfw.OpenGL) image.Selection {
	var (
		img       = gl.NewImage(10, 10)
		selection = img.WholeImageSelection()
		color     = colornames.Lightyellow
	)
	selection.SetColor(2, 2, color)
	selection.SetColor(4, 2, color)
	selection.SetColor(3, 4, color)
	selection.SetColor(2, 6, color)
	selection.SetColor(3, 6, color)
	selection.SetColor(4, 6, color)
	return selection
}
