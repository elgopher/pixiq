package mouse

// EventSource is a source of mouse Events. On each Update() Mouse polls
// the EventSource by executing PollMouseEvent method multiple times - until PollMouseEvent()
// returns false. In other words Mouse#Update drains the EventSource.
type EventSource interface {
	// PollMouseEvent retrieves and removes next mouse Event. If there are no more
	// events false is returned.
	PollMouseEvent() (Event, bool)
}

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
	scroll         Scroll
}

func (m *Mouse) Update() {
	m.clearJustPressed()
	m.clearJustReleased()
	lastPosition := m.position
	m.scroll = Scroll{}
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
			m.scroll = Scroll{
				x: event.scrollX + m.scroll.x,
				y: event.scrollY + m.scroll.y,
			}
		}
	}
}

func (m *Mouse) clearJustPressed() {
	for button := range m.justPressed {
		delete(m.justPressed, button)
	}
}

func (m *Mouse) clearJustReleased() {
	for button := range m.justReleased {
		delete(m.justReleased, button)
	}
}

// Pressed returns true if given mouse button is currently pressed.
// If between two last mouse.Update calls the key was pressed and released
// then the this method returns false.
func (m *Mouse) Pressed(button Button) bool {
	_, found := m.pressed[button]
	return found
}

// PressedButtons returns a slice of all currently pressed buttons. It may be empty
// aka nil. This function can be used to get a button mapping for a given action
// in the game.
// If between two last mouse.Update calls the button was pressed and released
// then the button is not returned.
func (m *Mouse) PressedButtons() []Button {
	var pressedButtons []Button
	for button := range m.pressed {
		pressedButtons = append(pressedButtons, button)
	}
	return pressedButtons
}

// JustPressed returns true if the button was pressed between two last mouse.Update
// calls. If it was pressed and released at the same time between these calls
// this method return true.
func (m *Mouse) JustPressed(button Button) bool {
	return m.justPressed[button]
}

// JustReleased returns true if the button was released between two last mouse.Update
// calls. If it was released and pressed at the same time between these calls
// this method return true.
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

func (m *Mouse) Scroll() Scroll {
	return m.scroll
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

type Scroll struct {
	x, y float64
}

func (s Scroll) X() float64 {
	return s.x
}

func (s Scroll) Y() float64 {
	return s.y
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
