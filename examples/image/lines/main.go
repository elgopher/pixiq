package main

import (
	"log"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
)

// This example shows how to set all pixels using Lines
func main() {
	glfw.RunOrDie(func(gl *glfw.OpenGL) {
		window, err := gl.OpenWindow(640, 360)
		if err != nil {
			log.Panicf("OpenWindow failed: %v", err)
		}
		for {
			screen := window.Screen()
			lines := screen.Lines()
			for y := 0; y < lines.Length(); y++ {
				line := lines.LineForWrite(y)
				for x := 0; x < len(line); x++ {
					color := image.RGBA(byte(x%255), byte(y%255), 255, 255)
					line[x] = color
				}
			}

			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}
