package keyboard

// EventsSource is a source of keyboard Events.
type EventsSource interface {
	// Fills output slice with events and returns number of events.
	Poll(output []Event) int
}

// UnknownKey returns instance of unknown Key.
// Scancode is platform-specific but consistent over time.
func UnknownKey(scanCode int) Key {
	return Key{}
}

// NewKey returns new instance of immutable Key.
func NewKey(token Token, scanCode int) Key {
	return Key{}
}

// NewReleasedEvent returns new instance of Event when key was released
func NewReleasedEvent(key Key) Event {
	return Event{
		typ: released,
		key: key,
	}
}

// NewReleasedEvent returns new instance of Event when key was pressed
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

// Token is platform-independent number identifying the key. It may be
// Unknown, then ScanCode should be used instead.
type Token int

const (
	// Unknown means that key cannot be mapped to well known Token.
	Unknown Token = 0
	// A is 65
	A Token = 65
)

// New creates Keyboard instance.
func New(source EventsSource) *Keyboard {
	return &Keyboard{}
}

// Keyboard provides a read-only information about the current state of the
// keyboard, such as what keys are currently pressed.
type Keyboard struct {
}

// Update updates the state of the keyboard by polling events queued since last time
// the function was executed.
func (k *Keyboard) Update() {

}

// Pressed returns true if given key is currently pressed.
func (k Keyboard) Pressed(keyToken Token) bool {
	return false
}
