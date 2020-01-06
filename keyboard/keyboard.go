package keyboard

// EventsSource is a source of keyboard Events.
type EventsSource interface {
	// Fills output slice with events and returns number of events.
	Poll(output []Event) int
}

// Event describes what happened with the key. Whether it was pressed or released.
type Event struct {
	Type EventType
	Key  Key
}

// EventType describes the type of event
type EventType int

const (
	PRESSED  EventType = 1
	RELEASED EventType = 2
)

// Key tells what token has a given key. It may be UNKNOWN, then ScanCode should
// be used instead.
type Key struct {
	Token Token
	// Scancode is platform-specific but consistent over time, so keys will
	// have different scancodes depending on the platform but they are safe
	// to save to disk
	ScanCode int
}

// Token is platform-independent number identifying a key
type Token int

const (
	UNKNOWN Token = 0
	A       Token = 65
)

func New(poller EventsSource) *Keyboard {
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

// Pressed returns true if given key is currently pressed
func (k Keyboard) Pressed(keyToken Token) bool {
	return false
}
