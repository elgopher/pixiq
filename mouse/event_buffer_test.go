package mouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/mouse"
)

func TestNewEventBuffer(t *testing.T) {
	t.Run("should create EventBuffer", func(t *testing.T) {
		sizes := []int{-1, 1, 1, 16}
		for _, size := range sizes {
			buffer := mouse.NewEventBuffer(size)
			assert.NotNil(t, buffer)
		}
	})
}

func TestEventBuffer_Poll(t *testing.T) {
	t.Run("should return EmptyEvent and false for empty EventBuffer", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(1)
		// when
		event, ok := buffer.Poll()
		// then
		assert.False(t, ok)
		assert.Equal(t, mouse.EmptyEvent, event)
	})
}

func TestEventBuffer_Add(t *testing.T) {
	event1 := mouse.NewPressedEvent(mouse.Left)
	event2 := mouse.NewScrolledEvent(1, 2)
	event3 := mouse.NewMovedEvent(1, 2, 1, 2, true)

	t.Run("should add events to EventBuffer with enough space", func(t *testing.T) {
		tests := map[string][]mouse.Event{
			"one event":    {event1},
			"two events":   {event1, event2},
			"three events": {event1, event2, event3},
		}
		for name, events := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := mouse.NewEventBuffer(3)
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
				assert.Equal(t, mouse.EmptyEvent, actualEvent)
			})
		}
	})
	t.Run("should override old events when EventBuffer has not enough space", func(t *testing.T) {
		tests := map[string]struct {
			buffer         *mouse.EventBuffer
			events         []mouse.Event
			expectedEvents []mouse.Event
		}{
			"size 1": {
				buffer:         mouse.NewEventBuffer(1),
				events:         []mouse.Event{event1, event2},
				expectedEvents: []mouse.Event{event2},
			},
			"size 2": {
				buffer:         mouse.NewEventBuffer(2),
				events:         []mouse.Event{event1, event2, event3},
				expectedEvents: []mouse.Event{event2, event3},
			},
			"already added and polled": {
				buffer: prepare(mouse.NewEventBuffer(2), func(q *mouse.EventBuffer) {
					q.Add(event1)
					q.Poll()
				}),
				events:         []mouse.Event{event2, event3},
				expectedEvents: []mouse.Event{event2, event3},
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
				assert.Equal(t, mouse.EmptyEvent, actualEvent)
			})
		}
	})
}

func prepare(b *mouse.EventBuffer, f func(q *mouse.EventBuffer)) *mouse.EventBuffer {
	f(b)
	return b
}
