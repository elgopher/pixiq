package opengl

import "github.com/jacekolszak/pixiq"

// Run should be executed in main thread.
func Run(runInDifferentGoroutine func(images *pixiq.Images)) {
	images := pixiq.NewImages()
	runInDifferentGoroutine(images)
}
