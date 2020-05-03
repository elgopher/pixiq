package mouse_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/mouse"
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
