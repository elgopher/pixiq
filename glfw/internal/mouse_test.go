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
		assert.NotNil(t, internal.NewMouseEvents(buffer, &fakeWindow{}))
	})
	t.Run("should panic for nil buffer", func(t *testing.T) {
		assert.Panics(t, func() {
			assert.NotNil(t, internal.NewMouseEvents(nil, &fakeWindow{}))
		})
	})
	t.Run("should panic for nil window", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(1)
		assert.Panics(t, func() {
			assert.NotNil(t, internal.NewMouseEvents(buffer, nil))
		})
	})
}

func TestMouseEvents_Poll(t *testing.T) {
	t.Run("should return EmptyEvent when there are no events", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(1)
		events := internal.NewMouseEvents(buffer, &fakeWindow{})
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
				events := internal.NewMouseEvents(buffer, &fakeWindow{})
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
				events := internal.NewMouseEvents(buffer, &fakeWindow{})
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
		events := internal.NewMouseEvents(buffer, &fakeWindow{})
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

	t.Run("should generate MoveEvent", func(t *testing.T) {
		tests := map[string]struct {
			window        internal.Window
			expectedEvent mouse.Event
		}{
			"1,2": {
				window: &fakeWindow{
					posX:   1,
					posY:   2,
					width:  2,
					height: 3,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(1, 2, 1, 2, true),
			},
			"zoom 2": {
				window: &fakeWindow{
					posX:   2.0,
					posY:   4.0,
					width:  3,
					height: 5,
					zoom:   2,
				},
				expectedEvent: mouse.NewMovedEvent(1, 2, 2.0, 4.0, true),
			},
			"outside window, x == width": {
				window: &fakeWindow{
					posX:   1,
					posY:   0,
					width:  1,
					height: 1,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(1, 0, 1, 0, false),
			},
			"outside window, x >= width": {
				window: &fakeWindow{
					posX:   2,
					posY:   0,
					width:  1,
					height: 1,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(2, 0, 2, 0, false),
			},
			"outside window, x < width": {
				window: &fakeWindow{
					posX:   -1,
					posY:   0,
					width:  1,
					height: 1,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(-1, 0, -1, 0, false),
			},
			"outside window, y == height": {
				window: &fakeWindow{
					posX:   0,
					posY:   1,
					width:  1,
					height: 1,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(0, 1, 0, 1, false),
			},
			"outside window, y >= height": {
				window: &fakeWindow{
					posX:   0,
					posY:   2,
					width:  1,
					height: 1,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(0, 2, 0, 2, false),
			},
			"outside window, y < 0": {
				window: &fakeWindow{
					posX:   0,
					posY:   -1,
					width:  1,
					height: 1,
					zoom:   1,
				},
				expectedEvent: mouse.NewMovedEvent(0, -1, 0, -1, false),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := mouse.NewEventBuffer(1)
				events := internal.NewMouseEvents(buffer, test.window)
				// when
				event, ok := events.Poll()
				// then
				assert.True(t, ok)
				assert.Equal(t, test.expectedEvent, event)
			})
		}

	})

	t.Run("same position should not generate MoveEvent", func(t *testing.T) {
		buffer := mouse.NewEventBuffer(1)
		window := &fakeWindow{
			posX:   1,
			posY:   1,
			width:  1,
			height: 1,
			zoom:   1,
		}
		events := internal.NewMouseEvents(buffer, window)
		_, _ = events.Poll()
		event, ok := events.Poll()
		// then
		assert.False(t, ok)
		assert.Equal(t, mouse.EmptyEvent, event)
	})

	t.Run("when position changes should generate MoveEvent", func(t *testing.T) {
		tests := map[string]struct {
			newPosX, newPosY float64
			expectedEvent    mouse.Event
		}{
			"inside window": {
				newPosX:       2,
				newPosY:       3,
				expectedEvent: mouse.NewMovedEvent(2, 3, 2, 3, true),
			},
			"outside window": {
				newPosX:       10,
				newPosY:       20,
				expectedEvent: mouse.NewMovedEvent(10, 20, 10, 20, false),
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := mouse.NewEventBuffer(1)
				window := &fakeWindow{
					posX:   1,
					posY:   1,
					width:  6,
					height: 6,
					zoom:   1,
				}
				events := internal.NewMouseEvents(buffer, window)
				_, _ = events.Poll()
				window.posX = test.newPosX
				window.posY = test.newPosY
				event, ok := events.Poll()
				// then
				assert.True(t, ok)
				assert.Equal(t, test.expectedEvent, event)
			})
		}
	})
}

func assertNoMoreMouseEvents(t *testing.T, events *internal.MouseEvents) {
	event, ok := events.Poll()
	require.False(t, ok)
	assert.Equal(t, mouse.EmptyEvent, event)
}

type fakeWindow struct {
	posX, posY    float64
	width, height int
	zoom          int
}

func (f *fakeWindow) CursorPosition() (float64, float64) {
	return f.posX, f.posY
}

func (f *fakeWindow) Size() (int, int) {
	return f.width, f.height
}

func (f *fakeWindow) Zoom() int {
	return f.zoom
}
