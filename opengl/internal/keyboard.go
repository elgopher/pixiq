package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/jacekolszak/pixiq/keyboard"
)

// KeyboardEvents maps GLFW events to keyboard.Event. Mapped events can be
// polled using keyboard.EventSource interface. This implementation if efficient
// because it is reusing event slice.
type KeyboardEvents struct {
	events     []keyboard.Event
	readIndex  int
	writeIndex int
}

// NewKeyboardEvents returns KeyboardEvents of given initial size. It will
// be expanded if necessary. Will panic if initial size smaller than 1.
func NewKeyboardEvents(initialSize int) *KeyboardEvents {
	if initialSize < 1 {
		panic("initial size was too small")
	}
	return &KeyboardEvents{
		events: make([]keyboard.Event, initialSize),
	}
}

// OnKeyCallback passes GLFW key event
func (e *KeyboardEvents) OnKeyCallback(_ *glfw.Window, glfwKey glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
	key, ok := keymap[glfwKey]
	if !ok {
		key = keyboard.NewUnknownKey(scanCode)
	}
	switch action {
	case glfw.Press:
		e.append(keyboard.NewPressedEvent(key))
	case glfw.Release:
		e.append(keyboard.NewReleasedEvent(key))
	}
}

func (e *KeyboardEvents) append(event keyboard.Event) {
	noSpace := len(e.events) == e.writeIndex
	if noSpace {
		e.expand()
	}
	e.events[e.writeIndex] = event
	e.writeIndex++
}

func (e *KeyboardEvents) expand() {
	largerEvents := make([]keyboard.Event, len(e.events)*2)
	copy(largerEvents, e.events)
	e.events = largerEvents
}

// Poll return next mapped event
func (e *KeyboardEvents) Poll() (keyboard.Event, bool) {
	somethingToRead := e.writeIndex > e.readIndex
	if somethingToRead {
		event := e.events[e.readIndex]
		e.readIndex++
		return event, true
	}
	e.Clear()
	return keyboard.EmptyEvent, false
}

// Clear effectively clears all collected events
func (e *KeyboardEvents) Clear() {
	e.readIndex = 0
	e.writeIndex = 0
}