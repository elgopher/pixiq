package pixiq

// NewScreens returns a new instance of Screens
func NewScreens(images *Images) *Screens {
	return &Screens{images: images}
}

// Screens is an abstraction for interaction with screens in a platform-agnostic way
type Screens struct {
	images *Images
}

// Screen is a place where images can be drawn
type Screen interface {
	// Draw draws image spanning the whole screen. If double buffering is used it may draw to the invisible buffer.
	Draw(image *Image)
	// SwapImages makes last drawn image visible (if double buffering was used, otherwise it may be a no-op)
	SwapImages()
	// Width returns the width of the screen in pixels. If zooming is used the width is not multiplied by zoom.
	Width() int
	// Height returns the height of the screen in pixels. If zooming is used the height is not multiplied by zoom.
	Height() int
}

// Loop starts the main loop. It will execute onEachFrame function for each frame, as soon as loop is stopped. This
// function blocks the current goroutine.
func (w *Screens) Loop(screen Screen, onEachFrame func(frame *Frame)) {
	frame := &Frame{}
	image := w.images.New(screen.Width(), screen.Height())
	frame.screen = image.WholeImageSelection()
	for !frame.loopStopped {
		onEachFrame(frame)
		screen.Draw(image)
		screen.SwapImages()
	}
}

// Frame contains information about the screen's image which will be drawn on screen after making modifications
type Frame struct {
	loopStopped bool
	screen      Selection
}

// StopLoopEventually stops the loop as soon as onEachFrame function is finished
func (f *Frame) StopLoopEventually() {
	f.loopStopped = true
}

// Screen returns the image Selection, which can be modified
func (f *Frame) Screen() Selection {
	return f.screen
}
