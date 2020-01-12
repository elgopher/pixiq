package keyboard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/keyboard"
)

func TestNewEventBuffer(t *testing.T) {
	t.Run("should create EventBuffer", func(t *testing.T) {
		sizes := []int{-1, 1, 1, 16}
		for _, size := range sizes {
			buffer := keyboard.NewEventBuffer(size)
			assert.NotNil(t, buffer)
		}
	})
}

func TestEventBuffer_Poll(t *testing.T) {
	t.Run("should return EmptyEvent and false for empty EventBuffer", func(t *testing.T) {
		buffer := keyboard.NewEventBuffer(1)
		// when
		event, ok := buffer.Poll()
		// then
		assert.False(t, ok)
		assert.Equal(t, keyboard.EmptyEvent, event)
	})
}

func TestEventBuffer_Add(t *testing.T) {
	event1 := keyboard.NewPressedEvent(keyboard.One)
	event2 := keyboard.NewPressedEvent(keyboard.Two)
	event3 := keyboard.NewPressedEvent(keyboard.Three)

	t.Run("should add events to EventBuffer with enough space", func(t *testing.T) {
		tests := map[string][]keyboard.Event{
			"one event":    {event1},
			"two events":   {event1, event2},
			"three events": {event1, event2, event3},
		}
		for name, events := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := keyboard.NewEventBuffer(3)
				// when
				for _, event := range events {
					buffer.Add(event)
				}
				// then
				for _, event := range events {
					actualEvent, found := buffer.Poll()
					assert.True(t, found)
					assert.Equal(t, event, actualEvent)
				}
				// and
				actualEvent, found := buffer.Poll()
				assert.False(t, found)
				assert.Equal(t, keyboard.EmptyEvent, actualEvent)
			})
		}
	})
	t.Run("should override old events when EventBuffer has not enough space", func(t *testing.T) {
		tests := map[string]struct {
			buffer         *keyboard.EventBuffer
			events         []keyboard.Event
			expectedEvents []keyboard.Event
		}{
			"size 1": {
				buffer:         keyboard.NewEventBuffer(1),
				events:         []keyboard.Event{event1, event2},
				expectedEvents: []keyboard.Event{event2},
			},
			"size 2": {
				buffer:         keyboard.NewEventBuffer(2),
				events:         []keyboard.Event{event1, event2, event3},
				expectedEvents: []keyboard.Event{event2, event3},
			},
			"already added and polled": {
				buffer: prepare(keyboard.NewEventBuffer(2), func(q *keyboard.EventBuffer) {
					q.Add(event1)
					q.Poll()
				}),
				events:         []keyboard.Event{event2, event3},
				expectedEvents: []keyboard.Event{event2, event3},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				events := test.events
				// when
				for _, event := range events {
					test.buffer.Add(event)
				}
				// then
				for _, event := range test.expectedEvents {
					actualEvent, found := test.buffer.Poll()
					assert.True(t, found)
					assert.Equal(t, event, actualEvent)
				}
				// and
				actualEvent, found := test.buffer.Poll()
				assert.False(t, found)
				assert.Equal(t, keyboard.EmptyEvent, actualEvent)
			})
		}
	})
}

func prepare(q *keyboard.EventBuffer, f func(q *keyboard.EventBuffer)) *keyboard.EventBuffer {
	f(q)
	return q
}
