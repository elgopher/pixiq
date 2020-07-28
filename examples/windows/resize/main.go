package main

import (
	"github.com/jacekolszak/pixiq/clear"
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/mouse"
	"log"
)

func main() {
	// This example shows how to open a window in a fullscreen mode
	glfw.StartMainThreadLoop(func(mainThreadLoop *glfw.MainThreadLoop) {
		gl, err := glfw.NewOpenGL(mainThreadLoop)
		if err != nil {
			panic(err)
		}

		width := 121
		height := 101
		zoom := 5
		win, err := gl.OpenWindow(width, height, glfw.Zoom(zoom), glfw.Title("Scroll the mouse wheel to zoom in/out"), glfw.Resizable(true))
		if err != nil {
			panic(err)
		}

		mouseState := mouse.New(win)
		for {
			screen := win.Screen()
			clearTool := clear.New()
			clearTool.SetColor(colornames.Lightgray)
			clearTool.Clear(screen)
			screen.SetColor(screen.Width()/2, screen.Height()/2, colornames.Black)

			mouseState.Update()
			if mouseState.Scroll().Y() > 0 {
				zoom++
				if err := win.Resize(width, height, zoom); err != nil {
					log.Printf("Resize failed: %v", err)
				}
			}
			if mouseState.Scroll().Y() < 0 && zoom > 5 {
				zoom--
				if err := win.Resize(width, height, zoom); err != nil {
					log.Printf("Resize failed: %v", err)
				}
			}
			win.Draw()
			if win.ShouldClose() {
				break
			}
		}
	})
}
