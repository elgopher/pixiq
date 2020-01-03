package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

func TestNewWindows(t *testing.T) {
	t.Run("should return Windows object for creating windows", func(t *testing.T) {
		windows := pixiq.NewWindows(pixiq.NewImages(&fakeAcceleratedImages{}), &systemWindowsMock{})
		assert.NotNil(t, windows)
	})
}

func TestWindow_New(t *testing.T) {
	windows := pixiq.NewWindows(pixiq.NewImages(&fakeAcceleratedImages{}), &systemWindowsMock{})
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
	t.Run("New should not open the system window", func(t *testing.T) {
		systemWindows := &systemWindowsMock{}
		windows := pixiq.NewWindows(pixiq.NewImages(&fakeAcceleratedImages{}), systemWindows)
		// when
		windows.New(0, 0)
		// then
		assert.Empty(t, systemWindows.openWindows)
	})
	t.Run("should pass hints to system window", func(t *testing.T) {
		t.Run("one hint", func(t *testing.T) {
			systemWindows := &systemWindowsMock{}
			windows := pixiq.NewWindows(pixiq.NewImages(&fakeAcceleratedImages{}), systemWindows)
			hint := &fakeHint{}
			// when
			windows.New(0, 0, hint).Loop(func(frame *pixiq.Frame) {
				frame.CloseWindowEventually()
			})
			// then
			require.Len(t, systemWindows.openWindows[0].hints, 1)
			assert.Same(t, hint, systemWindows.openWindows[0].hints[0].(*fakeHint))
		})
		t.Run("two hints", func(t *testing.T) {
			systemWindows := &systemWindowsMock{}
			windows := pixiq.NewWindows(pixiq.NewImages(&fakeAcceleratedImages{}), systemWindows)
			hint1 := &fakeHint{}
			hint2 := &fakeHint{}
			// when
			windows.New(0, 0, hint1, hint2).Loop(func(frame *pixiq.Frame) {
				frame.CloseWindowEventually()
			})
			// then
			window := systemWindows.openWindows[0]
			require.Len(t, window.hints, 2)
			assert.Same(t, hint1, window.hints[0].(*fakeHint))
			assert.Same(t, hint2, window.hints[1].(*fakeHint))
		})
	})
}

type fakeHint struct {
}

func TestWindow_Loop(t *testing.T) {

	t.Run("should run callback function until window is closed", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		windows := pixiq.NewWindows(images, &systemWindowsMock{})
		window := windows.New(0, 0)
		frameNumber := 0
		// when
		window.Loop(func(frame *pixiq.Frame) {
			if frameNumber == 2 {
				frame.CloseWindowEventually()
			} else {
				frameNumber += 1
			}
		})
		// then
		assert.Equal(t, 2, frameNumber)
	})

	t.Run("frame should provide Image for the whole window", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		windows := pixiq.NewWindows(images, &systemWindowsMock{})
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

	t.Run("should open system window", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		systemWindows := &systemWindowsMock{}
		windows := pixiq.NewWindows(images, systemWindows)
		window := windows.New(1, 2)
		// when
		window.Loop(func(frame *pixiq.Frame) {
			frame.CloseWindowEventually()
		})
		// then
		require.Len(t, systemWindows.openWindows, 1)
		win := systemWindows.openWindows[0]
		assert.Equal(t, 1, win.width)
		assert.Equal(t, 2, win.height)
	})

	t.Run("should draw image for each frame", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		systemWindows := &systemWindowsMock{}
		windows := pixiq.NewWindows(images, systemWindows)
		window := windows.New(0, 0)
		firstFrame := true
		var recordedImages []*pixiq.Image
		// when
		window.Loop(func(frame *pixiq.Frame) {
			if !firstFrame {
				frame.CloseWindowEventually()
			}
			firstFrame = false
			recordedImages = append(recordedImages, frame.Image())
		})
		// then
		require.Len(t, systemWindows.openWindows, 1)
		win := systemWindows.openWindows[0]
		assert.Equal(t, recordedImages, win.imagesDrawn)
	})

	t.Run("initial window image is transparent", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		windows := pixiq.NewWindows(images, &systemWindowsMock{})
		window := windows.New(1, 1)
		var image *pixiq.Image
		// when
		window.Loop(func(frame *pixiq.Frame) {
			image = frame.Image()
			frame.CloseWindowEventually()
		})
		// then
		assert.Equal(t, transparent, image.WholeImageSelection().Color(0, 0))
	})

	t.Run("should draw modified window image", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			color := pixiq.RGBA(10, 20, 30, 40)
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			systemWindows := &systemWindowsMock{}
			windows := pixiq.NewWindows(images, systemWindows)
			window := windows.New(1, 1)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				selection := frame.Image().Selection(0, 0)
				selection.SetColor(0, 0, color)
				frame.CloseWindowEventually()
			})
			// then
			require.Len(t, systemWindows.openWindows, 1)
			win := systemWindows.openWindows[0]
			require.Len(t, win.imagesDrawn, 1)
			assert.Equal(t, color, win.imagesDrawn[0].WholeImageSelection().Color(0, 0))
		})
		t.Run("after second frame", func(t *testing.T) {
			color1 := pixiq.RGBA(10, 20, 30, 40)
			color2 := pixiq.RGBA(10, 20, 30, 40)
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			systemWindows := &systemWindowsMock{}
			windows := pixiq.NewWindows(images, systemWindows)
			window := windows.New(1, 1)
			firstFrame := true
			// when
			window.Loop(func(frame *pixiq.Frame) {
				selection := frame.Image().WholeImageSelection()
				if firstFrame {
					selection.SetColor(0, 0, color1)
					firstFrame = false
				} else {
					selection.SetColor(0, 0, color2)
					frame.CloseWindowEventually()
				}
			})
			// then
			require.Len(t, systemWindows.openWindows, 1)
			win := systemWindows.openWindows[0]
			require.Len(t, win.imagesDrawn, 2)
			assert.Equal(t, color1, win.imagesDrawn[0].WholeImageSelection().Color(0, 0))
			assert.Equal(t, color2, win.imagesDrawn[1].WholeImageSelection().Color(0, 0))
		})
	})

	t.Run("should open window with given dimensions", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		systemWindows := &systemWindowsMock{}
		windows := pixiq.NewWindows(images, systemWindows)
		window := windows.New(1, 2)
		// when
		window.Loop(func(frame *pixiq.Frame) {
			frame.CloseWindowEventually()
		})
		// then
		assert.Equal(t, 1, systemWindows.openWindows[0].width)
		assert.Equal(t, 2, systemWindows.openWindows[0].height)
	})

	t.Run("should swap images", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			systemWindows := &systemWindowsMock{}
			windows := pixiq.NewWindows(images, systemWindows)
			window := windows.New(0, 0)
			// when
			window.Loop(func(frame *pixiq.Frame) {
				frame.CloseWindowEventually()
			})
			// then
			openWindow := systemWindows.openWindows[0]
			require.Len(t, openWindow.imagesDrawn, 1)
			assert.Same(t, openWindow.imagesDrawn[0], openWindow.visibleImage)
		})
		t.Run("after second frame", func(t *testing.T) {
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			systemWindows := &systemWindowsMock{}
			windows := pixiq.NewWindows(images, systemWindows)
			window := windows.New(0, 0)
			firstFrame := true
			// when
			window.Loop(func(frame *pixiq.Frame) {
				if !firstFrame {
					frame.CloseWindowEventually()
				}
				firstFrame = false
			})
			// then
			openWindow := systemWindows.openWindows[0]
			require.Len(t, openWindow.imagesDrawn, 2)
			assert.Same(t, openWindow.imagesDrawn[1], openWindow.visibleImage)
		})
	})

	t.Run("should close the window after loop is finished", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		systemWindows := &systemWindowsMock{}
		windows := pixiq.NewWindows(images, systemWindows)
		window := windows.New(0, 0)
		frameNumber := 0
		// when
		window.Loop(func(frame *pixiq.Frame) {
			if frameNumber == 2 {
				frame.CloseWindowEventually()
			} else {
				// then
				require.Len(t, systemWindows.openWindows, 1)
				assert.False(t, systemWindows.openWindows[0].closed)
				frameNumber += 1
			}
		})
		// then
		require.Len(t, systemWindows.openWindows, 1)
		assert.True(t, systemWindows.openWindows[0].closed)
	})

}

type systemWindowsMock struct {
	openWindows []*systemWindowMock
}

func (s *systemWindowsMock) Open(width, height int, hints ...pixiq.WindowHint) pixiq.SystemWindow {
	win := &systemWindowMock{width: width, height: height, hints: hints}
	s.openWindows = append(s.openWindows, win)
	return win
}

type systemWindowMock struct {
	imagesDrawn  []*pixiq.Image
	width        int
	height       int
	visibleImage *pixiq.Image
	closed       bool
	hints        []pixiq.WindowHint
}

func (f *systemWindowMock) Draw(image *pixiq.Image) {
	f.imagesDrawn = append(f.imagesDrawn, clone(image))
}

func (f *systemWindowMock) SwapImages() {
	f.visibleImage = f.imagesDrawn[len(f.imagesDrawn)-1]
}

func (f *systemWindowMock) Close() {
	f.closed = true
}

func clone(original *pixiq.Image) *pixiq.Image {
	images := pixiq.NewImages(&fakeAcceleratedImages{})
	clone := images.New(original.Width(), original.Height())
	originalSelection := original.WholeImageSelection()
	cloneSelection := clone.WholeImageSelection()
	for y := 0; y < originalSelection.Height(); y++ {
		for x := 0; x < originalSelection.Width(); x++ {
			cloneSelection.SetColor(x, y, originalSelection.Color(x, y))
		}
	}
	return clone
}
