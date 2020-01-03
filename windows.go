package pixiq

// NewWindows returns a factory of Window objects.
func NewWindows(images *Images) *Windows {
	return &Windows{images: images}
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

type Windows struct {
	images *Images
}

// Loop starts the main loop. It will execute onEachFrame function for each frame, as soon as window is closed. This
// function blocks the current goroutine.
func (w *Windows) Loop(window Window, onEachFrame func(frame *Frame)) {
	frame := &Frame{}
	frame.image = w.images.New(window.Width(), window.Height())
	for !frame.closeWindow {
		onEachFrame(frame)
		window.Draw(frame.image)
		window.SwapImages()
	}
	window.Close()
}

// Frame provides the whole window image which will be drawn on a screen after making modifications
type Frame struct {
	closeWindow bool
	image       *Image
}

// CloseWindowEventually closes the window as soon as onEachFrame function is finished
func (w *Frame) CloseWindowEventually() {
	w.closeWindow = true
}

// Image returns the whole window Image, which can be modified
func (w *Frame) Image() *Image {
	return w.image
}
