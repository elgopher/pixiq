package main

import (
	"fmt"

	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(gl *opengl.OpenGL) {
		win := gl.Open(320, 10, opengl.Title("Press any key..."))
		// Create keyboard instance for window.
		keys := keyboard.New(win)
		loop.Run(win, func(frame *loop.Frame) {
			// Poll keyboard events
			keys.Update()
			// PressedKeys will return all currently pressed keys
			pressedKeys := keys.PressedKeys()
			if len(pressedKeys) > 0 {
				fmt.Println(pressedKeys)
			}
		})
	})
}
