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
	"reflect"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	program2 "github.com/jacekolszak/pixiq/glsl/program"
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
	openGL := &OpenGL{
		mainThreadLoop:     mainThreadLoop,
		bindWindowToThread: mainThreadLoop.bind(mainWindow),
		stopPollingEvents:  make(chan struct{}),
		mainWindow:         mainWindow,
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
	mainThreadLoop.bind(win)()
	if err := gl.Init(); err != nil {
		return nil, err
	}
	return win, nil
}

// OpenGL provides method for creating OpenGL-accelerated image.Image and opening
// windows.
type OpenGL struct {
	mainThreadLoop     *MainThreadLoop
	bindWindowToThread func()
	stopPollingEvents  chan struct{}
	mainWindow         *glfw.Window
}

// Destroy cleans all the OpenGL resources associated with this instance.
// This method has to be called in integration tests to clean resources after
// each test. Otherwise on some platforms you may reach the limit of active
// OpenGL contexts.
func (g *OpenGL) Destroy() {
	g.stopPollingEvents <- struct{}{}
	g.mainThreadLoop.Execute(func() {
		g.mainThreadLoop.bind(g.mainWindow)()
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
// Will return error if width or height are negative or image of this dimensions
// cannot be created on a video card. (For instance when dimensions are not
// a power of two)
func (g *OpenGL) NewImage(width, height int) (*image.Image, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	acceleratedImage, err := g.NewAcceleratedImage(width, height)
	if err != nil {
		return nil, err
	}
	return image.New(width, height, acceleratedImage)
}

// NewAcceleratedImage returns an OpenGL-accelerated implementation of image.AcceleratedImage
// Will return error if width or height are negative or image of this dimensions
// cannot be created on a video card. (For instance when dimensions are not
// a power of two)
func (g *OpenGL) NewAcceleratedImage(width, height int) (image.AcceleratedImage, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	return g.newTexture(width, height)
}

func (g *OpenGL) newTexture(width, height int) (*texture, error) {
	var id uint32
	var err error
	g.mainThreadLoop.Execute(func() {
		g.bindWindowToThread()
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
	return &texture{
		id:                 id,
		width:              width,
		height:             height,
		mainThreadLoop:     g.mainThreadLoop,
		bindWindowToThread: g.bindWindowToThread,
	}, nil
}

type texture struct {
	id                 uint32
	width, height      int
	mainThreadLoop     *MainThreadLoop
	bindWindowToThread func()
}

func (t *texture) Modify(location image.AcceleratedFragmentLocation, call image.AcceleratedCall) {
	panic("implement me")
}

func (t *texture) TextureID() uint32 {
	return t.id
}

func (t *texture) Upload(input image.AcceleratedFragmentPixels) {
	// TODO EXPERIMENTAL IMPLEMENTATION JUST TO PROVE THAT MY ABSTRACTION CAN BE IMPLEMENTED!
	t.mainThreadLoop.Execute(func() {
		t.bindWindowToThread()
		gl.BindTexture(gl.TEXTURE_2D, t.id)
		gl.PixelStorei(gl.UNPACK_ROW_LENGTH, int32(input.Stride))
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			int32(input.Location.X),
			int32(input.Location.Y),
			int32(input.Location.Width),
			int32(input.Location.Height),
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(
				reflect.
					ValueOf(input.Pixels).
					Index(input.StartingPosition).
					UnsafeAddr(),
			),
		)
	})
}

func (t *texture) Download(output image.AcceleratedFragmentPixels) {
	//t.mainThreadLoop.Execute(func() {
	//	t.bindWindowToThread()
	//	gl.BindTexture(gl.TEXTURE_2D, t.id)
	//	gl.GetTexImage(
	//		gl.TEXTURE_2D,
	//		0,
	//		gl.RGBA,
	//		gl.UNSIGNED_BYTE,
	//		gl.Ptr(output),
	//	)
	//})
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
	screenTexture, err := g.newTexture(width, height)
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

func (g *OpenGL) DrawProgram() program2.Draw {
	return &drawProgram{mainThreadLoop: g.mainThreadLoop, glfwWindow: g.mainWindow}
}

type drawProgram struct {
	vertexShaderSrc   string
	fragmentShaderSrc string
	mainThreadLoop    *MainThreadLoop
	glfwWindow        *glfw.Window
}

func (g *drawProgram) Compile() (program2.CompiledDraw, error) {
	var err error
	var program *program
	g.mainThreadLoop.Execute(func() {
		g.mainThreadLoop.bind(g.glfwWindow)()
		program, err = compileProgram(g.vertexShaderSrc, g.fragmentShaderSrc)
	})
	if err != nil {
		return nil, err
	}
	return &compiledDraw{program: program, glfwWindow: g.glfwWindow, mainThreadLoop: g.mainThreadLoop}, nil
}

func (g *drawProgram) SetVertexShader(glsl string) {
	g.vertexShaderSrc = glsl
}

func (g *drawProgram) SetFragmentShader(glsl string) {
	g.fragmentShaderSrc = glsl
}

type compiledDraw struct {
	mainThreadLoop *MainThreadLoop
	glfwWindow     *glfw.Window
	program        *program
}

func (c *compiledDraw) NewVertexArrayObject() program2.VertexArrayObject {
	var id uint32
	c.mainThreadLoop.Execute(func() {
		c.mainThreadLoop.bind(c.glfwWindow)()
		gl.GenVertexArrays(1, &id)
	})
	return &vao{id: id, mainThreadLoop: c.mainThreadLoop, glfwWindow: c.glfwWindow}
}

func (c *compiledDraw) Delete() {
	panic("implement me")
}

func (c *compiledDraw) NewCall(f func(call program2.DrawCall)) image.AcceleratedCall {
	return func() {
		c.mainThreadLoop.Execute(func() {
			c.mainThreadLoop.bind(c.glfwWindow)()
			gl.UseProgram(c.program.id)
			f(&drawCall{})
		})
	}
}

func (c *compiledDraw) GetVertexAttributeLocation(name string) int {
	var location int
	c.mainThreadLoop.Execute(func() {
		c.mainThreadLoop.bind(c.glfwWindow)()
		location = int(gl.GetAttribLocation(c.program.id, gl.Str(name+"\x00")))
	})
	return location
}

func (c *compiledDraw) GetUniformLocation(name string) int {
	panic("implement me")
}

type drawCall struct {
}

func (d *drawCall) BindVertexArrayObject(program2.VertexArrayObject) {
	panic("implement me")
}

func (d *drawCall) BindTexture0(img *image.Image) {
	panic("implement me")
}

func (d *drawCall) SetFloatUniform(location int, val float32) {
	panic("implement me")
}

func (d *drawCall) SetIntUniform(location int, val int) {
	panic("implement me")
}

func (d *drawCall) SetMatrix4Uniform(location int, val [16]float32) {
	panic("implement me")
}

func (d *drawCall) Draw(mode program2.Mode, first, count int) {
	panic("implement me")
}

func (g *OpenGL) NewFloatVertexBuffer(program2.BufferUsage) program2.FloatVertexBuffer {
	var id uint32
	g.mainThreadLoop.Execute(func() {
		g.mainThreadLoop.bind(g.mainWindow)()
		gl.GenBuffers(1, &id)
	})
	return &floatVBO{id: id, mainThreadLoop: g.mainThreadLoop, glfwWindow: g.mainWindow}
}

func (g *OpenGL) CompileVertexShader(src string) (*shader, error) {
	return compileVertexShader(src)
}

func (g *OpenGL) CompileFragmentShader(src string) (*shader, error) {
	return compileFragmentShader(src)
}

func (g *OpenGL) LinkProgram(vertexShader *shader, fragmentShader *shader) {

}

type floatVBO struct {
	mainThreadLoop *MainThreadLoop
	glfwWindow     *glfw.Window
	id             uint32
}

func (f *floatVBO) Update(offset int, data []float32) {
	f.mainThreadLoop.Execute(func() {
		f.mainThreadLoop.bind(f.glfwWindow)()
		gl.BufferSubData(f.id, offset, len(data)*4, gl.Ptr(data))
	})
}

func (f *floatVBO) Pointer(start, length, stride int) program2.VertexBufferPointer {
	return program2.VertexBufferPointer{
		ID:     f.id,
		Size:   length,
		GLType: gl.FLOAT,
		Stride: stride * 4,
		Offset: start * 4,
	}
}

func (f *floatVBO) Delete() {
	panic("implement me")
}

type vao struct {
	id             uint32
	mainThreadLoop *MainThreadLoop
	glfwWindow     *glfw.Window
}

func (v *vao) SetVertexAttribute(location int, pointer program2.VertexBufferPointer) {
	v.mainThreadLoop.Execute(func() {
		v.mainThreadLoop.bind(v.glfwWindow)()
		gl.BindVertexArray(v.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, pointer.ID)
		gl.VertexAttribPointer(
			uint32(location),
			int32(pointer.Size),
			gl.FLOAT,
			false,
			int32(pointer.Stride),
			gl.PtrOffset(pointer.Offset),
		)
	})
}

func (v *vao) Delete() {
	panic("implement me")
}
