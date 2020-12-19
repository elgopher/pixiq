// Package glfw makes it possible to use Pixiq on PCs with Linux, Windows or MacOS.
// It provides a method for creating OpenGL-accelerated image.Image and Window which
// is an implementation of loop.Screen and keyboard.EventSource.
// Under the hood it is using OpenGL API and GLFW for manipulating windows
// and handling user input.
package glfw

import (
	"time"

	gl33 "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/goimage"
	"github.com/jacekolszak/pixiq/image"
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
		mainWindow, err = createWindow(mainThreadLoop, "OpenGL Pixiq Window", nil)
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
	openGL := &OpenGL{
		mainThreadLoop:    mainThreadLoop,
		runInOpenGLThread: runInOpenGLThread,
		stopPollingEvents: make(chan struct{}),
		mainWindow:        mainWindow,
		context:           gl.NewContext(newContext(mainThreadLoop, mainWindow)),
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

func createWindow(mainThreadLoop *MainThreadLoop, title string, share *glfw.Window) (*glfw.Window, error) {
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
	win, err := glfw.CreateWindow(3, 3, title, nil, share)
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
	context           *gl.Context
	windowsOpen       int
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
	glfwWindow *glfw.Window
	zoom       int
}

func (m *mouseWindow) CursorPosition() (float64, float64) {
	return m.glfwWindow.GetCursorPos()
}

// Size() is thread-safe
func (m *mouseWindow) Size() (int, int) {
	return m.glfwWindow.GetSize()
}

func (m *mouseWindow) Zoom() int {
	return m.zoom
}

// OpenWindow creates and shows Window.
func (g *OpenGL) OpenWindow(width, height int, options ...WindowOption) (*Window, error) {
	glfwWindow := g.mainWindow
	winContext := g.context
	if g.windowsOpen > 0 {
		var err error
		g.mainThreadLoop.Execute(func() {
			glfwWindow, err = createWindow(g.mainThreadLoop, "OpenGL Pixiq Window", g.mainWindow)
		})
		if err != nil {
			return nil, err
		}
		api := newContext(g.mainThreadLoop, glfwWindow)
		winContext = gl.NewContext(api)
	}
	onClose := func(window *Window) {
		if glfwWindow != g.mainWindow {
			g.mainThreadLoop.Execute(glfwWindow.Destroy)
		}
		g.windowsOpen--
	}
	win, err := newWindow(glfwWindow, g.mainThreadLoop, width, height, winContext, g.context, onClose, options...)
	if err != nil {
		return nil, err
	}
	g.windowsOpen++
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
	return g.context.API()
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
		window.title = title
		window.glfwWindow.SetTitle(title)
	}
}

// Zoom makes window/pixels bigger zoom times. For zoom <= 0, the zoom defaults to 1.
func Zoom(zoom int) WindowOption {
	return func(window *Window) {
		if zoom > 0 {
			window.zoom = zoom
		}
	}
}

// NewCursor creates a new custom cursor look that can be set for a Window with SetCursor.
// The look is taken from a Selection. The size of the cursor is based on the Selection size
// and zoom.
//
// By default cursor has hotspot=(0,0) and zoom=1 . These values can be modified
// by providing slice of CursorOption: glfw.Hotspot(x,y) and glfw.CursorZoom(x,y)
func (g *OpenGL) NewCursor(selection image.Selection, options ...CursorOption) *Cursor {
	opts := cursorOpts{
		zoom: 1,
	}
	for _, option := range options {
		opts = option(opts)
	}
	if opts.hotspotX < 0 {
		opts.hotspotX = 0
	}
	if opts.hotspotY < 0 {
		opts.hotspotY = 0
	}
	if opts.hotspotX > selection.Width() {
		opts.hotspotX = selection.Width()
	}
	if opts.hotspotY > selection.Height() {
		opts.hotspotY = selection.Height()
	}
	rgbaImage := goimage.FromSelection(selection, goimage.Zoom(opts.zoom))
	var glfwCursor *glfw.Cursor
	g.mainThreadLoop.Execute(func() {
		glfwCursor = glfw.CreateCursor(rgbaImage, opts.hotspotX*opts.zoom, opts.hotspotY*opts.zoom)
	})
	return &Cursor{glfwCursor: glfwCursor}
}

type cursorOpts struct {
	zoom               int
	hotspotX, hotspotY int
}

// Cursor is a mouse cursor which can be use use in the window by calling Window.SetCursor
type Cursor struct {
	glfwCursor *glfw.Cursor
}

// Destroy frees the resources allocated by Cursor. This method must be called when
// cursor is not used anymore to avoid memory leakage.
func (c *Cursor) Destroy() {
	c.glfwCursor.Destroy()
}

// CursorOption is an option used when calling NewCursor
type CursorOption func(opts cursorOpts) cursorOpts

// Hotspot sets coordinates, in pixels, of cursor hotspot. Coordinates are constrained
// to cursor size. Coordinates are set to 0 if negative. If zoom was used hotspot
// coordinates are multiplied by zoom.
func Hotspot(x, y int) CursorOption {
	return func(opts cursorOpts) cursorOpts {
		opts.hotspotX = x
		opts.hotspotY = y
		return opts
	}
}

// CursorZoom makes cursor bigger zoom times. For zoom <= 1, the zoom defaults to 1.
func CursorZoom(zoom int) CursorOption {
	return func(opts cursorOpts) cursorOpts {
		opts.zoom = zoom
		return opts
	}
}

var cursorMapping = map[CursorShape]glfw.StandardCursor{
	Arrow:     glfw.ArrowCursor,
	IBeam:     glfw.IBeamCursor,
	Crosshair: glfw.CrosshairCursor,
	Hand:      glfw.HandCursor,
	HResize:   glfw.HResizeCursor,
	VResize:   glfw.VResizeCursor,
}

// NewStandardCursor creates a standard cursor with specified shape
func (g *OpenGL) NewStandardCursor(shape CursorShape) *Cursor {
	var glfwCursor *glfw.Cursor
	g.mainThreadLoop.Execute(func() {
		cursor, ok := cursorMapping[shape]
		if !ok {
			cursor = glfw.ArrowCursor
		}
		glfwCursor = glfw.CreateStandardCursor(cursor)
	})
	return &Cursor{glfwCursor: glfwCursor}
}

// CursorShape is a shape used by NewStandardCursor
type CursorShape int

const (
	// Arrow is an arrow cursor shape which can be used in NewStandardCursor
	Arrow CursorShape = iota
	// IBeam is an ibeam cursor shape which can be used in NewStandardCursor
	IBeam
	// Crosshair is a crosshair cursor shape which can be used in NewStandardCursor
	Crosshair
	// Hand is a hand cursor shape which can be used in NewStandardCursor
	Hand
	// HResize is a hresize cursor shape which can be used in NewStandardCursor
	HResize
	// VResize is a vresize cursor shape which can be used in NewStandardCursor
	VResize
)
