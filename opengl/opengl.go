package opengl

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq"
)

// New creates OpenGL instance. MainThreadLoop is needed because some GLFW functions has to be called from main thread.
func New(mainThreadLoop *MainThreadLoop) *OpenGL {
	mainThreadLoop.Execute(func() {
		err := glfw.Init()
		if err != nil {
			panic(err)
		}
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.Visible, glfw.False)
		window, err := glfw.CreateWindow(640, 360, "dummy window needed for making the GL context", nil, nil)
		if err != nil {
			panic(err)
		}
		defer window.ShouldClose() // TODO
		window.MakeContextCurrent()
	})
	return &OpenGL{textures: &textures{}, glfwWindows: &glfwWindows{}}
}

// OpenGL provides opengl implementations of pixiq.AcceleratedImages and pixiq.SystemWindows
type OpenGL struct {
	textures    *textures
	glfwWindows *glfwWindows
}

// AcceleratedImages returns opengl implementation of pixiq.AcceleratedImages
func (g OpenGL) AcceleratedImages() pixiq.AcceleratedImages {
	return g.textures
}

// SystemWindows returns opengl implementation of pixiq.SystemWindows
func (g OpenGL) SystemWindows() pixiq.SystemWindows {
	return g.glfwWindows
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
