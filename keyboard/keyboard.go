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

// Token is platform-independent mapping identifying the key. It may be
// Unknown, then ScanCode should be used instead.
type Token uint

// Rune return the character assigned to key (capital letter)
func (t Token) Rune() rune {
	return rune(t)
}

var (
	// Space Key
	Space = newKey(' ')
	// Apostrophe Key
	Apostrophe = newKey('\'')
	// Comma Key
	Comma = newKey(',')
	// Minus Key
	Minus = newKey('-')
	// Period Key
	Period = newKey('.')
	// Slash Key
	Slash = newKey('/')
	// Zero Key
	Zero = newKey('0')
	// One Key
	One = newKey('1')
	// Two Key
	Two = newKey('2')
	// Three Key
	Three = newKey('3')
	// Four Key
	Four = newKey('4')
	// Five Key
	Five = newKey('5')
	// Six Key
	Six = newKey('6')
	// Seven Key
	Seven = newKey('7')
	// Eight Key
	Eight = newKey('8')
	// Nine Key
	Nine = newKey('9')
	// Semicolon Key
	Semicolon = newKey(';')
	// Equal Key
	Equal = newKey('=')
	// A key
	A = newKey('A')
	// B key
	B = newKey('B')
	// C Key
	C = newKey('C')
	// D Key
	D = newKey('D')
	// E Key
	E = newKey('E')
	// F Key
	F = newKey('F')
	// G Key
	G = newKey('G')
	// H Key
	H = newKey('H')
	// I Key
	I = newKey('I')
	// J Key
	J = newKey('J')
	// K Key
	K = newKey('K')
	// L Key
	L = newKey('L')
	// M Key
	M = newKey('M')
	// N Key
	N = newKey('N')
	// O Key
	O = newKey('O')
	// P Key
	P = newKey('P')
	// Q Key
	Q = newKey('Q')
	// R Key
	R = newKey('R')
	// S Key
	S = newKey('S')
	// T Key
	T = newKey('T')
	// U Key
	U = newKey('U')
	// V Key
	V = newKey('V')
	// W Key
	W = newKey('W')
	// X Key
	X = newKey('X')
	// Y Key
	Y = newKey('Y')
	// Z Key
	Z = newKey('Z')
	// LeftBracket Key
	LeftBracket = newKey('[')
	// Backslash Key
	Backslash = newKey('\\')
	// RightBracket Key
	RightBracket = newKey(']')
	// GraveAccent Key
	GraveAccent = newKey('`')
	// TODO keyworld 1 and 2? I think they should not be added.
	// Escape Key
	Escape = newKey(256)
	// Enter Key
	Enter = newKey(257)
	// Tab Key
	Tab = newKey(258)
	// Backspace Key
	Backspace = newKey(259)
	// Insert Key
	Insert = newKey(260)
	// Delete Key
	Delete = newKey(261)
	// Right Key
	Right = newKey(262)
	// Left Key
	Left = newKey(263)
	// Down Key
	Down = newKey(264)
	// Up Key
	Up = newKey(265)
	// PageUp Key
	PageUp = newKey(266)
	// PageDown Key
	PageDown = newKey(267)
	// Home Key
	Home = newKey(268)
	// End Key
	End = newKey(269)
	// CapsLock Key
	CapsLock = newKey(280)
	// ScrollLock Key
	ScrollLock = newKey(281)
	// NumLock Key
	NumLock = newKey(282)
	// PrintScreen Key
	PrintScreen = newKey(283) // TODO I think it does not work on Linux
	// Pause Key
	Pause = newKey(284)
	// F1 Key
	F1 = newKey(290)
	// F2 Key
	F2 = newKey(291)
	// F3 Key
	F3 = newKey(292)
	// F4 Key
	F4 = newKey(293)
	// F5 Key
	F5 = newKey(294)
	// F6 Key
	F6 = newKey(295)
	// F7 Key
	F7 = newKey(296)
	// F8 Key
	F8 = newKey(297)
	// F9 Key
	F9 = newKey(298)
	// F10 Key
	F10 = newKey(299)
	// F11 Key
	F11 = newKey(300)
	// F12 Key
	F12 = newKey(301)
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
