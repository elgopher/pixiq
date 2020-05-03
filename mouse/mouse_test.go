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
		keys := mouse.New(source)
		// then
		assert.NotNil(t, keys)
	})
}

func TestMouse_Pressed(t *testing.T) {
	t.Run("before Update was called, Pressed returns false for all buttons", func(t *testing.T) {
		tests := []mouse.Button{mouse.Left, mouse.Right, mouse.Middle, mouse.Button4, mouse.Button5, mouse.Button6, mouse.Button7, mouse.Button8}
		for _, key := range tests {
			testName := fmt.Sprintf("for key: %v", key)
			t.Run(testName, func(t *testing.T) {
				var (
					event  = mouse.NewPressedEvent(mouse.Left)
					source = newFakeEventSource(event)
					keys   = mouse.New(source)
				)
				// when
				pressed := keys.Pressed(key)
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
			"one PressedEvent for A; one ReleasedEvent for B": {
				source:             newFakeEventSource(leftPressed, rightReleased),
				expectedPressed:    []mouse.Button{mouse.Left},
				expectedNotPressed: []mouse.Button{mouse.Right},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				keys := mouse.New(test.source)
				// when
				keys.Update()
				// then
				for _, expectedPressedKey := range test.expectedPressed {
					assert.True(t, keys.Pressed(expectedPressedKey))
				}
				for _, expectedNotPressedKey := range test.expectedNotPressed {
					assert.False(t, keys.Pressed(expectedNotPressedKey))
				}
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
