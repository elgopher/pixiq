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
	glfwWindow      *glfw.Window
	mainThreadLoop  *MainThreadLoop
	keyboardEvents  *internal.KeyboardEvents
	mouseEvents     *internal.MouseEvents
	requestedWidth  int
	requestedHeight int
	zoom            int
	title           string
	mouseWindow     *mouseWindow
	onClose         func(*Window)
	closed          bool
	drawer          windowDrawer
}

type windowDrawer struct {
	glfwWindow      *glfw.Window
	mainThreadLoop  *MainThreadLoop
	screenPolygon   *screenPolygon
	screenImage     *image.Image
	screenTextureID uint32
	sharedContext   *gl.Context // API for main context shared between all windows
	context         *gl.Context
	program         *gl.Program
}

func newWindow(glfwWindow *glfw.Window, mainThreadLoop *MainThreadLoop, width, height int, context, sharedContext *gl.Context, onClose func(*Window), options ...WindowOption) (*Window, error) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	// FIXME: EventBuffer size should be configurable
	keyboardEvents := internal.NewKeyboardEvents(keyboard.NewEventBuffer(32))
	drawer, err := newWindowDrawer(glfwWindow, mainThreadLoop, width, height, context, sharedContext)
	if err != nil {
		return nil, err
	}
	win := &Window{
		glfwWindow:      glfwWindow,
		mainThreadLoop:  mainThreadLoop,
		keyboardEvents:  keyboardEvents,
		requestedWidth:  width,
		requestedHeight: height,
		zoom:            1,
		title:           "OpenGL Pixiq Window",
		onClose:         onClose,
		drawer:          drawer,
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
	})
	<-sizeIsSet
	return win, nil
}

func newWindowDrawer(glfwWindow *glfw.Window, mainThreadLoop *MainThreadLoop, width, height int, context, sharedContext *gl.Context) (windowDrawer, error) {
	screenAcceleratedImage := sharedContext.NewAcceleratedImage(width, height)
	program, err := compileProgram(context, vertexShaderSrc, fragmentShaderSrc)
	if err != nil {
		return windowDrawer{}, err
	}
	return windowDrawer{
		glfwWindow:      glfwWindow,
		mainThreadLoop:  mainThreadLoop,
		screenPolygon:   newScreenPolygon(context),
		screenImage:     image.New(screenAcceleratedImage),
		screenTextureID: screenAcceleratedImage.TextureID(),
		sharedContext:   sharedContext,
		context:         context,
		program:         program,
	}, nil
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
	w.drawer.drawIntoBackBuffer()
}

func (d *windowDrawer) drawIntoBackBuffer() {
	d.screenImage.Upload()
	// Finish actively polls GPU which may consume a lot of CPU power.
	// That's why Finish is called only if context synchronization is required
	api := d.context.API()
	if d.sharedContext.API() != api {
		d.sharedContext.API().Finish()
	}
	var width, height int
	d.mainThreadLoop.Execute(func() {
		width, height = d.glfwWindow.GetFramebufferSize()
	})
	api.Disable(gl33.BLEND)
	api.Disable(gl33.SCISSOR_TEST)
	api.BindFramebuffer(gl33.FRAMEBUFFER, 0)
	api.Viewport(0, 0, int32(width), int32(height))
	api.BindTexture(gl33.TEXTURE_2D, d.screenTextureID)
	api.UseProgram(d.program.ID())
	d.screenPolygon.draw()
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
		w.glfwWindow.SetKeyCallback(nil)
		w.glfwWindow.SetMouseButtonCallback(nil)
		w.glfwWindow.SetScrollCallback(nil)
		w.glfwWindow.Hide()
	})
	w.drawer.close()
	w.onClose(w)
	w.closed = true
}

func (d *windowDrawer) close() {
	d.screenPolygon.delete()
	d.program.Delete()
	d.screenImage.Delete()
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
func (w *Window) Width() int {
	width, _ := w.mouseWindow.Size()
	return width
}

// Height returns the actual height of the window in pixels. It may be different
// than requested height used when window was open due to platform limitation.
// If zooming is used the height is multiplied by zoom.
func (w *Window) Height() int {
	_, height := w.mouseWindow.Size()
	return height
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
	return w.drawer.screenImage.WholeImageSelection()
}

// ContextAPI returns window-specific OpenGL's context. Useful for accessing
// window's framebuffer.
func (w *Window) ContextAPI() gl.API {
	return w.drawer.context.API()
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
