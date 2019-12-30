package opengl

import "github.com/jacekolszak/pixiq"

// Run should be executed in main goroutine.
// It takes control over executing goroutine until runInDifferentGoroutine finishes.
func Run(runInDifferentGoroutine func(gl *OpenGL)) {
	runInDifferentGoroutine(&OpenGL{})
}

// OpenGL provides opengl implementations of pixiq.AcceleratedImages and pixiq.SystemWindows
type OpenGL struct {
}

// AcceleratedImages returns opengl implementation of pixiq.AcceleratedImages
func (g OpenGL) AcceleratedImages() pixiq.AcceleratedImages {
	return &textures{}
}

// SystemWindows returns opengl implementation of pixiq.SystemWindows
func (g OpenGL) SystemWindows() pixiq.SystemWindows {
	return glfwWindows{}
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
