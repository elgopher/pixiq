package internal_test

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glfw/internal"
	"github.com/jacekolszak/pixiq/mouse"
)

func TestNewMouseEvents(t *testing.T) {
	t.Run("should create MouseEvents when buffer is given", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(1)
		// expect
		assert.NotNil(t, internal.NewMouseEvents(buffer))
	})
	t.Run("should panic for nil buffer", func(t *testing.T) {
		assert.Panics(t, func() {
			assert.NotNil(t, internal.NewMouseEvents(nil))
		})
	})
}

func TestMouseEvents_Poll(t *testing.T) {
	t.Run("should return EmptyEvent when there are no events", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(1)
		events := internal.NewMouseEvents(buffer)
		// when
		event, ok := events.Poll()
		// then
		require.False(t, ok)
		assert.Equal(t, mouse.EmptyEvent, event)
	})

	t.Run("should map mouse button event", func(t *testing.T) {
		tests := map[string]struct {
			button        glfw.MouseButton
			action        glfw.Action
			expectedEvent mouse.Event
		}{
			"press left": {
				button:        glfw.MouseButtonLeft,
				action:        glfw.Press,
				expectedEvent: mouse.NewPressedEvent(mouse.Left),
			},
			"press right": {
				button:        glfw.MouseButtonRight,
				action:        glfw.Press,
				expectedEvent: mouse.NewPressedEvent(mouse.Right),
			},
			"press 1": {
				button:        glfw.MouseButton1,
				action:        glfw.Press,
				expectedEvent: mouse.NewPressedEvent(mouse.Left),
			},
			"release left": {
				button:        glfw.MouseButtonLeft,
				action:        glfw.Release,
				expectedEvent: mouse.NewReleasedEvent(mouse.Left),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := mouse.NewEventBuffer(1)
				events := internal.NewMouseEvents(buffer)
				// when
				events.OnMouseButtonCallback(nil, test.button, test.action, 0)
				event, ok := events.Poll()
				// then
				require.True(t, ok)
				assert.Equal(t, test.expectedEvent, event)
			})
		}
	})
	t.Run("should map scroll event", func(t *testing.T) {
		tests := map[string]struct {
			x, y          float64
			expectedEvent mouse.Event
		}{
			"0,0": {
				expectedEvent: mouse.NewScrolledEvent(0, 0),
			},
			"1,2": {
				x:             1,
				y:             2,
				expectedEvent: mouse.NewScrolledEvent(1, 2),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := mouse.NewEventBuffer(1)
				events := internal.NewMouseEvents(buffer)
				// when
				events.OnScrollCallback(nil, test.x, test.y)
				event, ok := events.Poll()
				// then
				require.True(t, ok)
				assert.Equal(t, test.expectedEvent, event)
			})
		}
	})

	t.Run("should return two button events", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(2)
		events := internal.NewMouseEvents(buffer)
		events.OnMouseButtonCallback(nil, glfw.MouseButtonLeft, glfw.Press, 0)
		events.OnMouseButtonCallback(nil, glfw.MouseButtonRight, glfw.Release, 0)
		// when
		event, ok := events.Poll()
		// then
		require.True(t, ok)
		assert.Equal(t, mouse.NewPressedEvent(mouse.Left), event)
		// and
		event, ok = events.Poll()
		require.True(t, ok)
		assert.Equal(t, mouse.NewReleasedEvent(mouse.Right), event)
		// and
		assertNoMoreMouseEvents(t, events)
	})

}

func assertNoMoreMouseEvents(t *testing.T, events *internal.MouseEvents) {
	event, ok := events.Poll()
	require.False(t, ok)
	assert.Equal(t, mouse.EmptyEvent, event)
}
