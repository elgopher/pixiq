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
	if len(e.events) == e.writeIndex {
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
	if len(e.events) > 0 && e.writeIndex > e.readIndex {
		event := e.events[e.readIndex]
		e.readIndex++
		return event, true
	}
	e.readIndex = 0
	e.writeIndex = 0
	return keyboard.EmptyEvent, false
}
