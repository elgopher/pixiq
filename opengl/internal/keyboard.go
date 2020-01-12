package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/keyboard"
)

// KeyboardEvents maps GLFW events to keyboard.Event. Mapped events can be
// polled using keyboard.EventSource interface.
type KeyboardEvents struct {
	queue *keyboard.EventQueue
}

// NewKeyboardEventsOfSize creates *KeyboardEvents of given size.
func NewKeyboardEvents(queue *keyboard.EventQueue) *KeyboardEvents {
	if queue == nil {
		panic("nil queue")
	}
	return &KeyboardEvents{queue: queue}
}

// OnKeyCallback passes GLFW key event
func (e *KeyboardEvents) OnKeyCallback(_ *glfw.Window, glfwKey glfw.Key, scanCode int, action glfw.Action, _ glfw.ModifierKey) {
	key, ok := keymap[glfwKey]
	if !ok {
		key = keyboard.NewUnknownKey(scanCode)
	}
	switch action {
	case glfw.Press:
		e.queue.Append(keyboard.NewPressedEvent(key))
	case glfw.Release:
		e.queue.Append(keyboard.NewReleasedEvent(key))
	}
}

// Poll return next mapped event
func (e *KeyboardEvents) Poll() (keyboard.Event, bool) {
	return e.queue.Poll()
}
