package glfw

import (
	"log"

	gl33 "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/glfw/internal"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/mouse"
)

// Window is an implementation of loop.Screen and keyboard.EventSource
type Window struct {
	glfwWindow             *glfw.Window
	mainThreadLoop         *MainThreadLoop
	screenPolygon          *screenPolygon
	keyboardEvents         *internal.KeyboardEvents
	mouseEvents            *internal.MouseEvents
	requestedWidth         int
	requestedHeight        int
	zoom                   int
	title                  string
	screenImage            *image.Image
	screenAcceleratedImage *gl.AcceleratedImage
	sharedContext          *gl.Context // API for main context shared between all windows
	context                *gl.Context
	program                *gl.Program
	mouseWindow            *mouseWindow
	onClose                func(*Window)
	closed                 bool
	fullScreenMode         *VideoMode
	sizeBefore             *sizeBeforeEnteringFullScreen
}

type sizeBeforeEnteringFullScreen struct {
	x, y, width, height, zoom int
}

func newWindow(glfwWindow *glfw.Window, mainThreadLoop *MainThreadLoop, width, height int, context *gl.Context, sharedContext *gl.Context, onClose func(*Window), options ...WindowOption) (*Window, error) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	screenAcceleratedImage := sharedContext.NewAcceleratedImage(width, height)
	program, err := compileProgram(context, vertexShaderSrc, fragmentShaderSrc)
	if err != nil {
		return nil, err
	}
	win := &Window{
		glfwWindow:             glfwWindow,
		mainThreadLoop:         mainThreadLoop,
		screenPolygon:          newScreenPolygon(context),
		keyboardEvents:         internal.NewKeyboardEvents(keyboard.NewEventBuffer(32)), // FIXME: EventBuffer size should be configurable
		requestedWidth:         width,
		requestedHeight:        height,
		zoom:                   1,
		title:                  "OpenGL Pixiq Window",
		screenImage:            image.New(screenAcceleratedImage),
		screenAcceleratedImage: screenAcceleratedImage,
		sharedContext:          sharedContext,
		context:                context,
		program:                program,
		onClose:                onClose,
	}
	var sizeIsSet <-chan bool
	mainThreadLoop.Execute(func() {
		for _, option := range options {
			if option == nil {
				log.Println("nil option given when opening the window")
				continue
			}
			option(win)
		}
		win.mouseWindow = &mouseWindow{
			glfwWindow: win.glfwWindow,
			zoom:       win.zoom,
		}
		win.mouseEvents = internal.NewMouseEvents(
			mouse.NewEventBuffer(32), // FIXME: EventBuffer size should be configurable
			win.mouseWindow)
		win.glfwWindow.SetKeyCallback(win.keyboardEvents.OnKeyCallback)
		win.glfwWindow.SetMouseButtonCallback(win.mouseEvents.OnMouseButtonCallback)
		win.glfwWindow.SetScrollCallback(win.mouseEvents.OnScrollCallback)
		sizeIsSet = updateSize(win)
		win.glfwWindow.Show()
		// monitor can be set only after window is shown
		videoMode := win.fullScreenMode
		if videoMode != nil {
			win.glfwWindow.SetMonitor(videoMode.monitor, 0, 0, videoMode.Width(), videoMode.Height(), videoMode.RefreshRate())
		}
	})
	<-sizeIsSet
	mainThreadLoop.Execute(func() {
		win.glfwWindow.SetSizeCallback(func(w *glfw.Window, width int, height int) {
			win.requestedWidth = width / win.zoom
			win.requestedHeight = height / win.zoom
			if width%win.zoom != 0 {
				win.requestedWidth++
			}
			if height%win.zoom != 0 {
				win.requestedHeight++
			}
		})
	})
	return win, nil
}

func updateSize(win *Window) <-chan bool {
	done := make(chan bool)
	newWidth := win.requestedWidth * win.zoom
	newHeight := win.requestedHeight * win.zoom
	currentWidth, currentHeight := win.glfwWindow.GetSize()
	if currentWidth != newWidth || currentHeight != newHeight {
		win.glfwWindow.SetSizeCallback(func(w *glfw.Window, width int, height int) {
			win.glfwWindow.SetSizeCallback(nil)
			close(done)
		})
		win.glfwWindow.SetSize(newWidth, newHeight)
	} else {
		close(done)
	}
	return done
}

// PollMouseEvent retrieves and removes next mouse Event. If there are no more
// events false is returned. It implements mouse.EventSource method.
func (w *Window) PollMouseEvent() (event mouse.Event, ok bool) {
	w.mainThreadLoop.Execute(func() {
		event, ok = w.mouseEvents.Poll()
	})
	return
}

// Draw draws a screen image in the window
func (w *Window) Draw() {
	if w.closed {
		panic("Draw forbidden for a closed window")
	}
	w.DrawIntoBackBuffer()
	w.SwapBuffers()
}

// DrawIntoBackBuffer draws a screen image into the back buffer. To make it visible
// to the user SwapBuffers must be executed.
func (w *Window) DrawIntoBackBuffer() {
	if w.closed {
		panic("DrawIntoBackBuffer forbidden for a closed window")
	}
	w.screenImage.Upload()
	// Finish actively polls GPU which may consume a lot of CPU power.
	// That's why Finish is called only if context synchronization is required
	api := w.context.API()
	if w.sharedContext.API() != api {
		w.sharedContext.API().Finish()
	}
	var width, height int
	w.mainThreadLoop.Execute(func() {
		width, height = w.glfwWindow.GetFramebufferSize()
	})
	api.Disable(gl33.BLEND)
	api.Disable(gl33.SCISSOR_TEST)
	api.BindFramebuffer(gl33.FRAMEBUFFER, 0)
	api.Viewport(0, 0, int32(width), int32(height))
	api.BindTexture(gl33.TEXTURE_2D, w.screenAcceleratedImage.TextureID())
	api.UseProgram(w.program.ID())
	xRight := float32(2*w.screenImage.Width()*w.zoom)/float32(width) - 1
	yBottom := -(float32(2*w.screenImage.Height()*w.zoom) / float32(height)) + 1
	w.screenPolygon.draw(xRight, yBottom)
}

// SwapBuffers makes current back buffer visible to the user.
func (w *Window) SwapBuffers() {
	if w.closed {
		panic("SwapBuffers forbidden for a closed window")
	}
	w.mainThreadLoop.Execute(w.glfwWindow.SwapBuffers)
}

// Close closes the window and cleans resources.
func (w *Window) Close() {
	if w.closed {
		return
	}
	w.mainThreadLoop.Execute(func() {
		w.glfwWindow.SetSizeCallback(nil)
		w.glfwWindow.SetKeyCallback(nil)
		w.glfwWindow.SetMouseButtonCallback(nil)
		w.glfwWindow.SetScrollCallback(nil)
		w.glfwWindow.Hide()
	})
	w.screenPolygon.delete()
	w.program.Delete()
	w.screenImage.Delete()
	w.onClose(w)
	w.closed = true
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
// If zooming is used the width is multiplied by zoom.
func (w *Window) Width() (width int) {
	w.mainThreadLoop.Execute(func() {
		width, _ = w.mouseWindow.Size()
	})
	return
}

// Height returns the actual height of the window in pixels. It may be different
// than requested height used when window was open due to platform limitation.
// If zooming is used the height is multiplied by zoom.
func (w *Window) Height() (height int) {
	w.mainThreadLoop.Execute(func() {
		_, height = w.mouseWindow.Size()
	})
	return
}

// X returns the X coordinate, in pixels, of the upper-left
// corner of the client area of the window.
func (w *Window) X() (x int) {
	w.mainThreadLoop.Execute(func() {
		x, _ = w.glfwWindow.GetPos()
	})
	return
}

// Y returns the Y coordinate, in pixels, of the upper-left
// corner of the client area of the window.
func (w *Window) Y() (y int) {
	w.mainThreadLoop.Execute(func() {
		_, y = w.glfwWindow.GetPos()
	})
	return
}

// Zoom returns the actual zoom. It is the zoom given during opening the window,
// unless zoom < 1 was given - then the actual zoom is 1.
func (w *Window) Zoom() int {
	return w.zoom
}

// PollKeyboardEvent retrieves and removes next keyboard Event. If there are no more
// events false is returned. It implements keyboard.EventSource method.
func (w *Window) PollKeyboardEvent() (event keyboard.Event, ok bool) {
	w.mainThreadLoop.Execute(func() {
		event, ok = w.keyboardEvents.Poll()
	})
	return
}

// Screen returns the image.Selection for the whole Window image
func (w *Window) Screen() image.Selection {
	var width, height int
	w.mainThreadLoop.Execute(func() {
		width = w.requestedWidth
		height = w.requestedHeight
	})
	w.ensureScreenSize(width, height)
	return w.screenImage.WholeImageSelection()
}

// ContextAPI returns window-specific OpenGL's context. Useful for accessing
// window's framebuffer.
func (w *Window) ContextAPI() gl.API {
	return w.context.API()
}

// SetCursor sets the window cursor
func (w *Window) SetCursor(cursor *Cursor) {
	if cursor == nil {
		panic("nil cursor")
	}
	w.mainThreadLoop.Execute(func() {
		w.glfwWindow.SetCursor(cursor.glfwCursor)
	})
}

// Title returns title of window
func (w *Window) Title() string {
	return w.title
}

// Resize changes the size of the window. Works only if full screen is off.
//
// Please note that retained Screen instance became obsolete after Resize.
// You have to call Window.Screen() again to get new screen
func (w *Window) Resize(width int, height, zoom int) {
	w.mainThreadLoop.Execute(func() {
		if w.fullScreenMode != nil {
			return
		}
		if zoom < 1 {
			zoom = 1
		}
		w.requestedWidth = width
		w.requestedHeight = height
		w.zoom = zoom
		newWidth := width * w.zoom
		newHeight := height * w.zoom
		w.glfwWindow.SetSize(newWidth, newHeight)
	})
}

func (w *Window) ensureScreenSize(width int, height int) {
	if w.screenImage.Width() != width || w.screenImage.Height() != height {
		newAcceleratedImage := w.sharedContext.NewAcceleratedImage(width, height)
		newImage := image.New(newAcceleratedImage)
		newSelection := newImage.WholeImageSelection()
		oldSelection := w.screenImage.WholeImageSelection()
		for y := 0; y < newImage.Height(); y++ {
			for x := 0; x < newImage.Width(); x++ {
				newSelection.SetColor(x, y, oldSelection.Color(x, y))
			}
		}
		w.screenImage.Delete()
		w.screenAcceleratedImage = newAcceleratedImage
		w.screenImage = newImage
	}
}

// SetPosition sets the position, in pixels, of the upper-left corner
// of the client area of the window.
func (w *Window) SetPosition(x int, y int) {
	w.mainThreadLoop.Execute(func() {
		w.glfwWindow.SetPos(x, y)
	})
}

// SetDecorationHint specifies whether the window will have window decorations
// such as a border, a close widget, etc.
func (w *Window) SetDecorationHint(enabled bool) {
	w.mainThreadLoop.Execute(func() {
		w.setBoolAttrib(glfw.Decorated, enabled)
	})
}

// Decorated returns true if window has decorations such as a border, a close widget, etc.
func (w *Window) Decorated() (decorated bool) {
	w.mainThreadLoop.Execute(func() {
		decorated = w.boolAttrib(glfw.Decorated)
	})
	return
}

func (w *Window) setBoolAttrib(hint glfw.Hint, enabled bool) {
	if enabled {
		w.glfwWindow.SetAttrib(hint, glfw.True)
	} else {
		w.glfwWindow.SetAttrib(hint, glfw.False)
	}
}

func (w *Window) boolAttrib(hint glfw.Hint) bool {
	return w.glfwWindow.GetAttrib(hint) == glfw.True
}

// EnterFullScreen makes window full screen using given display video mode
func (w *Window) EnterFullScreen(mode VideoMode, zoom int) {
	w.sizeBefore = &sizeBeforeEnteringFullScreen{
		x:      w.X(),
		y:      w.Y(),
		width:  w.requestedWidth,
		height: w.requestedHeight,
		zoom:   w.zoom,
	}
	w.zoom = zoom
	w.mainThreadLoop.Execute(func() {
		w.fullScreenMode = &mode
		w.glfwWindow.SetMonitor(mode.monitor, 0, 0, mode.Width(), mode.Height(), mode.RefreshRate())
	})
}

// ExitFullScreen exits from full screen and resizes the window to previous size
func (w *Window) ExitFullScreen() {
	var x, y, width, height, zoom int
	if w.sizeBefore != nil {
		x = w.sizeBefore.x
		y = w.sizeBefore.y
		width = w.sizeBefore.width
		height = w.sizeBefore.height
		zoom = w.sizeBefore.zoom
	} else {
		x = 0
		y = 0
		zoom = w.zoom
		width = w.requestedWidth
		height = w.requestedHeight
	}
	w.ExitFullScreenUsing(x, y, width, height, zoom)
}

// ExitFullScreenUsing exits from full screen and resizes the window
func (w *Window) ExitFullScreenUsing(x, y, width, height, zoom int) {
	w.mainThreadLoop.Execute(func() {
		w.fullScreenMode = nil
		w.requestedWidth = width
		w.requestedHeight = height
		w.zoom = zoom
		w.glfwWindow.SetMonitor(nil, x, y, width*zoom, height*zoom, 0)
	})
}

// SetAutoIconifyHint specifies whether fullscreen windows automatically iconify
// (and restore the previous video mode) on focus loss.
func (w *Window) SetAutoIconifyHint(enabled bool) {
	w.mainThreadLoop.Execute(func() {
		w.setBoolAttrib(glfw.AutoIconify, enabled)
	})
}
