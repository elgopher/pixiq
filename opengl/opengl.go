package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq"
)

// New creates OpenGL instance providing implementation of AcceleratedImages and Windows for opening system windows.
// Under the hood it is using OpenGL API and GLFW for manipulating windows and handling user input.
// MainThreadLoop is needed because some GLFW functions has to be called from main thread.
func New(loop *MainThreadLoop) *OpenGL {
	var program *program
	var window *glfw.Window
	loop.Execute(func() {
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
		window, err = glfw.CreateWindow(1, 1, "OpenGL Pixiq Window", nil, nil)
		if err != nil {
			panic(err)
		}
		window.MakeContextCurrent()
		if err := gl.Init(); err != nil {
			panic(err)
		}
		program, err = compileProgram()
		if err != nil {
			panic(err)
		}
	})
	return &OpenGL{
		textures: &textures{mainThreadLoop: loop},
		windows: &Windows{
			window:         window,
			program:        program,
			mainThreadLoop: loop,
		},
	}
}

// OpenGL provides opengl implementations of AcceleratedImages. It also provides Windows for opening system windows
type OpenGL struct {
	textures *textures
	windows  *Windows
}

// AcceleratedImages returns opengl implementation of AcceleratedImages
func (g OpenGL) AcceleratedImages() pixiq.AcceleratedImages {
	return g.textures
}

// Windows returns object for opening system windows
func (g OpenGL) Windows() *Windows {
	return g.windows
}

// Windows is for opening system windows
type Windows struct {
	program        *program
	mainThreadLoop *MainThreadLoop
	window         *glfw.Window
}

// Open creates and show window
func (g Windows) Open(width, height int, hints ...WindowHint) pixiq.Window {
	frameImagePolygon := newFrameImagePolygon(g.mainThreadLoop, g.program.vertexPositionLocation, g.program.texturePositionLocation)
	g.mainThreadLoop.Execute(func() {
		for _, hint := range hints {
			hint.apply(g.window)
		}
		g.window.SetSize(width, height)
		g.window.Show()
	})
	return &window{window: g.window, program: g.program, mainThreadLoop: g.mainThreadLoop, frameImagePolygon: frameImagePolygon}
}

// WindowHint is a hint which may (or may not) be applied to window (depending on operating system and other factors).
type WindowHint interface {
	apply(window *glfw.Window)
}

// NoDecorated is window hint hiding the window's title bar
type NoDecorated struct{}

func (NoDecorated) apply(window *glfw.Window) {
	window.SetAttrib(glfw.Decorated, glfw.False)
}

type window struct {
	window            *glfw.Window
	program           *program
	mainThreadLoop    *MainThreadLoop
	frameImagePolygon *frameImagePolygon
}

func (g *window) Draw(image *pixiq.Image) {
	texture, isGL := image.Upload().(GLTexture)
	if !isGL {
		panic("opengl Windows implementation can only draw images accelerated with opengl.GLTexture")
	}
	g.mainThreadLoop.Execute(func() {
		g.program.use()
		w, h := g.window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(w), int32(h))
		gl.BindTexture(gl.TEXTURE_2D, texture.TextureID())
		g.frameImagePolygon.draw()
	})
}

func (g *window) SwapImages() {
	g.mainThreadLoop.Execute(func() {
		g.window.SwapBuffers()
		glfw.PollEvents()
	})
}

func (g *window) Close() {
	g.mainThreadLoop.Execute(func() {
		g.window.Destroy()
	})
}

func (g *window) Width() int {
	w, _ := g.window.GetSize()
	return w
}

func (g *window) Height() int {
	_, h := g.window.GetSize()
	return h
}

type textures struct {
	mainThreadLoop *MainThreadLoop
}

func (g *textures) New(width, height int) pixiq.AcceleratedImage {
	var id uint32
	g.mainThreadLoop.Execute(func() {
		gl.GenTextures(1, &id)
		gl.BindTexture(gl.TEXTURE_2D, id)
		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(width),
			int32(height),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(nil),
		)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	})
	return &texture{id: id, width: width, height: height, mainThreadLoop: g.mainThreadLoop}
}

// GLTexture is an OpenGL texture which can be sampled to create rectangles on screen
type GLTexture interface {
	pixiq.AcceleratedImage
	TextureID() uint32
}

type texture struct {
	pixels         []pixiq.Color
	id             uint32
	width, height  int
	mainThreadLoop *MainThreadLoop
}

func (t *texture) TextureID() uint32 {
	return t.id
}

func (t *texture) Upload(pixels []pixiq.Color) {
	t.mainThreadLoop.Execute(func() {
		gl.BindTexture(gl.TEXTURE_2D, t.id)
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			int32(0),
			int32(0),
			int32(t.width),
			int32(t.height),
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(pixels),
		)
	})
	t.pixels = pixels
}
func (t *texture) Download(output []pixiq.Color) {
	t.mainThreadLoop.Execute(func() {
		gl.BindTexture(gl.TEXTURE_2D, t.id)
		gl.GetTexImage(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(output),
		)
	})
}
