package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/keyboard"
)

// KeyboardEvents maps GLFW events to keyboard.Event. Mapped events can be
// polled using keyboard.EventSource interface.
type KeyboardEvents struct {
	events []keyboard.Event
}

// OnKeyCallback passes GLFW key event
func (e *KeyboardEvents) OnKeyCallback(_ *glfw.Window, glfwKey glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
	key, ok := keymap[glfwKey]
	if !ok {
		key = keyboard.NewUnknownKey(scanCode)
	}
	switch action {
	case glfw.Press:
		e.events = append(e.events, keyboard.NewPressedEvent(key))
	case glfw.Release:
		e.events = append(e.events, keyboard.NewReleasedEvent(key))
	}
}

// Poll return next mapped event
func (e *KeyboardEvents) Poll() (keyboard.Event, bool) {
	if len(e.events) > 0 {
		event := e.events[0]
		e.events = e.events[1:]
		return event, true
	}
	return keyboard.EmptyEvent, false
}
