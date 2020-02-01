// Package opengl makes it possible to use Pixiq on PCs with Linux, Windows or MacOS.
// It provides a method for creating OpenGL-accelerated image.Image and Window which
// is an implementation of loop.Screen and keyboard.EventSource.
// Under the hood it is using OpenGL API and GLFW for manipulating windows
// and handling user input.
package opengl

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl/internal"
)

// New creates OpenGL instance.
// MainThreadLoop is needed because some GLFW functions has to be called
// from the main thread.
//
// There is a possibility to create multiple OpenGL objects. Please note though
// that some platforms may limit this number. In integration tests you should
// always remember to destroy the object after test by executing Destroy method,
// because eventually the number of objects may reach the mentioned limit.
//
// New may return error for different reasons, such as OpenGL is not supported
// on the platform.
//
// New will panic if mainThreadLoop is nil.
func New(mainThreadLoop *MainThreadLoop) (*OpenGL, error) {
	if mainThreadLoop == nil {
		panic("nil MainThreadLoop")
	}
	var (
		mainWindow *glfw.Window
		err        error
	)
	mainThreadLoop.Execute(func() {
		err = glfw.Init()
		if err != nil {
			return
		}
		mainWindow, err = createWindow(mainThreadLoop, nil)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	runInOpenGLContextThread := func(job func()) {
		mainThreadLoop.Execute(func() {
			mainThreadLoop.bind(mainWindow)
			job()
		})
	}
	openGL := &OpenGL{
		mainThreadLoop:           mainThreadLoop,
		runInOpenGLContextThread: runInOpenGLContextThread,
		stopPollingEvents:        make(chan struct{}),
		mainWindow:               mainWindow,
	}
	go openGL.startPollingEvents(openGL.stopPollingEvents)
	return openGL, nil
}

// RunOrDie is a shorthand method for starting MainThreadLoop and creating
// OpenGL instance. It runs the given callback function and blocks. It was created
// mainly for educational purposes to save a few keystrokes. In production
// quality code you should write this code yourself and implement a proper error
// handling.
//
// Will panic if OpenGL cannot be created.
func RunOrDie(main func(gl *OpenGL)) {
	StartMainThreadLoop(func(mainThreadLoop *MainThreadLoop) {
		openGL, err := New(mainThreadLoop)
		if err != nil {
			panic(err)
		}
		defer openGL.Destroy()
		main(openGL)
	})
}

func createWindow(mainThreadLoop *MainThreadLoop, share *glfw.Window) (*glfw.Window, error) {
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.False)
	// FIXME: For some reason XVFB does not change the frame buffer size after
	// resizing the window to higher values. That's why the window created
	// here has size equal to the biggest window used in integration tests
	// See: TestWindow_Draw() in opengl_test.go
	win, err := glfw.CreateWindow(3, 3, "OpenGL Pixiq Window", nil, share)
	if err != nil {
		return nil, err
	}
	mainThreadLoop.bind(win)
	if err := gl.Init(); err != nil {
		return nil, err
	}
	return win, nil
}

// OpenGL provides method for creating OpenGL-accelerated image.Image and opening
// windows.
type OpenGL struct {
	mainThreadLoop           *MainThreadLoop
	runInOpenGLContextThread func(func())
	stopPollingEvents        chan struct{}
	mainWindow               *glfw.Window
}

// Destroy cleans all the OpenGL resources associated with this instance.
// This method has to be called in integration tests to clean resources after
// each test. Otherwise on some platforms you may reach the limit of active
// OpenGL contexts.
func (g *OpenGL) Destroy() {
	g.stopPollingEvents <- struct{}{}
	g.runInOpenGLContextThread(func() {
		g.mainWindow.Destroy()
	})
}

func (g *OpenGL) startPollingEvents(stop <-chan struct{}) {
	// fixme: make it configurable
	ticker := time.NewTicker(4 * time.Millisecond) // 250Hz
	for {
		<-ticker.C
		select {
		case <-stop:
			return
		default:
			g.mainThreadLoop.Execute(glfw.PollEvents)
		}
	}
}

// NewImage creates an *image.Image which is using OpenGL acceleration
// under-the-hood.
//
// Example:
//
//	   gl := opengl.New(loop)
//	   defer gl.Destroy()
//	   img, err := gl.NewImage(2, 2)
//
// To avoid coupling with opengl you should define your own factory function
// for creating images and use it instead of directly accessing opengl.OpenGL:
//
//	   type NewImage func(width, height) (*image.Image, error)
//
// Will return error if width or height are negative or image of these dimensions
// cannot be created on a video card. (For instance when dimensions are not
// a power of two)
func (g *OpenGL) NewImage(width, height int) (*image.Image, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	acceleratedImage, err := g.NewTexture(width, height)
	if err != nil {
		return nil, err
	}
	return image.New(width, height, acceleratedImage)
}

// NewTexture returns an OpenGL-accelerated implementation of image.AcceleratedImage
// Will return error if width or height are negative or image of these dimensions
// cannot be created on a video card. (For instance when dimensions are not
// a power of two)
func (g *OpenGL) NewTexture(width, height int) (*Texture, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	var id uint32
	var err error
	g.runInOpenGLContextThread(func() {
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
		if glError := gl.GetError(); glError != gl.NO_ERROR {
			err = fmt.Errorf("OpenGL texture creation failed: %d", glError)
			return
		}
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	})
	if err != nil {
		return nil, err
	}
	return &Texture{
		id:                       id,
		width:                    width,
		height:                   height,
		runInOpenGLContextThread: g.runInOpenGLContextThread,
	}, nil
}

type Texture struct {
	id                       uint32
	width, height            int
	runInOpenGLContextThread func(func())
}

func (t *Texture) TextureID() uint32 {
	return t.id
}

func (t *Texture) Upload(pixels []image.Color) {
	t.runInOpenGLContextThread(func() {
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
}

func (t *Texture) Download(output []image.Color) {
	t.runInOpenGLContextThread(func() {
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

// OpenWindow creates and shows Window.
func (g *OpenGL) OpenWindow(width, height int, options ...WindowOption) (*Window, error) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	// FIXME: EventBuffer size should be configurable
	keyboardEvents := internal.NewKeyboardEvents(keyboard.NewEventBuffer(32))
	screenTexture, err := g.NewTexture(width, height)
	if err != nil {
		return nil, err
	}
	screenImage, err := image.New(width, height, screenTexture)
	if err != nil {
		return nil, err
	}
	win := &Window{
		mainThreadLoop:  g.mainThreadLoop,
		keyboardEvents:  keyboardEvents,
		requestedWidth:  width,
		requestedHeight: height,
		screenTexture:   screenTexture,
		screenImage:     screenImage,
		zoom:            1,
	}
	g.mainThreadLoop.Execute(func() {
		win.glfwWindow, err = createWindow(g.mainThreadLoop, g.mainWindow)
		if err != nil {
			return
		}
		win.glfwWindow.SetKeyCallback(win.keyboardEvents.OnKeyCallback)
		win.program, err = compileProgram(vertexShaderSrc, fragmentShaderSrc)
		if err != nil {
			return
		}
		win.screenPolygon = newScreenPolygon(
			win.program.vertexPositionLocation,
			win.program.texturePositionLocation)
		for _, option := range options {
			if option == nil {
				log.Println("nil option given when opening the window")
				continue
			}
			option(win)
		}
		win.glfwWindow.SetSize(win.requestedWidth*win.zoom, win.requestedHeight*win.zoom)
		win.glfwWindow.Show()
	})
	if err != nil {
		return nil, err
	}
	return win, nil
}

// CompileFragmentShader compiles fragment shader source code written in GLSL.
func (g *OpenGL) CompileFragmentShader(sourceCode string) (*FragmentShader, error) {
	var shader *shader
	var err error
	g.runInOpenGLContextThread(func() {
		shader, err = compileFragmentShader(sourceCode)
	})
	if err != nil {
		return nil, err
	}
	return &FragmentShader{shader: shader}, nil
}

// CompileVertexShader compiles vertex shader source code written in GLSL.
func (g *OpenGL) CompileVertexShader(sourceCode string) (*VertexShader, error) {
	var shader *shader
	var err error
	g.runInOpenGLContextThread(func() {
		shader, err = compileVertexShader(sourceCode)
	})
	if err != nil {
		return nil, err
	}
	return &VertexShader{shader: shader}, err
}

// LinkProgram links an OpenGL program from shaders. Created program can be used
// in image.Modify
func (g *OpenGL) LinkProgram(vertexShader *VertexShader, fragmentShader *FragmentShader) (*Program, error) {
	if vertexShader == nil {
		return nil, errors.New("nil vertexShader")
	}
	if fragmentShader == nil {
		return nil, errors.New("nil fragmentShader")
	}
	var program *program
	var err error
	g.runInOpenGLContextThread(func() {
		program, err = linkProgram(vertexShader.shader, fragmentShader.shader)
	})
	if err != nil {
		return nil, err
	}
	return &Program{program: program}, err
}

// FragmentShader is a part of an OpenGL program which transforms each fragment
// (pixel) color into another one
type FragmentShader struct {
	*shader
}

// VertexShader is a part of an OpenGL program which applies transformations
// to drawn vertices.
type VertexShader struct {
	*shader
}

// Program is shaders linked together
type Program struct {
	*program
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

// Zoom makes window/pixels bigger zoom times.
func Zoom(zoom int) WindowOption {
	return func(window *Window) {
		if zoom > 0 {
			window.zoom = zoom
		}
	}
}
