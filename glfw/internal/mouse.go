package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/elgopher/pixiq/mouse"
)

// MouseEvents maps GLFW events to mouse.Event. Mapped events can be
// polled using mouse.EventSource interface.
type MouseEvents struct {
	buffer             *mouse.EventBuffer
	window             Window
	lastPosX, lastPosY float64
}

// NewMouseEvents creates *MouseEvents using given buffer and window. Based on the
// information returned by Window mouse move events are generated.
func NewMouseEvents(buffer *mouse.EventBuffer, window Window) *MouseEvents {
	if buffer == nil {
		panic("nil buffer")
	}
	if window == nil {
		panic("nil window")
	}
	return &MouseEvents{buffer: buffer, window: window}
}

var mouseButtonMapping = map[glfw.MouseButton]mouse.Button{
	glfw.MouseButtonLeft:   mouse.Left,
	glfw.MouseButtonRight:  mouse.Right,
	glfw.MouseButtonMiddle: mouse.Middle,
	glfw.MouseButton4:      mouse.Button4,
	glfw.MouseButton5:      mouse.Button5,
	glfw.MouseButton6:      mouse.Button6,
	glfw.MouseButton7:      mouse.Button7,
	glfw.MouseButton8:      mouse.Button8,
}

// OnMouseButtonCallback passes GLFW mouse event
func (e *MouseEvents) OnMouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	btn, ok := mouseButtonMapping[button]
	if !ok {
		return
	}
	switch action {
	case glfw.Press:
		e.buffer.Add(mouse.NewPressedEvent(btn))
	case glfw.Release:
		e.buffer.Add(mouse.NewReleasedEvent(btn))
	}
}

// OnScrollCallback passes GLFW mouse event
func (e *MouseEvents) OnScrollCallback(_ *glfw.Window, xoff float64, yoff float64) {
	e.buffer.Add(mouse.NewScrolledEvent(-xoff, -yoff))
}

// Window is an abstraction for getting information about cursor position, size and zoom.
// It is needed for generating mouse move events
type Window interface {
	CursorPosition() (float64, float64)
	Size() (int, int)
	Zoom() int
}

// Poll return next mapped event
func (e *MouseEvents) Poll() (mouse.Event, bool) {
	event, ok := e.buffer.Poll()
	if ok {
		return event, ok
	}
	// generate move event, because GLFW does not provide move events for Linux
	// and Windows when cursor is outside window.
	realX, realY := e.window.CursorPosition()
	if e.lastPosX != realX || e.lastPosY != realY {
		w, h := e.window.Size()
		zoom := float64(e.window.Zoom())
		insideWindow := true
		if int(realX) >= w || int(realY) >= h || realX < 0 || realY < 0 {
			insideWindow = false
		}
		e.lastPosX = realX
		e.lastPosY = realY
		return mouse.NewMovedEvent(int(realX/zoom), int(realY/zoom), realX, realY, insideWindow), true
	}
	return mouse.EmptyEvent, false
}
