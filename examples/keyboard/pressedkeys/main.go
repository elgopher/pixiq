package main

import (
	"fmt"
	"log"

	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/keyboard"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		win, err := openGL.OpenWindow(320, 10, glfw.Title("Press any key..."))
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		// Create keyboard instance for window.
		keys := keyboard.New(win)
		for {
			// Poll keyboard events
			keys.Update()
			// PressedKeys will return all currently pressed keys
			pressedKeys := keys.PressedKeys()
			if len(pressedKeys) > 0 {
				fmt.Println(pressedKeys)
			}

			if win.ShouldClose() {
				break
			}
		}
	})
}
