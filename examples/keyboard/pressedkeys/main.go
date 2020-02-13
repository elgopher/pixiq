package main

import (
	"fmt"
	"log"

	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		win, err := gl.OpenWindow(320, 10, opengl.Title("Press any key..."))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
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
