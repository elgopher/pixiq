package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

func TestNewWindows(t *testing.T) {
	t.Run("should return Windows object for creating windows", func(t *testing.T) {
		windows := pixiq.NewWindows()
		assert.NotNil(t, windows)
	})
}

func TestWindow_New(t *testing.T) {
	windows := pixiq.NewWindows()
	t.Run("should clamp width to 0 if negative", func(t *testing.T) {
		win := windows.New(-1, 0)
		require.NotNil(t, win)
		assert.Equal(t, 0, win.Width())
	})
	t.Run("should clamp height to 0 if negative", func(t *testing.T) {
		win := windows.New(0, -1)
		require.NotNil(t, win)
		assert.Equal(t, 0, win.Height())
	})
	t.Run("should create window", func(t *testing.T) {
		win := windows.New(1, 2)
		require.NotNil(t, win)
		assert.Equal(t, 1, win.Width())
		assert.Equal(t, 2, win.Height())
	})
}

func TestWindow_Loop(t *testing.T) {

	t.Run("should run callback function until window is closed", func(t *testing.T) {
		windows := pixiq.NewWindows()
		window := windows.New(0, 0)
		executionCount := 0
		// when
		window.Loop(func(frame *pixiq.Frame) {
			executionCount += 1
			if executionCount == 2 {
				frame.CloseWindowEventually()
			}
		})
		// then
		assert.Equal(t, 2, executionCount)
	})

	t.Run("frame should provide Image for the whole window", func(t *testing.T) {
		windows := pixiq.NewWindows()
		tests := map[string]struct {
			width, height int
		}{
			"0x0": {
				width:  0,
				height: 0,
			},
			"1x2": {
				width:  1,
				height: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				window := windows.New(test.width, test.height)
				var image *pixiq.Image
				// when
				window.Loop(func(frame *pixiq.Frame) {
					image = frame.Image()
					frame.CloseWindowEventually()
				})
				// then
				assert.Equal(t, test.width, image.Width())
				assert.Equal(t, test.height, image.Height())
			})

		}
	})

}
