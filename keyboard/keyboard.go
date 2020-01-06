package keyboard

import "fmt"

// EventsSource is a source of keyboard Events.
type EventsSource interface {
	// Fills output slice with events and returns number of events.
	// This function read new events and removes them from source.
	Poll(output []Event) int
}

// NewKey returns new instance of immutable Key. This construct may be used for
// creating keys after deserialization.
func NewKey(token Token) Key {
	if token < 65 || token > 66 {
		panic(fmt.Sprintf("invalid token %v", token))
	}
	return Key{
		token: token,
	}
}

func NewUnknownKey(scanCode int) Key {
	return Key{
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

// Key identifies the pressed or release key.
type Key struct {
	token    Token
	scanCode int
}

// IsUnknown returns true if mapping has not been found. ScanCode should be used
// instead.
func (k Key) IsUnknown() bool {
	return k.token == unknown
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

const unknown Token = 0

var (
	// A key
	A = NewKey(65)
	// B key
	B = NewKey(66)
)

// New creates Keyboard instance.
func New(source EventsSource) *Keyboard {
	if source == nil {
		panic("nil EventsSource")
	}
	return &Keyboard{
		source:                source,
		events:                make([]Event, 32), // TODO magic number
		keysPressedByToken:    make(map[Token]bool),
		keysPressedByScanCode: make(map[int]bool),
	}
}

// Keyboard provides a read-only information about the current state of the
// keyboard, such as what keys are currently pressed.
type Keyboard struct {
	source                EventsSource
	events                []Event
	keysPressedByToken    map[Token]bool
	keysPressedByScanCode map[int]bool
}

// Update updates the state of the keyboard by polling events queued since last
// time the function was executed.
func (k *Keyboard) Update() {
	k.source.Poll(k.events)
	for _, event := range k.events {
		if event.typ == pressed {
			if event.key.IsUnknown() {
				k.keysPressedByScanCode[event.key.scanCode] = true
			} else {
				k.keysPressedByToken[event.key.token] = true
			}
		}
	}
}

// Pressed returns true if given key is currently pressed.
func (k *Keyboard) Pressed(key Key) bool {
	if key.IsUnknown() {
		return k.keysPressedByScanCode[key.scanCode]
	}
	return k.keysPressedByToken[key.token]
}
