// Package opengl makes it possible to use Pixiq on PCs with Linux, Windows or MacOS.
// It provides implementation of both pixiq.AcceleratedImages and pixiq.Screen.
// Under the hood it is using OpenGL API and GLFW for manipulating windows
// and handling user input.
package opengl

import (
	"errors"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl/internal"
)

// New creates OpenGL instance.
// MainThreadLoop is needed because some GLFW functions has to be called
// from the main thread.
func New(loop *MainThreadLoop) *OpenGL {
	if loop == nil {
		panic("nil MainThreadLoop")
	}
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
func (w Windows) Open(width, height int, options ...WindowOption) *Window {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	var err error
	win := &Window{
		mainThreadLoop:  w.mainThreadLoop,
		keyboardEvents:  internal.NewKeyboardEvents(16),
		requestedWidth:  width,
		requestedHeight: height,
		zoom:            1,
	}
	w.mainThreadLoop.Execute(func() {
		win.glfwWindow = createWindow(w.mainWindow)
		win.glfwWindow.SetKeyCallback(win.keyboardEvents.OnKeyCallback)
		win.program, err = compileProgram()
		if err != nil {
			return
		}
		win.screenPolygon = newScreenPolygon(
			win.program.vertexPositionLocation,
			win.program.texturePositionLocation)
		for _, option := range options {
			if option == nil {
				err = errors.New("nil option")
				return
			}
			option(win)
		}
		win.glfwWindow.SetSize(win.requestedWidth, win.requestedHeight)
		win.glfwWindow.Show()
	})
	if err != nil {
		panic(err)
	}
	return win
}

// WindowOption is an option used when opening the window.
type WindowOption func(window *Window)

// NoDecorationHint is Window hint hiding the border, close widget, etc.
// Exact behaviour depends on the platform.
func NoDecorationHint() WindowOption {
	return func(win *Window) {
		win.glfwWindow.SetAttrib(glfw.Decorated, glfw.False)
	}
}

// Title sets the window title.
func Title(title string) WindowOption {
	return func(window *Window) {
		window.glfwWindow.SetTitle(title)
	}
}

func Zoom(zoom int) WindowOption {
	return func(window *Window) {
		if zoom > 0 {
			window.zoom = zoom
			window.requestedWidth = window.requestedWidth * zoom
			window.requestedHeight = window.requestedHeight * zoom
		}
	}
}

// Window is an implementation of pixiq.Screen.
type Window struct {
	glfwWindow      *glfw.Window
	program         *program
	mainThreadLoop  *MainThreadLoop
	screenPolygon   *screenPolygon
	keyboardEvents  *internal.KeyboardEvents
	requestedWidth  int
	requestedHeight int
	zoom            int
}

// Draw draws image spanning the whole window to the invisible buffer.
func (w *Window) Draw(image *pixiq.Image) {
	texture, isGL := image.Upload().(GLTexture)
	if !isGL {
		panic("opengl Window can only draw images accelerated with opengl.GLTexture")
	}
	w.mainThreadLoop.Execute(func() {
		w.glfwWindow.MakeContextCurrent()
		w.program.use()
		width, height := w.glfwWindow.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.BindTexture(gl.TEXTURE_2D, texture.TextureID())
		w.screenPolygon.draw()
	})
}

// SwapImages makes last drawn image visible.
func (w *Window) SwapImages() {
	w.mainThreadLoop.Execute(func() {
		w.glfwWindow.SwapBuffers()
	})
}

// Close closes the window and cleans resources.
func (w *Window) Close() {
	w.mainThreadLoop.Execute(func() {
		w.glfwWindow.Destroy()
	})
}

// ShouldClose reports the value of the close flag of the window.
// The flag is set to true when user clicks Close button or hits ALT+F4/CMD+Q.
func (w *Window) ShouldClose() bool {
	var shouldClose bool
	w.mainThreadLoop.Execute(func() {
		shouldClose = w.glfwWindow.ShouldClose()
	})
	return shouldClose
}

// Width returns the actual width of the window in pixels. It may be different
// than requested width used when window was open due to platform limitation.
// If zooming is used the width is not multiplied by zoom.
func (w *Window) Width() int {
	var width int
	w.mainThreadLoop.Execute(func() {
		width, _ = w.glfwWindow.GetSize()
	})
	return width / w.zoom
}

// Height returns the actual height of the window in pixels. It may be different
// than requested height used when window was open due to platform limitation.
// If zooming is used the height is not multiplied by zoom.
func (w *Window) Height() int {
	var height int
	w.mainThreadLoop.Execute(func() {
		_, height = w.glfwWindow.GetSize()
	})
	return height / w.zoom
}

// Poll retrieves and removes next keyboard Event. If there are no more
// events false is returned. It implements keyboard.EventSource method.
func (w *Window) Poll() (keyboard.Event, bool) {
	if w.keyboardEvents.Drained() {
		w.mainThreadLoop.Execute(func() {
			glfw.PollEvents()
		})
	}
	return w.keyboardEvents.Poll()
}

type textures struct {
	mainThreadLoop *MainThreadLoop
}

func (t *textures) New(width, height int) pixiq.AcceleratedImage {
	var id uint32
	t.mainThreadLoop.Execute(func() {
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
		mainThreadLoop: t.mainThreadLoop,
	}
}

// GLTexture is an OpenGL texture which can be sampled to create rectangles on screen.
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
