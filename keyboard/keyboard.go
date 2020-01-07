package keyboard

import "fmt"

// EventSource is a source of keyboard Events.
type EventSource interface {
	// Poll retrieves and removes next keyboard Event. If there are no more
	// events false is returned.
	Poll() (Event, bool)
}

// NewKey returns new instance of immutable Key. This construct may be used for
// creating keys after deserialization. Otherwise package variables should be used.
func NewKey(token Token) Key {
	if token < 65 || token > 66 {
		panic(fmt.Sprintf("invalid token %v", token))
	}
	return Key{
		token: token,
	}
}

// NewUnknownKey creates key with platform-specific scanCode.
func NewUnknownKey(scanCode int) Key {
	return Key{
		scanCode: scanCode,
	}
}

// Key identifies the pressed or released key.
type Key struct {
	token    Token
	scanCode int
}

// IsUnknown returns true if mapping has not been found. ScanCode should be used
// instead.
func (k Key) IsUnknown() bool {
	return k.token == 0
}

// ScanCode returns the platform-specific code.
func (k Key) ScanCode() int {
	return k.scanCode
}

// Token return platform-independent mapping.
func (k Key) Token() Token {
	return k.token
}

// Token is platform-independent mapping identifying the key. It may be
// Unknown, then ScanCode should be used instead.
type Token uint

// Rune return the character assigned to key (capital letter)
func (t Token) Rune() rune {
	return rune(t)
}

var (
	// A key
	A = NewKey(65)
	// B key
	B = NewKey(66)

	// EmptyEvent should be returned by EventSource when it does not have more events
	EmptyEvent = Event{}
)

// NewPressedEvent returns new instance of Event when key was pressed.
func NewPressedEvent(key Key) Event {
	return Event{
		typ: pressed,
		key: key,
	}
}

// NewReleasedEvent returns new instance of Event when key was released.
func NewReleasedEvent(key Key) Event {
	return Event{
		typ: released,
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

// New creates Keyboard instance.
func New(source EventSource) *Keyboard {
	if source == nil {
		panic("nil EventSource")
	}
	return &Keyboard{
		source:                source,
		keysPressedByToken:    make(map[Token]bool),
		keysPressedByScanCode: make(map[int]bool),
	}
}

// Keyboard provides a read-only information about the current state of the
// keyboard, such as what keys are currently pressed.
type Keyboard struct {
	source                EventSource
	keysPressedByToken    map[Token]bool
	keysPressedByScanCode map[int]bool
}

// Update updates the state of the keyboard by polling events queued since last
// time the function was executed.
func (k *Keyboard) Update() {
	for {
		event, ok := k.source.Poll()
		if !ok {
			return
		}
		switch event.typ {
		case pressed:
			if event.key.IsUnknown() {
				k.keysPressedByScanCode[event.key.scanCode] = true
			} else {
				k.keysPressedByToken[event.key.token] = true
			}
		case released:
			if event.key.IsUnknown() {
				k.keysPressedByScanCode[event.key.scanCode] = false
			} else {
				k.keysPressedByToken[event.key.token] = false
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
