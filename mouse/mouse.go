package mouse

type EventSource interface {
}

// Event describes what happened with the mouse button. Whether it was pressed or released.
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

func New(source EventSource) *Mouse {
	if source == nil {
		panic("nil EventSource")
	}
	return &Mouse{
		source: source,
	}
}

type Mouse struct {
	source EventSource
}

func (m *Mouse) Update() {

}

func (m *Mouse) Pressed(a Key) bool {
	return false
}

func (m *Mouse) Position() Position {
	return Position{}
}

type Position struct {
}

// X returns the pixel position
func (p Position) X() int {
	return 0
}

func (p Position) Y() int {
	return 0
}

// Xf is useful when zoom was used.
func (p Position) Xf() float32 {
	return 0
}

type Key struct {
	name string
}

func newKey(name string) Key {
	return Key{name: name}
}

var (
	Left  = newKey("Left")
	Right = newKey("Right")
)
