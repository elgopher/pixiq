package internal_test

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl/internal"
)

func TestNewKeyboardEvents(t *testing.T) {
	t.Run("should panic for initial size -1", func(t *testing.T) {
		assert.Panics(t, func() {
			internal.NewKeyboardEvents(-1)
		})
	})
	t.Run("should panic for initial size 0", func(t *testing.T) {
		assert.Panics(t, func() {
			internal.NewKeyboardEvents(0)
		})
	})
}

func TestKeyboardEvents_Poll(t *testing.T) {
	t.Run("should return EmptyEvent when there are no events", func(t *testing.T) {
		events := internal.NewKeyboardEvents(16)
		// when
		event, ok := events.Poll()
		// then
		require.False(t, ok)
		assert.Equal(t, keyboard.EmptyEvent, event)
	})
	t.Run("should return EmptyEvent for Repeat action", func(t *testing.T) {
		events := internal.NewKeyboardEvents(16)
		events.OnKeyCallback(nil, glfw.KeyA, 0, glfw.Repeat, 0)
		// when
		event, ok := events.Poll()
		// then
		require.False(t, ok)
		assert.Equal(t, keyboard.EmptyEvent, event)
	})
	t.Run("should return mapped event", func(t *testing.T) {
		tests := map[string]struct {
			glfwKey       glfw.Key
			scanCode      int
			action        glfw.Action
			expectedEvent keyboard.Event
		}{
			"pressed A": {
				glfwKey:       glfw.KeyA,
				action:        glfw.Press,
				expectedEvent: keyboard.NewPressedEvent(keyboard.A),
			},
			"pressed B": {
				glfwKey:       glfw.KeyB,
				action:        glfw.Press,
				expectedEvent: keyboard.NewPressedEvent(keyboard.B),
			},
			"pressed Unknown 0": {
				glfwKey:       glfw.KeyUnknown,
				scanCode:      0,
				action:        glfw.Press,
				expectedEvent: keyboard.NewPressedEvent(keyboard.NewUnknownKey(0)),
			},
			"pressed Unknown 1": {
				glfwKey:       glfw.KeyUnknown,
				scanCode:      1,
				action:        glfw.Press,
				expectedEvent: keyboard.NewPressedEvent(keyboard.NewUnknownKey(1)),
			},
			"released A": {
				glfwKey:       glfw.KeyA,
				action:        glfw.Release,
				expectedEvent: keyboard.NewReleasedEvent(keyboard.A),
			},
			"released Unknown 0": {
				glfwKey:       glfw.KeyUnknown,
				action:        glfw.Release,
				expectedEvent: keyboard.NewReleasedEvent(keyboard.NewUnknownKey(0)),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				events := internal.NewKeyboardEvents(16)
				events.OnKeyCallback(nil, test.glfwKey, test.scanCode, test.action, 0)
				// when
				event, ok := events.Poll()
				// then
				require.True(t, ok)
				assert.Equal(t, test.expectedEvent, event)
				// and
				assertNoMoreEvents(t, events)
			})
		}
	})

	t.Run("should return two mapped events", func(t *testing.T) {
		events := internal.NewKeyboardEvents(16)
		events.OnKeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
		events.OnKeyCallback(nil, glfw.KeyB, 0, glfw.Release, 0)
		// when
		event, ok := events.Poll()
		// then
		require.True(t, ok)
		assert.Equal(t, keyboard.NewPressedEvent(keyboard.A), event)
		// and
		event, ok = events.Poll()
		require.True(t, ok)
		assert.Equal(t, keyboard.NewReleasedEvent(keyboard.B), event)
		// and
		assertNoMoreEvents(t, events)
	})

}

func TestKeyboardEvents_OnKeyCallback(t *testing.T) {
	t.Run("should expand buffer when too many events", func(t *testing.T) {
		events := internal.NewKeyboardEvents(1)
		// when
		events.OnKeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
		events.OnKeyCallback(nil, glfw.KeyB, 0, glfw.Press, 0)
		// then
		event, ok := events.Poll()
		require.True(t, ok)
		assert.Equal(t, keyboard.NewPressedEvent(keyboard.A), event)
		// and
		event, ok = events.Poll()
		require.True(t, ok)
		assert.Equal(t, keyboard.NewPressedEvent(keyboard.B), event)
		// and
		assertNoMoreEvents(t, events)
	})
}

func TestKeyboardEvents_Clear(t *testing.T) {
	t.Run("should clear events", func(t *testing.T) {
		tests := map[string]int{
			"no events":  0,
			"one event":  1,
			"two events": 2,
		}
		for name, numberOfEents := range tests {
			t.Run(name, func(t *testing.T) {
				events := internal.NewKeyboardEvents(1)
				for i := 0; i < numberOfEents; i++ {
					events.OnKeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
				}
				// when
				events.Clear()
				// then
				assertNoMoreEvents(t, events)
			})
		}
	})
}

func assertNoMoreEvents(t *testing.T, events *internal.KeyboardEvents) {
	event, ok := events.Poll()
	require.False(t, ok)
	assert.Equal(t, keyboard.EmptyEvent, event)
}
