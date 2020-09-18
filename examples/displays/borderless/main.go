package main

import (
	"fmt"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glclear"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/mouse"
)

func main() {
	// This example shows how to switch from windowed to full screen mode
	glfw.StartMainThreadLoop(func(mainThreadLoop *glfw.MainThreadLoop) {
		// Displays instance requires mainThreadLoop because accessing information
		// about displays must be done from the main thread.
		displays, err := glfw.Displays(mainThreadLoop)
		if err != nil {
			panic(err)
		}

		gl, err := glfw.NewOpenGL(mainThreadLoop)
		if err != nil {
			panic(err)
		}
		// Open standard window
		screenWidth, screenHeight := 640, 360
		win, err := gl.OpenWindow(screenWidth, screenHeight, glfw.Title("Press left mouse button to borderless fullscreen"))
		if err != nil {
			panic(err)
		}

		mouseState := mouse.New(win)

		fullscreen := false

		prepareScreen(gl, win)

		for {
			mouseState.Update()
			if mouseState.JustReleased(mouse.Left) {
				if !fullscreen {
					display := currentDisplay(win, displays)
					mode := display.VideoMode()
					// Disable automatic iconify on focus loss.
					win.SetAutoIconifyHint(false)
					zoom := adjustZoom(mode, screenWidth, screenHeight)
					fmt.Println(zoom)
					// Turn on the full screen
					win.EnterFullScreen(mode, zoom)
					fullscreen = true
				} else {
					win.ExitFullScreen()
					fullscreen = false
				}
			}
			win.Draw()
			if win.ShouldClose() {
				break
			}
		}

	})
}

func prepareScreen(gl *glfw.OpenGL, win *glfw.Window) {
	clearTool := glclear.New(gl.Context())
	clearTool.SetColor(colornames.White)
	screen := win.Screen()
	clearTool.Clear(screen)
	clearTool.SetColor(colornames.Black)
	clearTool.Clear(screen.Selection(270, 130).WithSize(100, 100))
}

func currentDisplay(win *glfw.Window, displays *glfw.DisplaysAPI) glfw.Display {
	// TODO This functionality should be in a new package
	highestArea := 0
	all := displays.All()
	bestDisplay := all[0]
	for _, display := range all {
		workarea := display.Workarea()
		videoMode := display.VideoMode()
		left := max(win.X(), workarea.X())
		right := min(videoMode.Width()+workarea.X(), win.Width()+win.X())
		top := max(win.Y(), workarea.Y())
		bottom := min(videoMode.Height()+workarea.Y(), win.Height()+win.Y())
		w := right - left
		h := bottom - top
		if w > 0 && h > 0 {
			area := w * h
			if area > highestArea {
				bestDisplay = display
				highestArea = area
			}
		}
	}
	return bestDisplay
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func adjustZoom(mode glfw.VideoMode, width, height int) (zoom int) {
	zoom = 1
	w := width
	h := height
	for mode.Width() > w && mode.Height() > h {
		zoom++
		w = width * zoom
		h = height * zoom
	}
	if w > mode.Width() || h > mode.Height() {
		zoom--
	}
	return zoom
}
