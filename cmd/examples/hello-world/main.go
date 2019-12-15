package main

import (
	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.Run(func(images *pixiq.Images) {
		images.New(16, 16)
	})
}
