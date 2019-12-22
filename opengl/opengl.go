package opengl

import "github.com/jacekolszak/pixiq"

// Run should be executed in main thread.
func Run(runInDifferentGoroutine func(images *pixiq.Images)) {
	images := pixiq.NewImages(&glTextures{})
	runInDifferentGoroutine(images)
}

type glTextures struct {
}

func (g *glTextures) New(width, height int) pixiq.AcceleratedImage {
	panic("implement me")
}
