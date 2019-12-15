package opengl

import "github.com/jacekolszak/pixiq"

func Run(runInDifferentGoroutine func(images *pixiq.Images)) {
	images := pixiq.NewImages()
	runInDifferentGoroutine(images)
}
