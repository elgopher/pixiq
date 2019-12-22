package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

func TestNewWindows(t *testing.T) {
	t.Run("should return Windows object for creating windows", func(t *testing.T) {
		windows := pixiq.NewWindows(&fakeAcceleratedImages{}, openWindowMock)
		assert.NotNil(t, windows)
	})
}

func TestWindow_New(t *testing.T) {
	windows := pixiq.NewWindows(&fakeAcceleratedImages{}, openWindowMock)
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
		windows := pixiq.NewWindows(&fakeAcceleratedImages{}, openWindowMock)
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
		windows := pixiq.NewWindows(&fakeAcceleratedImages{}, openWindowMock)
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

	t.Run("should draw accelerated image for each frame", func(t *testing.T) {
		windowMock := &systemWindowMock{}
		openWindowMock := func(width, height int) pixiq.SystemWindow {
			return windowMock
		}
		images := &fakeAcceleratedImages{}
		windows := pixiq.NewWindows(images, openWindowMock)
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
		assert.Equal(t, []pixiq.AcceleratedImage{images.images[0], images.images[0]}, windowMock.imagesDrawn)
	})

	t.Run("should upload initial, transparent image", func(t *testing.T) {
		t.Run("0x0", func(t *testing.T) {
			images := &fakeAcceleratedImages{}
			windows := pixiq.NewWindows(images, openWindowMock)
			window := windows.New(0, 0)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				frame.CloseWindowEventually()
			})
			// then
			images.assertOneImageWithPixels(t, []pixiq.Color{})
		})
		t.Run("1x1", func(t *testing.T) {
			images := &fakeAcceleratedImages{}
			windows := pixiq.NewWindows(images, openWindowMock)
			window := windows.New(1, 1)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				frame.CloseWindowEventually()
			})
			// then
			images.assertOneImageWithPixels(t, []pixiq.Color{transparent})
		})
	})

	t.Run("should upload modified window image", func(t *testing.T) {
		t.Run("1x1", func(t *testing.T) {
			images := &fakeAcceleratedImages{}
			windows := pixiq.NewWindows(images, openWindowMock)
			window := windows.New(1, 1)
			color := pixiq.RGBA(10, 20, 30, 40)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				selection := frame.Image().Selection(0, 0)
				selection.SetColor(0, 0, color)
				frame.CloseWindowEventually()
			})
			// then
			images.assertOneImageWithPixels(t, []pixiq.Color{color})
		})
		t.Run("1x2", func(t *testing.T) {
			images := &fakeAcceleratedImages{}
			windows := pixiq.NewWindows(images, openWindowMock)
			window := windows.New(1, 2)
			color0 := pixiq.RGBA(10, 20, 30, 40)
			color1 := pixiq.RGBA(50, 60, 70, 80)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				selection := frame.Image().Selection(0, 0)
				selection.SetColor(0, 0, color0)
				selection.SetColor(0, 1, color1)
				frame.CloseWindowEventually()
			})
			// then
			images.assertOneImageWithPixels(t, []pixiq.Color{color0, color1})
		})
		t.Run("2x1", func(t *testing.T) {
			images := &fakeAcceleratedImages{}
			windows := pixiq.NewWindows(images, openWindowMock)
			window := windows.New(2, 1)
			color0 := pixiq.RGBA(10, 20, 30, 40)
			color1 := pixiq.RGBA(50, 60, 70, 80)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				selection := frame.Image().Selection(0, 0)
				selection.SetColor(0, 0, color0)
				selection.SetColor(1, 0, color1)
				frame.CloseWindowEventually()
			})
			// then
			images.assertOneImageWithPixels(t, []pixiq.Color{color0, color1})
		})
		t.Run("2x2", func(t *testing.T) {
			images := &fakeAcceleratedImages{}
			windows := pixiq.NewWindows(images, openWindowMock)
			window := windows.New(2, 2)
			color0 := pixiq.RGBA(10, 20, 30, 40)
			color1 := pixiq.RGBA(50, 60, 70, 80)
			color2 := pixiq.RGBA(90, 100, 110, 120)
			color3 := pixiq.RGBA(130, 140, 150, 160)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				selection := frame.Image().Selection(0, 0)
				selection.SetColor(0, 0, color0)
				selection.SetColor(1, 0, color1)
				selection.SetColor(0, 1, color2)
				selection.SetColor(1, 1, color3)
				frame.CloseWindowEventually()
			})
			// then
			images.assertOneImageWithPixels(t, []pixiq.Color{color0, color1, color2, color3})
		})
	})

}
