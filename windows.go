package pixiq

// NewWindows returns a new instance of Windows
func NewWindows(images *Images) *Windows {
	return &Windows{images: images}
}

// Windows is an abstraction for interaction with windows in a platform-agnostic way
type Windows struct {
	images *Images
}

// Window is a window where images can be drawn
type Window interface {
	// Draw draws image spanning the whole window. If double buffering is used it may draw to the invisible buffer.
	Draw(image *Image)
	// SwapImages makes last drawn image visible (if double buffering was used, otherwise it may be a no-op)
	SwapImages()
	// Close closes the window and cleans resources
	Close()
	// Width returns the width of the window in pixels. If zooming is used the width is not multiplied by zoom.
	Width() int
	// Height returns the height of the window in pixels. If zooming is used the height is not multiplied by zoom.
	Height() int
}

// Loop starts the main loop. It will execute onEachFrame function for each frame, as soon as window is closed. This
// function blocks the current goroutine.
func (w *Windows) Loop(window Window, onEachFrame func(frame *Frame)) {
	frame := &Frame{}
	image := w.images.New(window.Width(), window.Height())
	frame.screen = image.WholeImageSelection()
	for !frame.closeWindow {
		onEachFrame(frame)
		window.Draw(image)
		window.SwapImages()
	}
	window.Close()
}

// Frame contains information about the current screen which will be drawn inside window after making modifications
type Frame struct {
	closeWindow bool
	screen      Selection
}

// CloseWindowEventually closes the window as soon as onEachFrame function is finished
func (f *Frame) CloseWindowEventually() {
	f.closeWindow = true
}

// Screens returns the whole window Image, which can be modified
func (f *Frame) Screen() Selection {
	return f.screen
}
