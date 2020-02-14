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
		allImages:         allImages{},
	}
	go openGL.startPollingEvents(openGL.stopPollingEvents)
	return openGL, nil
}

// vertexBufferIDs contains all vertex buffer identifiers in OpenGL context
type vertexBufferIDs map[VertexBuffer]uint32

type allImages map[image.AcceleratedImage]*AcceleratedImage

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
	// resizing the window to higher values than initial ones. That's why the window
	// created here has size equal to the biggest window used in integration tests
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
	allImages         allImages
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
		return nil, illegalArgumentError("negative width")
	}
	if height < 0 {
		return nil, illegalArgumentError("negative height")
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
		return nil, illegalArgumentError("negative width")
	}
	if height < 0 {
		return nil, illegalArgumentError("negative height")
	}
	var id uint32
	var frameBufferID uint32
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

		gl.GenFramebuffers(1, &frameBufferID)
		gl.BindFramebuffer(gl.FRAMEBUFFER, frameBufferID)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, id, 0)
	})
	if err != nil {
		return nil, err
	}
	img := &AcceleratedImage{
		textureID:         id,
		frameBufferID:     frameBufferID,
		width:             width,
		height:            height,
		runInOpenGLThread: g.runInOpenGLThread,
	}
	g.allImages[img] = img
	return img, nil
}

// AcceleratedImage is an image.AcceleratedImage implementation storing pixels
// on a video card VRAM.
type AcceleratedImage struct {
	textureID         uint32
	frameBufferID     uint32
	width, height     int
	runInOpenGLThread func(func())
}

// Upload send pixels to video card
func (t *AcceleratedImage) Upload(pixels []image.Color) {
	if len(pixels) == 0 {
		return
	}
	t.runInOpenGLThread(func() {
		gl.BindTexture(gl.TEXTURE_2D, t.textureID)
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
		gl.BindTexture(gl.TEXTURE_2D, t.textureID)
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
		win.screenPolygon = newScreenPolygon()
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
		return nil, illegalArgumentError("nil vertexShader")
	}
	if fragmentShader == nil {
		return nil, illegalArgumentError("nil fragmentShader")
	}
	var (
		program          *program
		err              error
		uniformLocations map[string]int32
		attributes       map[int32]attribute
	)
	g.runInOpenGLThread(func() {
		program, err = linkProgram(vertexShader.shader, fragmentShader.shader)
		if err == nil {
			uniformLocations = program.activeUniformLocations()
			attributes = program.attributes()
		}
	})
	if err != nil {
		return nil, err
	}
	return &Program{
		program:           program,
		runInOpenGLThread: g.runInOpenGLThread,
		uniformLocations:  uniformLocations,
		attributes:        attributes,
		allImages:         g.allImages,
	}, err
}

// NewFloatVertexBuffer creates an OpenGL's Vertex Buffer Object (VBO) containing only float32 numbers.
func (g *OpenGL) NewFloatVertexBuffer(size int) (*FloatVertexBuffer, error) {
	if size < 0 {
		return nil, illegalArgumentError("negative size")
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

// VertexLayout defines data types of VertexArray locations.
type VertexLayout []Type

// Type is a kind of OpenGL's attribute.
type Type struct {
	components int32
	xtype      uint32
	name       string
}

func valueOf(xtype uint32) Type {
	switch xtype {
	case gl.FLOAT:
		return Float
	case gl.FLOAT_VEC2:
		return Vec2
	case gl.FLOAT_VEC3:
		return Vec3
	case gl.FLOAT_VEC4:
		return Vec4
	}
	panic("not supported type")
}

func (t Type) String() string {
	return t.name
}

var (
	// Float is single-precision floating point number.
	// Equivalent of Go's float32.
	Float = Type{components: 1, xtype: gl.FLOAT, name: "Float"}
	// Vec2 is a vector of two single-precision floating point numbers.
	// Equivalent of Go's [2]float32.
	Vec2 = Type{components: 2, xtype: gl.FLOAT, name: "Vec2"}
	// Vec3 is a vector of three single-precision floating point numbers.
	// Equivalent of Go's [3]float32.
	Vec3 = Type{components: 3, xtype: gl.FLOAT, name: "Vec3"}
	// Vec4 is a vector of four single-precision floating point numbers.
	// Equivalent of Go's [4]float32.
	Vec4 = Type{components: 4, xtype: gl.FLOAT, name: "Vec4"}
)

// NewVertexArray creates a new instance of VertexArray. All vertex attributes
// specified in layout will be enabled.
func (g *OpenGL) NewVertexArray(layout VertexLayout) (*VertexArray, error) {
	if len(layout) == 0 {
		return nil, errors.New("empty layout")
	}
	var id uint32
	g.runInOpenGLThread(func() {
		gl.GenVertexArrays(1, &id)
		gl.BindVertexArray(id)
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

// VertexArray is a thin abstraction for OpenGL's Vertex Array Object.
//
// https://www.khronos.org/opengl/wiki/Vertex_Specification#Vertex_Array_Object
type VertexArray struct {
	id                uint32
	runInOpenGLThread func(func())
	layout            VertexLayout
	vertexBufferIDs   vertexBufferIDs
}

// Delete should be called whenever you don't plan to use VertexArray anymore.
// VertexArray is an external resource (like file for example) and must be deleted manually.
func (a *VertexArray) Delete() {
	a.runInOpenGLThread(func() {
		gl.DeleteVertexArrays(1, &a.id)
	})
}

// VertexBufferPointer is a slice of VertexBuffer
type VertexBufferPointer struct {
	Buffer VertexBuffer
	Offset int
	Stride int
}

// VertexBuffer contains data about vertices.
type VertexBuffer interface {
	// ID returns OpenGL identifier/name.
	ID() uint32
}

// IsClientError returns true if the error returned by opengl package methods
// are due to improper use of API.
func IsClientError(error error) bool {
	if _, ok := error.(illegalArgumentError); ok {
		return true
	}
	if _, ok := error.(illegalStateError); ok {
		return true
	}
	return false
}

type illegalArgumentError string

func (e illegalArgumentError) Error() string {
	return string(e)
}

type illegalStateError string

func (e illegalStateError) Error() string {
	return string(e)
}

// Set sets a location of VertexArray pointing to VertexBuffer slice.
func (a *VertexArray) Set(location int, pointer VertexBufferPointer) error {
	if pointer.Offset < 0 {
		return illegalArgumentError("negative pointer offset")
	}
	if pointer.Stride < 0 {
		return illegalArgumentError("negative pointer stride")
	}
	if pointer.Buffer == nil {
		return illegalArgumentError("nil pointer buffer")
	}
	if location < 0 {
		return illegalArgumentError("negative location")
	}
	if location >= len(a.layout) {
		return illegalArgumentError("location out-of-bounds")
	}
	bufferID, ok := a.vertexBufferIDs[pointer.Buffer]
	if !ok {
		return illegalStateError("vertex buffer has not been created in this context")
	}
	a.runInOpenGLThread(func() {
		gl.BindVertexArray(a.id)
		gl.BindBuffer(gl.ARRAY_BUFFER, bufferID)
		typ := a.layout[location]
		components := typ.components
		gl.VertexAttribPointer(uint32(location), components, typ.xtype, false, int32(pointer.Stride*4), gl.PtrOffset(pointer.Offset*4))
	})
	return nil
}

// FloatVertexBuffer is a struct representing OpenGL's Vertex Buffer Object (VBO) containing only float32 numbers.
type FloatVertexBuffer struct {
	id                uint32
	deleted           bool
	size              int
	runInOpenGLThread func(func())
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
		return illegalArgumentError("negative offset")
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
		return illegalArgumentError("negative offset")
	}
	if b.size < len(data)+offset {
		return illegalArgumentError("FloatVertexBuffer is to small to store data")
	}
	var err error
	b.runInOpenGLThread(func() {
		gl.BindBuffer(gl.ARRAY_BUFFER, b.id)
		gl.BufferSubData(gl.ARRAY_BUFFER, offset*4, len(data)*4, gl.Ptr(data))
	})
	return err
}

// ID returns OpenGL identifier/name.
func (b *FloatVertexBuffer) ID() uint32 {
	return b.id
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
	allImages         allImages
	attributes        map[int32]attribute
}

// AcceleratedCommand returns a potentially cached instance of *AcceleratedCommand.
func (p *Program) AcceleratedCommand(command Command) *AcceleratedCommand {
	return &AcceleratedCommand{
		command:           command,
		runInOpenGLThread: p.runInOpenGLThread,
		program:           p,
		allImages:         p.allImages,
	}
}

func (p *Program) uniformAttributeLocation(name string) (int32, error) {
	location, ok := p.uniformLocations[name]
	if !ok {
		return 0, errors.New("not existing uniform attribute name")
	}
	return location, nil
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
