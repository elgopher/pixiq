package main

import (
	"fmt"
	"log"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/loop"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		win, err := openGL.OpenWindow(320, 10, glfw.Title("Press any key..."))
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
