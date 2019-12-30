package opengl

import "github.com/jacekolszak/pixiq"

// Run should be executed in main thread. It passes opengl implementations of AcceleratedImages and SystemWindows.
func Run(runInDifferentGoroutine func(acceleratedImages pixiq.AcceleratedImages, systemWindows pixiq.SystemWindows)) {
	runInDifferentGoroutine(&textures{}, &glfwWindows{})
}

type glfwWindows struct {
}

func (g glfwWindows) Open(width, height int) pixiq.SystemWindow {
	return &glfwWindow{}
}

type glfwWindow struct {
}

func (g *glfwWindow) Draw(image *pixiq.Image) {
	_, isGL := image.Upload().(GLTexture)
	if !isGL {
		panic("opengl SystemWindows implementation can only draw images accelerated with opengl.GLTexture")
	}
}

type textures struct {
}

func (g *textures) New(width, height int) pixiq.AcceleratedImage {
	return &texture{}
}

// GLTexture is an OpenGL texture which can be sampled to create rectangles on screen
type GLTexture interface {
	pixiq.AcceleratedImage
	TextureID() uint32
}

type texture struct {
}

func (t *texture) TextureID() uint32 {
	panic("implement me")
}

func (t *texture) Upload(pixels []pixiq.Color) {
	panic("implement me")
}
