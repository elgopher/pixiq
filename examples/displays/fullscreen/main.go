package main

import (
	"fmt"
	"github.com/jacekolszak/pixiq/clear"
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
)

func main() {
	// This example shows how to open a window in a fullscreen mode
	glfw.StartMainThreadLoop(func(mainThreadLoop *glfw.MainThreadLoop) {
		// Displays instance requires mainThreadLoop because accessing information
		// about displays must be done from the main thread.
		displays, err := glfw.Displays(mainThreadLoop)
		if err != nil {
			panic(err)
		}
		// Get Primary display. This is usually the display where elements like the Windows task bar
		// or the OS X menu bar is located.
		primary, ok := displays.Primary()
		if !ok {
			panic("no displays found")
		}
		// get current video mode which is usually the best one to pick
		videoMode := primary.VideoMode()
		// try to find the window size and zoom close enough to requested 640x360
		width, height, zoom := adjustSize(videoMode, 640, 360)
		fmt.Printf("Adjusted size is %d x %d, zoom=%d\n", width, height, zoom)

		gl, err := glfw.NewOpenGL(mainThreadLoop)
		if err != nil {
			panic(err)
		}

		// use glfw.FullScreen option to open window in full-screen
		win, err := gl.OpenWindow(width, height, glfw.Zoom(zoom), glfw.FullScreen(videoMode))
		if err != nil {
			panic(err)
		}

		prepareScreen(win)

		// Show full screen for 3 seconds
		fmt.Println("Refresh rate is", videoMode.RefreshRate())
		for x := 0; x < videoMode.RefreshRate()*3; x++ {
			win.Draw() // blocks until VSync
		}
		win.Close()
	})
}

// Adjusts the size of window based on the VideoMode. It first try to to increase the zoom, then will adjust
// the width and height if display has different ratio.
func adjustSize(mode glfw.VideoMode, width, height int) (newWidth, newHeight, zoom int) {
	// TODO This functionality should be in a new package
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
	horizontalMargin := (mode.Width() - (width * zoom)) / zoom
	verticalMargin := (mode.Height() - (height * zoom)) / zoom
	return width + horizontalMargin, height + verticalMargin, zoom
}

func prepareScreen(win *glfw.Window) {
	screen := win.Screen()
	clearTool := clear.New()
	clearTool.SetColor(colornames.Lightgray)
	clearTool.Clear(screen)
	screen.SetColor(100, 100, colornames.Black)
}
