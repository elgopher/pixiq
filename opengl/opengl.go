package opengl

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq"
)

// New creates OpenGL instance providing implementations of both AcceleratedImages and SystemWindows.
// MainThreadLoop is needed because some GLFW functions has to be called from main thread.
func New(loop *MainThreadLoop) *OpenGL {
	loop.Execute(func() {
		err := glfw.Init()
		if err != nil {
			panic(err)
		}
		//glfw.WindowHint(glfw.ContextCreationAPI, glfw.OSMesaContextAPI)
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.Visible, glfw.False)
		window, err := glfw.CreateWindow(1, 1, "dummy window needed for making the GL context", nil, nil)
		if err != nil {
			panic(err)
		}
		window.MakeContextCurrent()
	})
	return &OpenGL{textures: &textures{}, glfwWindows: &glfwWindows{}, mainThreadLoop: loop}
}

// OpenGL provides opengl implementations of AcceleratedImages and SystemWindows
type OpenGL struct {
	textures       *textures
	glfwWindows    *glfwWindows
	mainThreadLoop *MainThreadLoop
}

// AcceleratedImages returns opengl implementation of AcceleratedImages
func (g OpenGL) AcceleratedImages() pixiq.AcceleratedImages {
	return g.textures
}

// SystemWindows returns opengl implementation of SystemWindows
func (g OpenGL) SystemWindows() pixiq.SystemWindows {
	return g.glfwWindows
}

// Terminate closes all windows and frees resources
func (g OpenGL) Terminate() {
	g.mainThreadLoop.Execute(func() {
		glfw.Terminate()
	})
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
	pixels []pixiq.Color
}

func (t *texture) TextureID() uint32 {
	panic("implement me")
}

func (t *texture) Upload(pixels []pixiq.Color) {
	t.pixels = pixels
}
func (t *texture) Download(output []pixiq.Color) {
	for i := 0; i < len(output); i++ {
		output[i] = t.pixels[i]
	}
}
