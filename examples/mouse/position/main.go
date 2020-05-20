package main

import (
	"log"

	"github.com/jacekolszak/pixiq/clear"
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/mouse"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		window, err := openGL.OpenWindow(80, 20, glfw.Title("Move mouse left and right"), glfw.Zoom(7))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// Create mouse instance for window.
		mouseState := mouse.New(window)
		x := 40
		clearTool := clear.New()
		for {
			screen := window.Screen()
			clearTool.Clear(screen)
			// Poll mouse events
			mouseState.Update()
			x += mouseState.PositionChange().X()
			if x < 0 {
				x = 0
			}
			if x >= screen.Width() {
				x = screen.Width() - 1
			}
			drawVerticalLine(screen, x)
			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}

func drawVerticalLine(screen image.Selection, x int) {
	for y := 0; y < screen.Height(); y++ {
		screen.SetColor(x, y, colornames.Azure)
	}
}
