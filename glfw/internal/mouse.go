package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/mouse"
)

type MouseEvents struct {
	buffer *mouse.EventBuffer
}

func NewMouseEvents(buffer *mouse.EventBuffer) *MouseEvents {
	if buffer == nil {
		panic("nil buffer")
	}
	return &MouseEvents{buffer: buffer}
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

func (e *MouseEvents) OnMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
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

func (e *MouseEvents) OnScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	e.buffer.Add(mouse.NewScrolledEvent(xoff, yoff))
}

// Poll return next mapped event
func (e *MouseEvents) Poll() (mouse.Event, bool) {
	return e.buffer.Poll()
}
