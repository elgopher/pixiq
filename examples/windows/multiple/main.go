package main

import (
	"log"
	"sync"

	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
)

// This example shows how to open two windows at the same time.
//
// Please note that this functionality is experimental and may change in the
// near future. Such feature may be harmful for overall performance of Pixiq.
func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		redWindow, err := openGL.OpenWindow(320, 180, glfw.Title("red"))
		if err != nil {
			log.Panicf("red OpenWindow failed: %v", err)
		}
		blueWindow, err := openGL.OpenWindow(250, 90, glfw.Title("blue"))
		if err != nil {
			log.Panicf("blue OpenWindow failed: %v", err)
		}

		var waitUntilAllClosed sync.WaitGroup
		waitUntilAllClosed.Add(2)

		// Start the loop in the background, because loop.Run blocks
		// the current goroutine.
		go func() {
			loop.Run(redWindow, func(frame *loop.Frame) {
				fillScreenWith(frame.Screen(), colornames.Red)
				if redWindow.ShouldClose() {
					frame.StopLoopEventually()
				}
			})
			// clean resources
			redWindow.Close()
			waitUntilAllClosed.Done()
		}()

		// Start another one.
		go func() {
			loop.Run(blueWindow, func(frame *loop.Frame) {
				fillScreenWith(frame.Screen(), colornames.Blue)
				if blueWindow.ShouldClose() {
					frame.StopLoopEventually()
				}
			})
			blueWindow.Close()
			waitUntilAllClosed.Done()
		}()

		// wait for all windows to be closed
		waitUntilAllClosed.Wait()
	})
}

// fillScreenWith returns a function filling whole Screen with specific color.
func fillScreenWith(screen image.Selection, color image.Color) {
	for y := 0; y < screen.Height(); y++ {
		for x := 0; x < screen.Width(); x++ {
			screen.SetColor(x, y, color)
		}
	}
}
