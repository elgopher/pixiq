package mouse_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pixiq/mouse"
)

func TestNew(t *testing.T) {
	t.Run("should panic when source is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			mouse.New(nil)
		})
	})
	t.Run("should create a mouse instance", func(t *testing.T) {
		source := &fakeEventSource{}
		// when
		mouseState := mouse.New(source)
		// then
		assert.NotNil(t, mouseState)
	})
}

func TestMouse_Pressed(t *testing.T) {
	t.Run("before Update was called, Pressed returns false for all buttons", func(t *testing.T) {
		tests := []mouse.Button{mouse.Left, mouse.Right, mouse.Middle, mouse.Button4, mouse.Button5, mouse.Button6, mouse.Button7, mouse.Button8}
		for _, button := range tests {
			testName := fmt.Sprintf("for button: %v", button)
			t.Run(testName, func(t *testing.T) {
				var (
					event      = mouse.NewPressedEvent(mouse.Left)
					source     = newFakeEventSource(event)
					mouseState = mouse.New(source)
				)
				// when
				pressed := mouseState.Pressed(button)
				// then
				assert.False(t, pressed)
			})
		}
	})
	t.Run("after Update was called", func(t *testing.T) {
		var (
			leftPressed   = mouse.NewPressedEvent(mouse.Left)
			leftReleased  = mouse.NewReleasedEvent(mouse.Left)
			rightPressed  = mouse.NewPressedEvent(mouse.Right)
			rightReleased = mouse.NewReleasedEvent(mouse.Right)
		)
		tests := map[string]struct {
			source             mouse.EventSource
			expectedPressed    []mouse.Button
			expectedNotPressed []mouse.Button
		}{
			"one PressedEvent for Left": {
				source:             newFakeEventSource(leftPressed),
				expectedPressed:    []mouse.Button{mouse.Left},
				expectedNotPressed: []mouse.Button{mouse.Right},
			},
			"two PressedEvents for Right and Left": {
				source:          newFakeEventSource(rightPressed, leftPressed),
				expectedPressed: []mouse.Button{mouse.Left, mouse.Right},
			},
			"two PressedEvents for Left and Right": {
				source:          newFakeEventSource(leftPressed, rightPressed),
				expectedPressed: []mouse.Button{mouse.Left, mouse.Right},
			},
			"one PressedEvent; one ReleasedEvent for Left": {
				source:             newFakeEventSource(leftPressed, leftReleased),
				expectedNotPressed: []mouse.Button{mouse.Left},
			},
			"one PressedEvent for Left; one ReleasedEvent for Right": {
				source:             newFakeEventSource(leftPressed, rightReleased),
				expectedPressed:    []mouse.Button{mouse.Left},
				expectedNotPressed: []mouse.Button{mouse.Right},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				// when
				mouseState.Update()
				// then
				for _, expectedPressedButton := range test.expectedPressed {
					assert.True(t, mouseState.Pressed(expectedPressedButton))
				}
				for _, expectedNotPressedButton := range test.expectedNotPressed {
					assert.False(t, mouseState.Pressed(expectedNotPressedButton))
				}
			})
		}
	})
}

func TestMouse_PressedButtons(t *testing.T) {
	var (
		leftPressed  = mouse.NewPressedEvent(mouse.Left)
		leftReleased = mouse.NewReleasedEvent(mouse.Left)
		rightPressed = mouse.NewPressedEvent(mouse.Right)
	)
	t.Run("before Update pressed buttons are empty", func(t *testing.T) {
		source := newFakeEventSource(leftPressed)
		mouseState := mouse.New(source)
		// when
		pressed := mouseState.PressedButtons()
		// then
		assert.Empty(t, pressed)
	})
	t.Run("after Update", func(t *testing.T) {
		tests := map[string]struct {
			source          mouse.EventSource
			expectedPressed []mouse.Button
		}{
			"one PressedEvent for Left": {
				source:          newFakeEventSource(leftPressed),
				expectedPressed: []mouse.Button{mouse.Left},
			},
			"one PressedEvent for Right": {
				source:          newFakeEventSource(rightPressed),
				expectedPressed: []mouse.Button{mouse.Right},
			},
			"one PressedEvent for Left, one ReleaseEvent for Left": {
				source:          newFakeEventSource(leftPressed, leftReleased),
				expectedPressed: nil,
			},
			"Left pressed, then released and pressed again": {
				source:          newFakeEventSource(leftPressed, leftReleased, leftPressed),
				expectedPressed: []mouse.Button{mouse.Left},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				mouseState.Update()
				// when
				pressed := mouseState.PressedButtons()
				// then
				assert.Equal(t, test.expectedPressed, pressed)
			})
		}
	})
	t.Run("after second update", func(t *testing.T) {
		t.Run("when Left was pressed before first update", func(t *testing.T) {
			source := newFakeEventSource(leftPressed)
			mouseState := mouse.New(source)
			mouseState.Update()
			mouseState.Update()
			// when
			pressed := mouseState.PressedButtons()
			// then
			assert.Equal(t, []mouse.Button{mouse.Left}, pressed)
		})
		t.Run("when Left was pressed before first update, then released before second one", func(t *testing.T) {
			source := newFakeEventSource(leftPressed)
			mouseState := mouse.New(source)
			mouseState.Update()
			source.events = append(source.events, leftReleased)
			mouseState.Update()
			// when
			pressed := mouseState.PressedButtons()
			// then
			assert.Empty(t, pressed)
		})
	})
}

func TestJustPressed(t *testing.T) {
	var (
		leftPressed  = mouse.NewPressedEvent(mouse.Left)
		leftReleased = mouse.NewReleasedEvent(mouse.Left)
		rightPressed = mouse.NewPressedEvent(mouse.Right)
	)

	t.Run("before update should return false", func(t *testing.T) {
		tests := map[string]struct {
			source mouse.EventSource
		}{
			"for no events": {
				source: newFakeEventSource(),
			},
			"when Left was pressed": {
				source: newFakeEventSource(leftPressed),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				// when
				justPressed := mouseState.JustPressed(mouse.Left)
				// then
				assert.False(t, justPressed)
			})
		}
	})

	t.Run("after first update", func(t *testing.T) {
		tests := map[string]struct {
			source              mouse.EventSource
			expectedJustPressed bool
		}{
			"for no events": {
				source:              newFakeEventSource(),
				expectedJustPressed: false,
			},
			"when Left has been pressed": {
				source:              newFakeEventSource(leftPressed),
				expectedJustPressed: true,
			},
			"when Right has been pressed": {
				source:              newFakeEventSource(rightPressed),
				expectedJustPressed: false,
			},
			"when Left has been released": {
				source:              newFakeEventSource(leftReleased),
				expectedJustPressed: false,
			},
			"when Left has been pressed and released": {
				source:              newFakeEventSource(leftPressed, leftReleased),
				expectedJustPressed: true,
			},
			"when Left has been released and pressed": {
				source:              newFakeEventSource(leftReleased, leftPressed),
				expectedJustPressed: true,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				mouseState.Update()
				// when
				justPressed := mouseState.JustPressed(mouse.Left)
				// then
				assert.Equal(t, test.expectedJustPressed, justPressed)
			})
		}
	})

	t.Run("should return false after second update", func(t *testing.T) {
		source := newFakeEventSource(leftPressed)
		mouseState := mouse.New(source)
		mouseState.Update()
		mouseState.Update()
		// when
		pressed := mouseState.JustPressed(mouse.Left)
		// then
		assert.False(t, pressed)
	})
}

func TestJustReleased(t *testing.T) {
	var (
		leftReleased  = mouse.NewReleasedEvent(mouse.Left)
		leftPressed   = mouse.NewPressedEvent(mouse.Left)
		rightReleased = mouse.NewReleasedEvent(mouse.Right)
	)

	t.Run("before update should return false", func(t *testing.T) {
		tests := map[string]struct {
			source mouse.EventSource
		}{
			"for no events": {
				source: newFakeEventSource(),
			},
			"when Left was released": {
				source: newFakeEventSource(leftReleased),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				// when
				justReleased := mouseState.JustReleased(mouse.Left)
				// then
				assert.False(t, justReleased)
			})
		}
	})

	t.Run("after first update", func(t *testing.T) {
		tests := map[string]struct {
			source               mouse.EventSource
			expectedJustReleased bool
		}{
			"for no events": {
				source:               newFakeEventSource(),
				expectedJustReleased: false,
			},
			"when Left was released": {
				source:               newFakeEventSource(leftReleased),
				expectedJustReleased: true,
			},
			"when Right was released": {
				source:               newFakeEventSource(rightReleased),
				expectedJustReleased: false,
			},
			"when Left was pressed": {
				source:               newFakeEventSource(leftPressed),
				expectedJustReleased: false,
			},
			"when Left was released and pressed": {
				source:               newFakeEventSource(leftReleased, leftPressed),
				expectedJustReleased: true,
			},
			"when Left was pressed and released": {
				source:               newFakeEventSource(leftPressed, leftReleased),
				expectedJustReleased: true,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				mouseState.Update()
				// when
				justReleased := mouseState.JustReleased(mouse.Left)
				// then
				assert.Equal(t, test.expectedJustReleased, justReleased)
			})
		}
	})

	t.Run("should return false after second update", func(t *testing.T) {
		source := newFakeEventSource(leftReleased)
		mouseState := mouse.New(source)
		mouseState.Update()
		mouseState.Update()
		// when
		released := mouseState.JustReleased(mouse.Left)
		// then
		assert.False(t, released)
	})
}

func TestMouse_Update(t *testing.T) {
	tests := map[string]int{
		"1 event":     1,
		"2 events":    2,
		"1000 events": 1000,
	}
	for name, numberOfEvents := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("should drain EventSource", func(t *testing.T) {
				source := newFakeEventSourceWith(numberOfEvents)
				mouseState := mouse.New(source)
				// when
				mouseState.Update()
				// then
				assert.Empty(t, source.events)
			})

			t.Run("should drain EventSource after second Update()", func(t *testing.T) {
				source := newFakeEventSourceWith(numberOfEvents)
				mouseState := mouse.New(source)
				mouseState.Update()
				// when
				mouseState.Update()
				// then
				assert.Empty(t, source.events)
			})
		})
	}
}

type expectedPosition struct {
	x, y         int
	realX, realY float64
	insideWindow bool
}

type expectedPositionChange struct {
	x, y                      int
	realX, realY              float64
	windowEntered, windowLeft bool
	hasChanged                bool
}

func TestMouse_Position(t *testing.T) {
	t.Run("before Update was called, Position returns 0, 0", func(t *testing.T) {
		source := newFakeEventSource()
		mouseState := mouse.New(source)
		// when
		position := mouseState.Position()
		// then
		assert.Equal(t, 0, position.X())
		assert.Equal(t, 0, position.Y())
		assert.Equal(t, 0.0, position.RealX())
		assert.Equal(t, 0.0, position.RealY())
		assert.True(t, position.InsideWindow())
	})
	t.Run("after Update was called", func(t *testing.T) {
		// when
		tests := map[string]struct {
			source           mouse.EventSource
			expectedPosition expectedPosition
		}{
			"no events": {
				source: newFakeEventSource(),
				expectedPosition: expectedPosition{
					insideWindow: true,
				},
			},
			"one event": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(0, 0, 0.0, 0.0, true)),
				expectedPosition: expectedPosition{
					insideWindow: true,
				},
			},
			"one event with non default values": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(1, 2, 1.0, 2.0, false)),
				expectedPosition: expectedPosition{
					x:            1,
					y:            2,
					realX:        1.0,
					realY:        2.0,
					insideWindow: false,
				},
			},
			"one event with zoom": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(1, 2, 4, 8, true)),
				expectedPosition: expectedPosition{
					x:            1,
					y:            2,
					realX:        4.0,
					realY:        8.0,
					insideWindow: true,
				},
			},
			"one event with subpixels": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(1, 2, 1.5, 2.3, true)),
				expectedPosition: expectedPosition{
					x:            1,
					y:            2,
					realX:        1.5,
					realY:        2.3,
					insideWindow: true,
				},
			},
			"two events": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(2, 3, 2.0, 3.0, true),
					mouse.NewMovedEvent(1, 2, 1.0, 2.0, false),
				),
				expectedPosition: expectedPosition{
					x:            1,
					y:            2,
					realX:        1.0,
					realY:        2.0,
					insideWindow: false,
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				mouseState.Update()
				// when
				position := mouseState.Position()
				// then
				assert.Equal(t, test.expectedPosition.x, position.X())
				assert.Equal(t, test.expectedPosition.y, position.Y())
				// and
				assert.Equal(t, test.expectedPosition.realX, position.RealX())
				assert.Equal(t, test.expectedPosition.realY, position.RealY())
				// and
				assert.Equal(t, test.expectedPosition.insideWindow, position.InsideWindow())
			})
		}
	})
	t.Run("should return position change", func(t *testing.T) {
		// when
		tests := map[string]struct {
			source         mouse.EventSource
			expectedChange expectedPositionChange
		}{
			"no events": {
				source: newFakeEventSource(),
			},
			"one event": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(1, 2, 1.0, 2.0, false),
				),
				expectedChange: expectedPositionChange{
					x:          1,
					y:          2,
					realX:      1.0,
					realY:      2.0,
					windowLeft: true,
					hasChanged: true,
				},
			},
			"two events": {
				source: newFakeEventSource(
					mouse.NewMovedEvent(1, 2, 1.0, 2.0, false),
					mouse.NewMovedEvent(2, 3, 2.0, 3.0, false),
				),
				expectedChange: expectedPositionChange{
					x:          2,
					y:          3,
					realX:      2.0,
					realY:      3.0,
					windowLeft: true,
					hasChanged: true,
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				mouseState.Update()
				// when
				change := mouseState.PositionChange()
				// then
				assert.Equal(t, test.expectedChange.hasChanged, mouseState.PositionChanged())
				// and
				assert.Equal(t, test.expectedChange.x, change.X())
				assert.Equal(t, test.expectedChange.y, change.Y())
				// and
				assert.Equal(t, test.expectedChange.realX, change.RealX())
				assert.Equal(t, test.expectedChange.realY, change.RealY())
				// and
				assert.Equal(t, test.expectedChange.windowEntered, change.WindowEntered())
				assert.Equal(t, test.expectedChange.windowLeft, change.WindowLeft())
			})
		}
	})
	t.Run("should return position change after second update", func(t *testing.T) {
		// when
		tests := map[string]struct {
			newEvents      []mouse.Event
			expectedChange expectedPositionChange
		}{
			"no new events": {},
			"one new event": {
				newEvents: []mouse.Event{
					mouse.NewMovedEvent(2, 4, 2.5, 4.3, false),
				},
				expectedChange: expectedPositionChange{
					x:          1,
					y:          2,
					realX:      1.5,
					realY:      2.3,
					hasChanged: true,
				},
			},
			"one new event inside window": {
				newEvents: []mouse.Event{
					mouse.NewMovedEvent(1, 0, 1.0, 0, true),
				},
				expectedChange: expectedPositionChange{
					x:             0,
					y:             -2,
					realX:         0.0,
					realY:         -2.0,
					windowEntered: true,
					hasChanged:    true,
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				source := newFakeEventSource(
					mouse.NewMovedEvent(1, 2, 1.0, 2.0, false),
				)
				mouseState := mouse.New(source)
				mouseState.Update()
				source.events = test.newEvents
				mouseState.Update()
				// when
				change := mouseState.PositionChange()
				// then
				assert.Equal(t, test.expectedChange.hasChanged, mouseState.PositionChanged())
				// and
				assert.Equal(t, test.expectedChange.x, change.X())
				assert.Equal(t, test.expectedChange.y, change.Y())
				// and
				assert.Equal(t, test.expectedChange.realX, change.RealX())
				assert.Equal(t, test.expectedChange.realY, change.RealY())
				// and
				assert.Equal(t, test.expectedChange.windowEntered, change.WindowEntered())
				assert.Equal(t, test.expectedChange.windowLeft, change.WindowLeft())
			})
		}
	})
}

func TestMouse_Scroll(t *testing.T) {
	t.Run("before Update was called, Scroll returns 0, 0", func(t *testing.T) {
		source := newFakeEventSource()
		mouseState := mouse.New(source)
		// when
		scroll := mouseState.Scroll()
		// then
		assert.Equal(t, 0.0, scroll.X())
		assert.Equal(t, 0.0, scroll.Y())
	})
	t.Run("after Update was called", func(t *testing.T) {
		tests := map[string]struct {
			source               mouse.EventSource
			expectedX, expectedY float64
		}{
			"no events": {
				source: newFakeEventSource(),
			},
			"one event": {
				source:    newFakeEventSource(mouse.NewScrolledEvent(1.1, 2.2)),
				expectedX: 1.1,
				expectedY: 2.2,
			},
			"two events": {
				source:    newFakeEventSource(mouse.NewScrolledEvent(1.1, 2.2), mouse.NewScrolledEvent(0.3, 0.5)),
				expectedX: 1.4,
				expectedY: 2.7,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				mouseState := mouse.New(test.source)
				mouseState.Update()
				// when
				scroll := mouseState.Scroll()
				// then
				assert.InDelta(t, test.expectedX, scroll.X(), 0.000000001)
				assert.InDelta(t, test.expectedY, scroll.Y(), 0.000000001)
			})
		}
	})
	t.Run("after second Update was called", func(t *testing.T) {
		tests := map[string]struct {
			newEvents            []mouse.Event
			expectedX, expectedY float64
		}{
			"no new events": {
				expectedX: 0.0,
				expectedY: 0.0,
			},
			"one event": {
				newEvents: []mouse.Event{mouse.NewScrolledEvent(0.1, 0.2)},
				expectedX: 0.1,
				expectedY: 0.2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				source := newFakeEventSource(mouse.NewScrolledEvent(1.1, 2.2))
				mouseState := mouse.New(source)
				mouseState.Update()
				source.events = test.newEvents
				mouseState.Update()
				// when
				scroll := mouseState.Scroll()
				// then
				assert.InDelta(t, test.expectedX, scroll.X(), 0.000000001)
				assert.InDelta(t, test.expectedY, scroll.Y(), 0.000000001)
			})
		}
	})
}

func newFakeEventSource(events ...mouse.Event) *fakeEventSource {
	source := &fakeEventSource{}
	source.events = []mouse.Event{}
	source.events = append(source.events, events...)
	return source
}

type fakeEventSource struct {
	events []mouse.Event
}

func (f *fakeEventSource) PollMouseEvent() (mouse.Event, bool) {
	if len(f.events) > 0 {
		event := f.events[0]
		f.events = f.events[1:]
		return event, true
	}
	return mouse.EmptyEvent, false
}

func newFakeEventSourceWith(numberOfEvents int) *fakeEventSource {
	source := newFakeEventSource()
	for i := 0; i < numberOfEvents; i++ {
		source.events = append(source.events, mouse.NewPressedEvent(mouse.Left))
	}
	return source
}
