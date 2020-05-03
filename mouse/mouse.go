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
	button Button
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
		source:      source,
		pressed:     map[Button]struct{}{},
		justPressed: make(map[Button]bool),
	}
}

type Mouse struct {
	source      EventSource
	pressed     map[Button]struct{}
	justPressed map[Button]bool
}

func (m *Mouse) Update() {
	m.clearJustPressed()
	for {
		event, ok := m.source.PollMouseEvent()
		if !ok {
			return
		}
		switch event.typ {
		case pressed:
			m.pressed[event.button] = struct{}{}
			m.justPressed[event.button] = true
		case released:
			delete(m.pressed, event.button)
		case moved:
			fmt.Println("moved", event.position.pixelPosX, event.position.pixelPosY, event.position.subpixelPosX, event.position.subpixelPosY, event.position.insideWindow)
		case scrolled:
			fmt.Println("scrolled", event.scrollX, event.scrollY)
		}
	}
}

func (m *Mouse) clearJustPressed() {
	for key := range m.justPressed {
		delete(m.justPressed, key)
	}
}

func (m *Mouse) Pressed(button Button) bool {
	_, found := m.pressed[button]
	return found
}

func (m *Mouse) PressedButtons() []Button {
	var pressedButtons []Button
	for key := range m.pressed {
		pressedButtons = append(pressedButtons, key)
	}
	return pressedButtons
}

func (m *Mouse) JustPressed(button Button) bool {
	return m.justPressed[button]
}

func (m *Mouse) JustReleased(a Button) bool {
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

type Button int

const (
	Left    Button = 1
	Right   Button = 2
	Middle  Button = 3
	Button4 Button = 4
	Button5 Button = 5
	Button6 Button = 6
	Button7 Button = 7
	Button8 Button = 8
)

func NewReleasedEvent(button Button) Event {
	return Event{
		typ:    released,
		button: button,
	}
}

func NewPressedEvent(button Button) Event {
	return Event{
		typ:    pressed,
		button: button,
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
			insideWindow: insideWindow,
		},
	}
}
