package keyboard

import "fmt"

// EventsSource is a source of keyboard Events.
type EventsSource interface {
	// Fills output slice with events and returns number of events.
	// This function read new events and removes them from source.
	Poll(output []Event) int
}

// NewKey returns new instance of immutable Key.
func NewKey(token Token, scanCode int) Key {
	if token > 0 && (token < A || token > B) {
		panic(fmt.Sprintf("invalid token %v", token))
	}
	return Key{
		token:    token,
		scanCode: scanCode,
	}
}

// NewReleasedEvent returns new instance of Event when key was released.
func NewReleasedEvent(key Key) Event {
	return Event{}
}

// NewPressedEvent returns new instance of Event when key was pressed.
func NewPressedEvent(key Key) Event {
	return Event{
		typ: pressed,
		key: key,
	}
}

// Event describes what happened with the key. Whether it was pressed or released.
type Event struct {
	typ eventType
	key Key
}

// eventType is used because using polymorphism means heap allocation and we don't
// want to generate garbage (really? StarCraft e-sport players can perform up to
// 300 APM, which means 600 event objects per minute - maybe it is not that much).
type eventType byte

const (
	pressed  eventType = 1
	released eventType = 2
)

// Key contains numbers identifying the key.
type Key struct {
	token    Token
	scanCode int
}

// ScanCode returns the platform-specific code.
func (k Key) ScanCode() int {
	return k.scanCode
}

// Token return platform-independent mapping
func (k Key) Token() Token {
	return k.token
}

// Token is platform-independent mapping identifying the key. It may be
// Unknown, then ScanCode should be used instead.
type Token uint

const (
	// Unknown means that there is no known mapping for given key
	Unknown Token = 0
	// A token
	A Token = 65
	// B token
	B Token = 66
)

// New creates Keyboard instance.
func New(source EventsSource) *Keyboard {
	if source == nil {
		panic("nil EventsSource")
	}
	return &Keyboard{
		source:      source,
		events:      make([]Event, 32), // TODO magic number
		keysPressed: make(map[Token]bool),
	}
}

// Keyboard provides a read-only information about the current state of the
// keyboard, such as what keys are currently pressed.
type Keyboard struct {
	source      EventsSource
	events      []Event
	keysPressed map[Token]bool
}

// Update updates the state of the keyboard by polling events queued since last
// time the function was executed.
func (k *Keyboard) Update() {
	k.source.Poll(k.events)
	for _, event := range k.events {
		if event.typ == pressed {
			k.keysPressed[event.key.token] = true
		}
	}
}

// Pressed returns true if given key is currently pressed.
func (k *Keyboard) Pressed(keyToken Token) bool {
	return k.keysPressed[keyToken]
}
