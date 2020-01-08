package internal_test

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl/internal"
)

func TestKeyboardEvents_Poll(t *testing.T) {
	t.Run("should return EmptyEvent when there are no events", func(t *testing.T) {
		events := internal.KeyboardEvents{}
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
				events := internal.KeyboardEvents{}
				events.OnKeyCallback(nil, test.glfwKey, test.scanCode, test.action, 0)
				// when
				event, ok := events.Poll()
				// then
				require.True(t, ok)
				assert.Equal(t, test.expectedEvent, event)
				// and
				event, ok = events.Poll()
				require.False(t, ok)
				assert.Equal(t, keyboard.EmptyEvent, event)
			})
		}

	})
	t.Run("should return two mapped events", func(t *testing.T) {
		events := internal.KeyboardEvents{}
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
		event, ok = events.Poll()
		require.False(t, ok)
		assert.Equal(t, keyboard.EmptyEvent, event)
	})
}
