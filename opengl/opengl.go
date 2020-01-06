// Package opengl makes it possible to use Pixiq on PCs with Linux, Windows or MacOS.
// It provides implementation of both pixiq.AcceleratedImages and pixiq.Screen.
// Under the hood it is using OpenGL API and GLFW for manipulating windows
// and handling user input.
package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq"
)

// New creates OpenGL instance.
// MainThreadLoop is needed because some GLFW functions has to be called
// from the main thread.
func New(loop *MainThreadLoop) *OpenGL {
	var mainWindow *glfw.Window
	loop.Execute(func() {
		err := glfw.Init()
		if err != nil {
			panic(err)
		}
		mainWindow = createWindow(nil)
	})
	return &OpenGL{
		textures: &textures{mainThreadLoop: loop},
		windows: &Windows{
			mainWindow:     mainWindow,
			mainThreadLoop: loop,
		},
	}
}

// Run is a shorthand method for creating Pixiq objects with OpenGL acceleration
// and Windows. It runs the given callback function and blocks. It was created
// mainly for educational purposes to save a few keystrokes.
func Run(main func(gl *OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops)) {
	StartMainThreadLoop(func(loop *MainThreadLoop) {
		openGL := New(loop)
		images := pixiq.NewImages(openGL.AcceleratedImages())
		loops := pixiq.NewScreenLoops(images)
		main(openGL, images, loops)
	})
}

func createWindow(share *glfw.Window) *glfw.Window {
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.False)
	win, err := glfw.CreateWindow(1, 1, "OpenGL Pixiq Window", nil, share)
	if err != nil {
		panic(err)
	}
	win.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}
	return win
}

// OpenGL provides opengl implementations of pixiq.AcceleratedImages
// and pixiq.Screen. It provides Windows object for opening system windows.
type OpenGL struct {
	textures *textures
	windows  *Windows
}

// AcceleratedImages returns opengl implementation of pixiq.AcceleratedImages.
func (g OpenGL) AcceleratedImages() pixiq.AcceleratedImages {
	return g.textures
}

// Windows returns object for opening system windows. Each open Window
// is a pixiq.Screen implementation.
func (g OpenGL) Windows() *Windows {
	return g.windows
}

// Windows is used for opening system windows.
type Windows struct {
	mainThreadLoop *MainThreadLoop
	// mainWindow contains textures and user programs
	mainWindow *glfw.Window
}

// Open creates and shows Window.
func (g Windows) Open(width, height int, options ...WindowOption) *Window {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	var err error
	win := &Window{mainThreadLoop: g.mainThreadLoop}
	g.mainThreadLoop.Execute(func() {
		win.glfwWindow = createWindow(g.mainWindow)
		win.program, err = compileProgram()
		if err != nil {
			panic(err)
		}
		win.screenPolygon = newScreenPolygon(
			win.program.vertexPositionLocation,
			win.program.texturePositionLocation)
		win.glfwWindow.SetSize(width, height)
		for _, option := range options {
			option(win)
		}
		win.glfwWindow.Show()
	})
	return win
}

// WindowOption is an option used when opening the window
type WindowOption func(window *Window)

// NoDecorationHint is Window hint hiding the border, close widget, etc.
// Exact behaviour depends on the platform.
func NoDecorationHint() WindowOption {
	return func(win *Window) {
		win.glfwWindow.SetAttrib(glfw.Decorated, glfw.False)
	}
}

// Title sets the window title
func Title(title string) WindowOption {
	return func(window *Window) {
		window.glfwWindow.SetTitle(title)
	}
}

// Window is an implementation of pixiq.Screen
type Window struct {
	glfwWindow     *glfw.Window
	program        *program
	mainThreadLoop *MainThreadLoop
	screenPolygon  *screenPolygon
}

// Draw draws image spanning the whole window to the invisible buffer.
func (g *Window) Draw(image *pixiq.Image) {
	texture, isGL := image.Upload().(GLTexture)
	if !isGL {
		panic("opengl Window can only draw images accelerated with opengl.GLTexture")
	}
	g.mainThreadLoop.Execute(func() {
		g.glfwWindow.MakeContextCurrent()
		g.program.use()
		w, h := g.glfwWindow.GetFramebufferSize()
		gl.Viewport(0, 0, int32(w), int32(h))
		gl.BindTexture(gl.TEXTURE_2D, texture.TextureID())
		g.screenPolygon.draw()
	})
}

// SwapImages makes last drawn image visible
func (g *Window) SwapImages() {
	g.mainThreadLoop.Execute(func() {
		g.glfwWindow.SwapBuffers()
		glfw.PollEvents()
	})
}

// Close closes the window and cleans resources
func (g *Window) Close() {
	g.mainThreadLoop.Execute(func() {
		g.glfwWindow.Destroy()
	})
}

// ShouldClose reports the value of the close flag of the window.
// The flag is set to true when user clicks Close button or hits ALT+F4/CMD+Q.
func (g *Window) ShouldClose() bool {
	var shouldClose bool
	g.mainThreadLoop.Execute(func() {
		shouldClose = g.glfwWindow.ShouldClose()
	})
	return shouldClose
}

// Width returns the width of the window in pixels.
// If zooming is used the width is not multiplied by zoom.
func (g *Window) Width() int {
	var width int
	g.mainThreadLoop.Execute(func() {
		width, _ = g.glfwWindow.GetSize()
	})
	return width
}

// Height returns the height of the window in pixels.
// If zooming is used the height is not multiplied by zoom.
func (g *Window) Height() int {
	var height int
	g.mainThreadLoop.Execute(func() {
		_, height = g.glfwWindow.GetSize()
	})
	return height
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
	return &texture{
		id:             id,
		width:          width,
		height:         height,
		mainThreadLoop: g.mainThreadLoop,
	}
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
