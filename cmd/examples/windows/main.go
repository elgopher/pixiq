package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

var (
	red  = pixiq.RGBA(255, 0, 0, 255)
	blue = pixiq.RGBA(0, 0, 255, 255)
)

func main() {
	opengl.StartMainThreadLoop(func(loop *opengl.MainThreadLoop) {
		openGL := opengl.New(loop)
		windows := openGL.Windows()
		window1 := windows.Open(320, 180)
		window2 := windows.Open(160, 90)
		screens := pixiq.NewScreens(pixiq.NewImages(openGL.AcceleratedImages()))
		go screens.Loop(window1, fillWith(red))
		screens.Loop(window2, fillWith(blue))
	})
}

func fillWith(color pixiq.Color) func(frame *pixiq.Frame) {
	return func(frame *pixiq.Frame) {
		screen := frame.Screen()
		for y := 0; y < screen.Height(); y++ {
			for x := 0; x < screen.Width(); x++ {
				screen.SetColor(x, y, color)
			}
		}
	}
}
