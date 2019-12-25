package opengl

import "github.com/jacekolszak/pixiq"

// Run should be executed in main thread.
func Run(runInDifferentGoroutine func(images *pixiq.Images)) {
	images := pixiq.NewImages(&textures{})
	runInDifferentGoroutine(images)
}

type textures struct {
}

func (g *textures) New(width, height int) pixiq.AcceleratedImage {
	panic("implement me")
}

type texture struct {
}

func (t *texture) Upload(pixels []pixiq.Color) {
	panic("implement me")
}
