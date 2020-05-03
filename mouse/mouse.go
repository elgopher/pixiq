package mouse

import "fmt"

type EventSource interface {
	PollMouseEvent() (Event, bool)
}

// EmptyEvent should be returned by EventSource when it does not have more events.
var EmptyEvent = Event{}

// Event describes what happened with the mouse button. Whether it was pressed or released.
type Event struct {
	typ eventType
	// Pressed/Released
	key Key
	// Moved
	position Position
	// Scroll
	scrollX, scrollY float64
}

type eventType byte

const (
	pressed eventType = iota
	released
	moved
	scrolled
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
	for {
		event, ok := m.source.PollMouseEvent()
		if !ok {
			return
		}
		switch event.typ {
		case pressed:
			fmt.Println("pressed")
		case released:
			fmt.Println("released")
		case moved:
			fmt.Println("moved", event.position.pixelPosX, event.position.pixelPosY, event.position.insideWindow)
		case scrolled:
			fmt.Println("scrolled")
		}
	}

}

func (m *Mouse) Pressed(a Key) bool {
	return false
}

func (m *Mouse) Position() Position {
	return Position{}
}

func (m *Mouse) PositionDelta() PositionDelta {
	return PositionDelta{}
}

type PositionDelta struct {
}

type Position struct {
	pixelPosX, pixelPosY       int
	subpixelPosX, subpixelPosY float64
	insideWindow               bool
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

type Key int

const (
	Left Key = iota
	Right
	Middle
)

func NewReleasedEvent(key Key) Event {
	return Event{
		typ: released,
		key: key,
	}
}

func NewPressedEvent(key Key) Event {
	return Event{
		typ: pressed,
		key: key,
	}
}

func NewScrolledEvent(x, y float64) Event {
	return Event{
		typ:     scrolled,
		scrollX: x,
		scrollY: y,
	}
}

func NewMovedEvent(pixelPosX, pixelPosY int, subpixelPosX, subpixelPosY float64, insideWindow bool) Event {
	return Event{
		typ: moved,
		position: Position{
			pixelPosX:    pixelPosX,
			pixelPosY:    pixelPosY,
			subpixelPosX: subpixelPosX,
			subpixelPosY: subpixelPosY,
			insideWindow: true,
		},
	}
}
