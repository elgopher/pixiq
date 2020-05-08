// Package glfw makes it possible to use Pixiq on PCs with Linux, Windows or MacOS.
// It provides a method for creating OpenGL-accelerated image.Image and Window which
// is an implementation of loop.Screen and keyboard.EventSource.
// Under the hood it is using OpenGL API and GLFW for manipulating windows
// and handling user input.
package glfw

import (
	"log"
	"sync"
	"time"

	gl33 "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/glfw/internal"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/mouse"
)

// NewOpenGL creates OpenGL instance.
// MainThreadLoop is needed because some GLFW functions has to be called
// from the main thread.
//
// There is a possibility to create multiple OpenGL objects. Please note though
// that some platforms may limit this number. In integration tests you should
// always remember to destroy the object after test by executing Destroy method,
// because eventually the number of objects may reach the mentioned limit.
//
// NewOpenGL may return error for different reasons, such as OpenGL is not supported
// on the platform.
//
// NewOpenGL will panic if mainThreadLoop is nil.
func NewOpenGL(mainThreadLoop *MainThreadLoop) (*OpenGL, error) {
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
	})
	if err != nil {
		return nil, err
	}
	runInOpenGLThread := func(job func()) {
		mainThreadLoop.Execute(func() {
			mainThreadLoop.bind(mainWindow)
			job()
		})
	}
	api := &context{
		run: func(f func()) {
			mainThreadLoop.executeCommand(command{
				window:  mainWindow,
				execute: f,
			})
		},
		runAsync: func(f func()) {
			mainThreadLoop.executeAsyncCommand(command{
				window:  mainWindow,
				execute: f,
			})
		},
	}
	openGL := &OpenGL{
		mainThreadLoop:    mainThreadLoop,
		runInOpenGLThread: runInOpenGLThread,
		stopPollingEvents: make(chan struct{}),
		mainWindow:        mainWindow,
		api:               api,
		context:           gl.NewContext(api),
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
		openGL, err := NewOpenGL(mainThreadLoop)
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
	// resizing the window to higher values than initial ones. That's why the window
	// created here has size equal to the biggest window used in integration tests
	// See: TestWindow_Draw() in glfw_test.go
	win, err := glfw.CreateWindow(3, 3, "OpenGL Pixiq Window", nil, share)
	if err != nil {
		return nil, err
	}
	mainThreadLoop.bind(win)
	if err := gl33.Init(); err != nil {
		return nil, err
	}
	return win, nil
}

// OpenGL provides method for creating OpenGL-accelerated image.Image and opening
// windows.
type OpenGL struct {
	mainThreadLoop    *MainThreadLoop
	runInOpenGLThread func(func())
	stopPollingEvents chan struct{}
	mainWindow        *glfw.Window
	api               gl.API
	context           *gl.Context
}

// Destroy cleans all the OpenGL resources associated with this instance.
// This method has to be called in integration tests to clean resources after
// each test. Otherwise on some platforms you may reach the limit of active
// OpenGL contexts.
func (g *OpenGL) Destroy() {
	g.stopPollingEvents <- struct{}{}
	g.runInOpenGLThread(func() {
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
//	   openGL := glfw.NewOpenGL(loop)
//	   defer openGL.Destroy()
//	   img, err := openGL.NewImage(2, 2)
//
// To avoid coupling with glfw you should define your own factory function
// for creating images and use it instead of directly accessing glfw.OpenGL:
//
//	   type NewImage func(width, height int) *image.Image
//
// Will panic if width or height are negative or higher than MAX_TEXTURE_SIZE
func (g *OpenGL) NewImage(width, height int) *image.Image {
	if width < 0 {
		panic("negative width")
	}
	if height < 0 {
		panic("negative height")
	}
	acceleratedImage := g.context.NewAcceleratedImage(width, height)
	return image.New(acceleratedImage)
}

// mouseWindow implements mouse.Window
type mouseWindow struct {
	mutex          sync.Mutex
	glfwWindow     *glfw.Window
	mainThreadLoop *MainThreadLoop
	zoom           int
	width, height  int
}

func (m *mouseWindow) CursorPosition() (float64, float64) {
	var x, y float64
	m.mainThreadLoop.Execute(func() {
		x, y = m.glfwWindow.GetCursorPos()
	})
	return x, y
}

// Size() is thread-safe
func (m *mouseWindow) Size() (int, int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.width == 0 {
		m.mainThreadLoop.Execute(func() {
			m.width, m.height = m.glfwWindow.GetSize()
		})
	}
	return m.width, m.height
}

func (m *mouseWindow) Zoom() int {
	return m.zoom
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
	screenAcceleratedImage := g.context.NewAcceleratedImage(width, height)
	screenImage := image.New(screenAcceleratedImage)
	win := &Window{
		mainThreadLoop:   g.mainThreadLoop,
		keyboardEvents:   keyboardEvents,
		requestedWidth:   width,
		requestedHeight:  height,
		screenImage:      screenImage,
		screenContextAPI: g.context.API(),
		zoom:             1,
	}
	var err error
	g.mainThreadLoop.Execute(func() {
		win.glfwWindow, err = createWindow(g.mainThreadLoop, g.mainWindow)
		if err != nil {
			return
		}
		for _, option := range options {
			if option == nil {
				log.Println("nil option given when opening the window")
				continue
			}
			option(win)
		}
		win.mouseWindow = &mouseWindow{
			glfwWindow:     win.glfwWindow,
			mainThreadLoop: g.mainThreadLoop,
			zoom:           win.zoom,
		}
		win.mouseEvents = internal.NewMouseEvents(
			mouse.NewEventBuffer(32),
			win.mouseWindow)
		win.glfwWindow.SetKeyCallback(win.keyboardEvents.OnKeyCallback)
		win.glfwWindow.SetMouseButtonCallback(win.mouseEvents.OnMouseButtonCallback)
		win.glfwWindow.SetScrollCallback(win.mouseEvents.OnScrollCallback)
		win.glfwWindow.SetSize(win.requestedWidth*win.zoom, win.requestedHeight*win.zoom)
		win.glfwWindow.Show()
	})
	if err != nil {
		return nil, err
	}
	win.api = &context{
		run: func(f func()) {
			g.mainThreadLoop.executeCommand(command{
				window:  win.glfwWindow,
				execute: f,
			})
		},
		runAsync: func(f func()) {
			g.mainThreadLoop.executeAsyncCommand(command{
				window:  win.glfwWindow,
				execute: f,
			})
		},
	}
	win.context = gl.NewContext(win.api)
	win.screenPolygon = newScreenPolygon(win.context, win.api)
	win.program, err = compileProgram(win.context, vertexShaderSrc, fragmentShaderSrc)
	if err != nil {
		return nil, err
	}
	// in this window context there is only one program used with one texture
	win.api.UseProgram(win.program.ID())
	win.api.BindTexture(gl33.TEXTURE_2D, screenAcceleratedImage.TextureID())
	return win, nil
}

// Context returns OpenGL's context. It's methods can be invoked from any goroutine.
// Each invocation will return the same instance.
func (g *OpenGL) Context() *gl.Context {
	return g.context
}

// ContextAPI returns gl.API, which can be used to OpenGL direct access.
// It's methods can be invoked from any goroutine.
func (g *OpenGL) ContextAPI() gl.API {
	return g.api
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
