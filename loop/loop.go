// Package Loop provides screen looping functionality in a platform-agnostic way
package loop

import "github.com/jacekolszak/pixiq/image"

// Screen provides Image which can be manipulated and then drawn on the display.
type Screen interface {
	// Returns the image spanning the whole screen.
	Image() *image.Image
	// Draw draws the image on the screen.
	// If double buffering is used it may draw to the invisible buffer.
	Draw()
	// SwapImages makes last drawn image visible (if double buffering was used,
	// otherwise it may be a no-op)
	SwapImages()
}

// Run starts the screen loop. It will execute onEachFrame function for each frame,
// until loop is stopped. This function blocks the current goroutine.
func Run(screen Screen, onEachFrame func(frame *Frame)) {
	frame := &Frame{}
	for !frame.loopStopped {
		frame.screen = screen.Image().WholeImageSelection()
		onEachFrame(frame)
		screen.Draw()
		screen.SwapImages()
	}
}

// Frame contains information about the screen's image which will be drawn on screen
// after making modifications
type Frame struct {
	loopStopped bool
	screen      image.Selection
}

// StopLoopEventually stops the loop as soon as onEachFrame function is finished
func (f *Frame) StopLoopEventually() {
	f.loopStopped = true
}

// Screen returns the image Selection, which can be modified
func (f *Frame) Screen() image.Selection {
	return f.screen
}
