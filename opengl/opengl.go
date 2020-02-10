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
	"strings"
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
		vertexBufferIDs:   vertexBufferIDs{},
		textureIDs:        textureIDs{},
	}
	go openGL.startPollingEvents(openGL.stopPollingEvents)
	return openGL, nil
}

// vertexBufferIDs contains all vertex buffer identifiers in OpenGL context
type vertexBufferIDs map[*FloatVertexBuffer]uint32

// textureIDs contains all texture identifiers in OpenGL context
type textureIDs map[image.AcceleratedImage]uint32

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
	mainThreadLoop    *MainThreadLoop
	runInOpenGLThread func(func())
	stopPollingEvents chan struct{}
	mainWindow        *glfw.Window
	vertexBufferIDs   vertexBufferIDs
	textureIDs        textureIDs
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
	acceleratedImage, err := g.NewAcceleratedImage(width, height)
	if err != nil {
		return nil, err
	}
	return image.New(width, height, acceleratedImage)
}

// NewAcceleratedImage returns an OpenGL-accelerated implementation of image.AcceleratedImage
// Will return error if width or height are negative or image of these dimensions
// cannot be created on a video card. (For instance when dimensions are not
// a power of two)
func (g *OpenGL) NewAcceleratedImage(width, height int) (*AcceleratedImage, error) {
	if width < 0 {
		return nil, errors.New("negative width")
	}
	if height < 0 {
		return nil, errors.New("negative height")
	}
	var id uint32
	var err error
	g.runInOpenGLThread(func() {
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
	img := &AcceleratedImage{
		id:                id,
		width:             width,
		height:            height,
		runInOpenGLThread: g.runInOpenGLThread,
	}
	g.textureIDs[img] = id
	return img, nil
}

// AcceleratedImage is an image.AcceleratedImage implementation storing pixels
// on a video card VRAM.
type AcceleratedImage struct {
	id                uint32
	width, height     int
	runInOpenGLThread func(func())
}

// Upload send pixels to video card
func (t *AcceleratedImage) Upload(pixels []image.Color) {
	if len(pixels) == 0 {
		return
	}
	t.runInOpenGLThread(func() {
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

// Download gets pixels pixels from video card
func (t *AcceleratedImage) Download(output []image.Color) {
	if len(output) == 0 {
		return
	}
	t.runInOpenGLThread(func() {
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
	screenAcceleratedImage, err := g.NewAcceleratedImage(width, height)
	if err != nil {
		return nil, err
	}
	screenImage, err := image.New(width, height, screenAcceleratedImage)
	if err != nil {
		return nil, err
	}
	win := &Window{
		mainThreadLoop:         g.mainThreadLoop,
		keyboardEvents:         keyboardEvents,
		requestedWidth:         width,
		requestedHeight:        height,
		screenAcceleratedImage: screenAcceleratedImage,
		screenImage:            screenImage,
		zoom:                   1,
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
	g.runInOpenGLThread(func() {
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
	g.runInOpenGLThread(func() {
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
	var uniformLocations = map[string]int32{}
	g.runInOpenGLThread(func() {
		program, err = linkProgram(vertexShader.shader, fragmentShader.shader)
		if err == nil {
			uniformLocations = program.uniformAttributeLocations()
		}
	})
	if err != nil {
		return nil, err
	}
	return &Program{
		program:           program,
		runInOpenGLThread: g.runInOpenGLThread,
		uniformLocations:  uniformLocations,
		textureIDs:        g.textureIDs,
	}, err
}

// NewFloatVertexBuffer creates an OpenGL's Vertex Buffer Object (VBO) containing only float32 numbers.
func (g *OpenGL) NewFloatVertexBuffer(size int) (*FloatVertexBuffer, error) {
	if size < 0 {
		return nil, errors.New("negative size")
	}
	var id uint32
	g.runInOpenGLThread(func() {
		gl.GenBuffers(1, &id)
		gl.BindBuffer(gl.ARRAY_BUFFER, id)
		gl.BufferData(gl.ARRAY_BUFFER, size*4, gl.Ptr(nil), gl.STATIC_DRAW) // FIXME: Parametrize usage
	})
	vb := &FloatVertexBuffer{
		id:                id,
		size:              size,
		runInOpenGLThread: g.runInOpenGLThread,
	}
	g.vertexBufferIDs[vb] = id
	return vb, nil
}

type VertexLayout []Type

type Type struct {
	components int32
	xtype      uint32
}

var (
	Float  = Type{components: 1, xtype: gl.FLOAT}
	Float2 = Type{components: 2, xtype: gl.FLOAT}
	Float3 = Type{components: 3, xtype: gl.FLOAT}
	Float4 = Type{components: 4, xtype: gl.FLOAT}
)

func (g *OpenGL) NewVertexArray(layout VertexLayout) (*VertexArray, error) {
	if len(layout) == 0 {
		return nil, errors.New("empty layout")
	}
	var id uint32
	g.runInOpenGLThread(func() {
		// TODO: not tested at all
		gl.GenVertexArrays(1, &id)
		for i := 0; i < len(layout); i++ {
			gl.EnableVertexAttribArray(uint32(i))
		}
	})
	return &VertexArray{
		id:                id,
		layout:            layout,
		runInOpenGLThread: g.runInOpenGLThread,
		vertexBufferIDs:   g.vertexBufferIDs,
	}, nil
}

type VertexArray struct {
	id                uint32
	runInOpenGLThread func(func())
	layout            VertexLayout
	vertexBufferIDs   vertexBufferIDs
}

func (a *VertexArray) Delete() {
}

type VertexBufferPointer struct {
	Buffer *FloatVertexBuffer
	Offset int
	Stride int
}

func (a *VertexArray) Set(location int, pointer VertexBufferPointer) error {
	if pointer.Offset < 0 {
		return errors.New("negative pointer offset")
	}
	if pointer.Stride < 0 {
		return errors.New("negative pointer stride")
	}
	if pointer.Buffer == nil {
		return errors.New("nil pointer buffer")
	}
	if location < 0 {
		return errors.New("negative location")
	}
	if location >= len(a.layout) {
		return errors.New("location out-of-bounds")
	}
	bufferID, ok := a.vertexBufferIDs[pointer.Buffer]
	if !ok {
		return errors.New("vertex buffer has not been created in this context")
	}
	a.runInOpenGLThread(func() {
		// TODO: not tested at all
		gl.BindVertexArray(a.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, bufferID)
		typ := a.layout[location]
		components := typ.components
		gl.VertexAttribPointer(uint32(location), components, typ.xtype, false, int32(pointer.Stride*4), gl.PtrOffset(int(components*4)))
	})
	return nil
}

// FloatVertexBuffer is a struct representing OpenGL's Vertex Buffer Object (VBO) containing only float32 numbers.
type FloatVertexBuffer struct {
	id                uint32
	size              int
	runInOpenGLThread func(func())
	deleted           bool
}

// Size is the number of float values defined during creation time.
func (b *FloatVertexBuffer) Size() int {
	return b.size
}

// Download gets data starting at a given offset in VRAM and put them into slice. Whole output slice will be filled with data,
// unless output slice is bigger then the vertex buffer.
func (b *FloatVertexBuffer) Download(offset int, output []float32) error {
	if b.deleted {
		return errors.New("deleted buffer")
	}
	if offset < 0 {
		return errors.New("negative offset")
	}
	if len(output) == 0 {
		return nil
	}
	size := len(output)
	if size+offset > b.size {
		size = b.size - offset
	}
	b.runInOpenGLThread(func() {
		gl.BindBuffer(gl.ARRAY_BUFFER, b.id)
		gl.GetBufferSubData(gl.ARRAY_BUFFER, offset*4, size*4, gl.Ptr(output))
	})
	return nil
}

// Delete should be called whenever you don't plan to use vertex buffer anymore. Vertex Buffer is external resource
// (like file for example) and must be deleted manually
func (b *FloatVertexBuffer) Delete() {
	b.runInOpenGLThread(func() {
		gl.DeleteBuffers(1, &b.id)
	})
	b.deleted = true
}

// Upload sends data to the vertex buffer. All slice data will be inserted starting at a given offset position.
//
// Returns error when vertex buffer is too small to hold the data or offset is negative.
func (b *FloatVertexBuffer) Upload(offset int, data []float32) error {
	if offset < 0 {
		return errors.New("negative offset")
	}
	if b.size < len(data)+offset {
		return errors.New("FloatVertexBuffer is to small to store data")
	}
	var err error
	b.runInOpenGLThread(func() {
		gl.BindBuffer(gl.ARRAY_BUFFER, b.id)
		gl.BufferSubData(gl.ARRAY_BUFFER, offset*4, len(data)*4, gl.Ptr(data))
	})
	return err
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
	uniformLocations  map[string]int32
	runInOpenGLThread func(func())
	textureIDs        textureIDs
}

func (p *Program) AcceleratedCommand(command Command) (*AcceleratedCommand, error) {
	if command == nil {
		return nil, errors.New("nil command")
	}

	acceleratedCommand := &AcceleratedCommand{
		command:           command,
		runInOpenGLThread: p.runInOpenGLThread,
		program:           p,
		textureIDs:        p.textureIDs,
	}
	return acceleratedCommand, nil
}

func (p *Program) uniformAttributeLocation(name string) (int32, error) {
	location, ok := p.uniformLocations[name]
	if !ok {
		return 0, errors.New("not existing uniform attribute name")
	}
	return location, nil
}

type Command interface {
	// Implementations must not retain renderer and selections.
	RunGL(renderer *Renderer, selections []image.AcceleratedImageSelection) error
}

type Renderer struct {
	program           *Program
	runInOpenGLThread func(func())
	textureIDs        textureIDs
}

func (r *Renderer) BindTexture(textureUnit int, uniformAttributeName string, image image.AcceleratedImage) error {
	if textureUnit < 0 {
		return errors.New("negative textureUnit")
	}
	trimmed := strings.TrimSpace(uniformAttributeName)
	if trimmed == "" {
		return errors.New("empty uniformAttributeName")
	}
	textureLocation, err := r.program.uniformAttributeLocation(uniformAttributeName)
	if err != nil {
		return err
	}
	textureID, ok := r.textureIDs[image]
	if !ok {
		return errors.New("image has not been created in this OpenGL context")
	}
	r.runInOpenGLThread(func() {
		// TODO: not tested at all
		gl.Uniform1i(textureLocation, int32(textureUnit))
		gl.ActiveTexture(uint32(gl.TEXTURE0 + textureUnit))
		gl.BindTexture(gl.TEXTURE_2D, textureID)
	})
	return nil
}

type Mode struct {
	glMode uint32
}

var Triangles = Mode{
	glMode: gl.TRIANGLES,
}

var Points = Mode{
	glMode: gl.POINTS,
}

func (r *Renderer) DrawArrays(array *VertexArray, mode Mode, first, count int) {
	r.runInOpenGLThread(func() {
		gl.BindVertexArray(array.id)
		gl.DrawArrays(mode.glMode, int32(first), int32(count))
	})
}

func (r *Renderer) Clear(color image.Color) {
	r.runInOpenGLThread(func() {
		// TODO Move to image.Color
		r := float32(color.R()) / 255.0
		g := float32(color.G()) / 255.0
		b := float32(color.B()) / 255.0
		a := float32(color.A()) / 255.0
		gl.ClearColor(r, g, b, a)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

// AcceleratedCommand is an image.AcceleratedCommand implementation.
type AcceleratedCommand struct {
	command           Command
	program           *Program
	runInOpenGLThread func(func())
	textureIDs        textureIDs
}

func (c *AcceleratedCommand) Run(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) error {
	var err error
	var frameBuffer uint32
	c.runInOpenGLThread(func() {
		c.program.use()
		gl.Enable(gl.SCISSOR_TEST)
		gl.GenFramebuffers(1, &frameBuffer)
		gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)
		outputTextureID := c.textureIDs[output.Image]
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, outputTextureID, 0)
		loc := output.Location
		gl.Scissor(int32(loc.X), int32(loc.Y), int32(loc.Width), int32(loc.Height))
		gl.Viewport(int32(loc.X), int32(loc.Y), int32(loc.Width), int32(loc.Height))
	})
	renderer := &Renderer{
		program:           c.program,
		runInOpenGLThread: c.runInOpenGLThread,
		textureIDs:        c.textureIDs,
	}
	err = c.command.RunGL(renderer, selections)
	c.runInOpenGLThread(func() {
		gl.DeleteBuffers(1, &frameBuffer)
	})
	return err
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
