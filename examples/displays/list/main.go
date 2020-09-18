package main

import (
	"fmt"
	"github.com/jacekolszak/pixiq/glfw"
)

func main() {
	// This example shows how to list all displays
	glfw.StartMainThreadLoop(func(mainThreadLoop *glfw.MainThreadLoop) {
		// Displays instance requires mainThreadLoop because accessing information
		// about displays must be done from the main thread.
		displays, err := glfw.Displays(mainThreadLoop)
		if err != nil {
			panic(err)
		}

		all := displays.All()
		for _, display := range all {
			fmt.Println("Name:", display.Name())
			physicalSize := display.PhysicalSize()
			fmt.Printf("Phyical size: %d mm x %d mm\n", physicalSize.Width(), physicalSize.Height())
			videoMode := display.VideoMode()
			fmt.Printf("Current resolution: %d x %d, %d Hz\n", videoMode.Width(), videoMode.Height(), videoMode.RefreshRate())
			fmt.Println()
		}
	})
}
