package mouse

import "fmt"

type EventSource interface {
	PollMouseEvent() (Event, bool)
}

// EmptyEvent should be returned by EventSource when it does not have more events.
var EmptyEvent = Event{}

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
		source:       source,
		pressed:      map[Button]struct{}{},
		justPressed:  make(map[Button]bool),
		justReleased: make(map[Button]bool),
		position: Position{
			insideWindow: true,
		},
	}
}

type Mouse struct {
	source         EventSource
	pressed        map[Button]struct{}
	justPressed    map[Button]bool
	justReleased   map[Button]bool
	position       Position
	positionChange PositionChange
}

func (m *Mouse) Update() {
	m.clearJustPressed()
	m.clearJustReleased()
	lastPosition := m.position
	defer func() {
		if lastPosition != m.position {
			windowLeft := false
			if !m.position.insideWindow && lastPosition.insideWindow {
				windowLeft = true
			}
			windowEntered := false
			if m.position.insideWindow && !lastPosition.insideWindow {
				windowEntered = true
			}
			m.positionChange = PositionChange{
				x:             m.position.x - lastPosition.x,
				y:             m.position.y - lastPosition.y,
				realX:         m.position.realX - lastPosition.realX,
				realY:         m.position.realY - lastPosition.realY,
				windowLeft:    windowLeft,
				windowEntered: windowEntered,
			}
		} else {
			m.positionChange = PositionChange{}
		}
	}()
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
			m.justReleased[event.button] = true
		case moved:
			m.position = event.position
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

func (m *Mouse) clearJustReleased() {
	for key := range m.justReleased {
		delete(m.justReleased, key)
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

func (m *Mouse) JustReleased(button Button) bool {
	return m.justReleased[button]
}

func (m *Mouse) Position() Position {
	return m.position
}

func (m *Mouse) PositionChange() PositionChange {
	return m.positionChange
}

func (m *Mouse) PositionChanged() bool {
	return m.positionChange != PositionChange{}
}

type Position struct {
	x, y         int
	realX, realY float64
	insideWindow bool
}

// X returns the pixel position
func (p Position) X() int {
	return p.x
}

func (p Position) Y() int {
	return p.y
}

func (p Position) RealX() float64 {
	return p.realX
}

func (p Position) RealY() float64 {
	return p.realY
}

func (p Position) InsideWindow() bool {
	return p.insideWindow
}

type PositionChange struct {
	x, y          int
	realX, realY  float64
	windowEntered bool
	windowLeft    bool
}

// X returns the pixel position
func (p PositionChange) X() int {
	return p.x
}

func (p PositionChange) Y() int {
	return p.y
}

func (p PositionChange) RealX() float64 {
	return p.realX
}

func (p PositionChange) RealY() float64 {
	return p.realY
}

func (p PositionChange) WindowEntered() bool {
	return p.windowEntered
}

func (p PositionChange) WindowLeft() bool {
	return p.windowLeft
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

func NewMovedEvent(posX, posY int, realPosX, realPosY float64, insideWindow bool) Event {
	return Event{
		typ: moved,
		position: Position{
			x:            posX,
			y:            posY,
			realX:        realPosX,
			realY:        realPosY,
			insideWindow: insideWindow,
		},
	}
}
