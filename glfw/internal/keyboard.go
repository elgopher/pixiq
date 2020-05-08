package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/keyboard"
)

// KeyboardEvents maps GLFW events to keyboard.Event. Mapped events can be
// polled using keyboard.EventSource interface.
type KeyboardEvents struct {
	buffer *keyboard.EventBuffer
}

// NewKeyboardEvents creates *KeyboardEvents using given buffer
func NewKeyboardEvents(buffer *keyboard.EventBuffer) *KeyboardEvents {
	if buffer == nil {
		panic("nil buffer")
	}
	return &KeyboardEvents{buffer: buffer}
}

// OnKeyCallback passes GLFW key event
func (e *KeyboardEvents) OnKeyCallback(_ *glfw.Window, glfwKey glfw.Key, scanCode int, action glfw.Action, _ glfw.ModifierKey) {
	key, ok := keymap[glfwKey]
	if !ok {
		key = keyboard.NewUnknownKey(scanCode)
	}
	switch action {
	case glfw.Press:
		e.buffer.Add(keyboard.NewPressedEvent(key))
	case glfw.Release:
		e.buffer.Add(keyboard.NewReleasedEvent(key))
	}
}

// Poll return next mapped event
func (e *KeyboardEvents) Poll() (keyboard.Event, bool) {
	return e.buffer.Poll()
}
