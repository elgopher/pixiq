package keyboard

// EventSource is a source of keyboard Events.
type EventSource interface {
	// Poll retrieves and removes next keyboard Event. If there are no more
	// events false is returned.
	Poll() (Event, bool)
}

// TODO Add serialization deserialization of key
func newKey(token Token) Key {
	return Key{token: token}
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
	return k.token == ""
}

// ScanCode returns the platform-specific code.
func (k Key) ScanCode() int {
	return k.scanCode
}

// Token returns platform-independent token.
func (k Key) Token() Token {
	return k.token
}

func (k Key) setPressed(keyboard *Keyboard, value bool) {
	if k.IsUnknown() {
		keyboard.keysPressedByScanCode[k.scanCode] = value
	} else {
		keyboard.keysPressedByToken[k.token] = value
	}
}

func (k Key) pressed(keyboard *Keyboard) bool {
	if k.IsUnknown() {
		return keyboard.keysPressedByScanCode[k.scanCode]
	}
	return keyboard.keysPressedByToken[k.token]
}

// Token is a string representation of the key. It the key is
// unknown then token is empty.
type Token string

var (
	// Space Key " "
	Space = newKey(" ")
	// Apostrophe Key "'"
	Apostrophe = newKey("'")
	// Comma Key ","
	Comma = newKey(",")
	// Minus Key "-"
	Minus = newKey("-")
	// Period Key "."
	Period = newKey(".")
	// Slash Key "/"
	Slash = newKey("/")
	// Zero Key "0"
	Zero = newKey("0")
	// One Key "1"
	One = newKey("1")
	// Two Key "2"
	Two = newKey("2")
	// Three Key "3"
	Three = newKey("3")
	// Four Key "4"
	Four = newKey("4")
	// Five Key "5"
	Five = newKey("5")
	// Six Key "6"
	Six = newKey("6")
	// Seven Key "7"
	Seven = newKey("7")
	// Eight Key "8"
	Eight = newKey("8")
	// Nine Key "9"
	Nine = newKey("9")
	// Semicolon Key ";"
	Semicolon = newKey(";")
	// Equal Key "="
	Equal = newKey("=")
	// A key
	A = newKey("A")
	// B key
	B = newKey("B")
	// C Key
	C = newKey("C")
	// D Key
	D = newKey("D")
	// E Key
	E = newKey("E")
	// F Key
	F = newKey("F")
	// G Key
	G = newKey("G")
	// H Key
	H = newKey("H")
	// I Key
	I = newKey("I")
	// J Key
	J = newKey("J")
	// K Key
	K = newKey("K")
	// L Key
	L = newKey("L")
	// M Key
	M = newKey("M")
	// N Key
	N = newKey("N")
	// O Key
	O = newKey("O")
	// P Key
	P = newKey("P")
	// Q Key
	Q = newKey("Q")
	// R Key
	R = newKey("R")
	// S Key
	S = newKey("S")
	// T Key
	T = newKey("T")
	// U Key
	U = newKey("U")
	// V Key
	V = newKey("V")
	// W Key
	W = newKey("W")
	// X Key
	X = newKey("X")
	// Y Key
	Y = newKey("Y")
	// Z Key
	Z = newKey("Z")
	// LeftBracket Key "["
	LeftBracket = newKey("[")
	// Backslash Key "\"
	Backslash = newKey("\\")
	// RightBracket Key "]"
	RightBracket = newKey("]")
	// GraveAccent Key "`"
	GraveAccent = newKey("`")
	// TODO keyworld 1 and 2? I think they should not be added.
	// Esc Key
	Esc = newKey("Esc")
	// Enter Key
	Enter = newKey("Enter")
	// Tab Key
	Tab = newKey("Tab")
	// Backspace Key
	Backspace = newKey("Backspace")
	// Insert Key
	Insert = newKey("Insert")
	// Delete Key
	Delete = newKey("Delete")
	// Right Key
	Right = newKey("Right")
	// Left Key
	Left = newKey("Left")
	// Down Key
	Down = newKey("Down")
	// Up Key
	Up = newKey("Up")
	// PageUp Key
	PageUp = newKey("PageUp")
	// PageDown Key
	PageDown = newKey("PageDown")
	// Home Key
	Home = newKey("Home")
	// End Key
	End = newKey("End")
	// CapsLock Key
	CapsLock = newKey("CapsLock")
	// ScrollLock Key
	ScrollLock = newKey("ScrollLock")
	// NumLock Key
	NumLock = newKey("NumLock")
	// PrintScreen Key
	PrintScreen = newKey("PrintScreen")
	// Pause Key
	Pause = newKey("Pause")
	// F1 Key
	F1 = newKey("F1")
	// F2 Key
	F2 = newKey("F2")
	// F3 Key
	F3 = newKey("F3")
	// F4 Key
	F4 = newKey("F4")
	// F5 Key
	F5 = newKey("F5")
	// F6 Key
	F6 = newKey("F6")
	// F7 Key
	F7 = newKey("F7")
	// F8 Key
	F8 = newKey("F8")
	// F9 Key
	F9 = newKey("F9")
	// F10 Key
	F10 = newKey("F10")
	// F11 Key
	F11 = newKey("F11")
	// F12 Key
	F12 = newKey("F12")
	// F13 Key
	F13 = newKey("F13")
	// F14 Key
	F14 = newKey("F14")
	// F15 Key
	F15 = newKey("F15")
	// F16 Key
	F16 = newKey("F16")
	// F17 Key
	F17 = newKey("F17")
	// F18 Key
	F18 = newKey("F18")
	// F19 Key
	F19 = newKey("F19")
	// F20 Key
	F20 = newKey("F20")
	// F21 Key
	F21 = newKey("F21")
	// F22 Key
	F22 = newKey("F22")
	// F23 Key
	F23 = newKey("F23")
	// F24 Key
	F24 = newKey("F24")
	// F25 Key
	F25 = newKey("F25")
	// Keypad0 Key
	Keypad0 = newKey("Keypad 0")
	// Keypad1 Key
	Keypad1 = newKey("Keypad 1")
	// Keypad2 Key
	Keypad2 = newKey("Keypad 2")
	// Keypad3 Key
	Keypad3 = newKey("Keypad 3")
	// Keypad4 Key
	Keypad4 = newKey("Keypad 4")
	// Keypad5 Key
	Keypad5 = newKey("Keypad 5")
	// Keypad6 Key
	Keypad6 = newKey("Keypad 6")
	// Keypad7 Key
	Keypad7 = newKey("Keypad 7")
	// Keypad8 Key
	Keypad8 = newKey("Keypad 8")
	// Keypad9 Key
	Keypad9 = newKey("Keypad 9")
	// TODO

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
// want to generate garbage.
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
			event.key.setPressed(k, true)
		case released:
			event.key.setPressed(k, false)
		}
	}
}

// Pressed returns true if given key is currently pressed.
func (k *Keyboard) Pressed(key Key) bool {
	return key.pressed(k)
}
