package pixiq

// NewWindows returns a factory of Window objects.
func NewWindows() *Windows {
	return &Windows{}
}

// Windows is a factory of Window objects
type Windows struct {
}

// New creates a new window with width and height given in pixels.
func (w *Windows) New(width, height int) *Window {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return &Window{
		width:  width,
		height: height,
	}
}

// Window is area where image will be drawn
type Window struct {
	width, height int
}

// Width returns width of the window in pixels
func (w *Window) Width() int {
	return w.width
}

// Height returns height of the window in pixels
func (w *Window) Height() int {
	return w.height
}

// Loop starts the main loop. It will execute onEachFrame function for each frame, as soon as window is closed. This
// function blocks the current goroutine.
func (w *Window) Loop(onEachFrame func(frame *Frame)) {
	frame := &Frame{}
	for !frame.closeWindow {
		onEachFrame(frame)
	}
}

// Frame is a current
type Frame struct {
	closeWindow bool
}

// CloseWindowEventually closes the window as soon as onEachFrame function is finished
func (w *Frame) CloseWindowEventually() {
	w.closeWindow = true
}
