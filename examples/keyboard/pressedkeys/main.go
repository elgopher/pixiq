package main

import (
	"fmt"

	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
		win := gl.Windows().Open(320, 10, opengl.Title("Press any key..."))
		keys := keyboard.New(win)
		loops.Loop(win, func(frame *pixiq.Frame) {
			keys.Update()
			pressedKeys := keys.PressedKeys()
			if len(pressedKeys) > 0 {
				fmt.Println(pressedKeys)
			}
		})
	})
}
