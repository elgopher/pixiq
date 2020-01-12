package keyboard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/keyboard"
)

func TestNewEventQueue(t *testing.T) {
	t.Run("should create EventQueue", func(t *testing.T) {
		sizes := []int{-1, 1, 1, 16}
		for _, size := range sizes {
			buffer := keyboard.NewEventQueue(size)
			assert.NotNil(t, buffer)
		}
	})
}

func TestEventQueue_Poll(t *testing.T) {
	t.Run("should return EmptyEvent and false for empty EventQueue", func(t *testing.T) {
		queue := keyboard.NewEventQueue(1)
		// when
		event, ok := queue.Poll()
		// then
		assert.False(t, ok)
		assert.Equal(t, keyboard.EmptyEvent, event)
	})
}

func TestEventQueue_Append(t *testing.T) {
	event1 := keyboard.NewPressedEvent(keyboard.A)
	event2 := keyboard.NewPressedEvent(keyboard.B)
	event3 := keyboard.NewPressedEvent(keyboard.C)

	t.Run("should append events to EventQueue with enough space", func(t *testing.T) {
		tests := map[string][]keyboard.Event{
			"one event":    {event1},
			"two events":   {event1, event2},
			"three events": {event1, event2, event3},
		}
		for name, events := range tests {
			t.Run(name, func(t *testing.T) {
				queue := keyboard.NewEventQueue(3)
				// when
				for _, event := range events {
					queue.Append(event)
				}
				// then
				for _, event := range events {
					actualEvent, found := queue.Poll()
					assert.True(t, found)
					assert.Equal(t, event, actualEvent)
				}
				// and
				actualEvent, found := queue.Poll()
				assert.False(t, found)
				assert.Equal(t, keyboard.EmptyEvent, actualEvent)
			})
		}
	})
	t.Run("should override old events when EventQueue has not enough space", func(t *testing.T) {
		tests := map[string]struct {
			size           int
			events         []keyboard.Event
			expectedEvents []keyboard.Event
		}{
			"size 1": {
				size:           1,
				events:         []keyboard.Event{event1, event2},
				expectedEvents: []keyboard.Event{event2},
			},
			"size 2": {
				size:           2,
				events:         []keyboard.Event{event1, event2, event3},
				expectedEvents: []keyboard.Event{event2, event3},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				events := test.events
				queue := keyboard.NewEventQueue(test.size)
				// when
				for _, event := range events {
					queue.Append(event)
				}
				// then
				for _, event := range test.expectedEvents {
					actualEvent, found := queue.Poll()
					assert.True(t, found)
					assert.Equal(t, event, actualEvent)
				}
				// and
				actualEvent, found := queue.Poll()
				assert.False(t, found)
				assert.Equal(t, keyboard.EmptyEvent, actualEvent)
			})
		}
	})
}
