package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

func TestNewWindows(t *testing.T) {
	t.Run("should return Windows object", func(t *testing.T) {
		windows := pixiq.NewWindows(pixiq.NewImages(&fakeAcceleratedImages{}))
		assert.NotNil(t, windows)
	})
}

func TestWindow_Loop(t *testing.T) {

	t.Run("should run callback function until window is closed", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		windows := pixiq.NewWindows(images)
		frameNumber := 0
		// when
		windows.Loop(&windowMock{}, func(frame *pixiq.Frame) {
			if frameNumber == 2 {
				frame.CloseWindowEventually()
			} else {
				frameNumber += 1
			}
		})
		// then
		assert.Equal(t, 2, frameNumber)
	})

	t.Run("frame should provide window's screen", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		windows := pixiq.NewWindows(images)
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
				var screen pixiq.Selection
				win := &windowMock{width: test.width, height: test.height}
				// when
				windows.Loop(win, func(frame *pixiq.Frame) {
					screen = frame.Screen()
					frame.CloseWindowEventually()
				})
				// then
				assert.Equal(t, test.width, screen.Width())
				assert.Equal(t, test.height, screen.Height())
				assert.Equal(t, 0, screen.ImageX())
				assert.Equal(t, 0, screen.ImageY())
				assert.NotNil(t, screen.Image())
				assert.Equal(t, test.width, screen.Image().Width())
				assert.Equal(t, test.height, screen.Image().Height())
			})
		}
	})

	t.Run("should draw image for each frame", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		win := &windowMock{}
		windows := pixiq.NewWindows(images)
		firstFrame := true
		var recordedImages []*pixiq.Image
		// when
		windows.Loop(win, func(frame *pixiq.Frame) {
			if !firstFrame {
				frame.CloseWindowEventually()
			}
			firstFrame = false
			recordedImages = append(recordedImages, frame.Screen().Image())
		})
		// then
		assert.Equal(t, recordedImages, win.imagesDrawn)
	})

	t.Run("initial screen is transparent", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		windows := pixiq.NewWindows(images)
		win := &windowMock{width: 1, height: 1}
		var screen pixiq.Selection
		// when
		windows.Loop(win, func(frame *pixiq.Frame) {
			screen = frame.Screen()
			frame.CloseWindowEventually()
		})
		// then
		assert.Equal(t, transparent, screen.Color(0, 0))
	})

	t.Run("should draw modified window image", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			color := pixiq.RGBA(10, 20, 30, 40)
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			win := &windowMock{width: 1, height: 1}
			windows := pixiq.NewWindows(images)
			// when
			windows.Loop(win, func(frame *pixiq.Frame) {
				frame.Screen().SetColor(0, 0, color)
				frame.CloseWindowEventually()
			})
			// then
			require.Len(t, win.imagesDrawn, 1)
			assert.Equal(t, color, win.imagesDrawn[0].WholeImageSelection().Color(0, 0))
		})
		t.Run("after second frame", func(t *testing.T) {
			color1 := pixiq.RGBA(10, 20, 30, 40)
			color2 := pixiq.RGBA(10, 20, 30, 40)
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			win := &windowMock{width: 1, height: 1}
			windows := pixiq.NewWindows(images)
			firstFrame := true
			// when
			windows.Loop(win, func(frame *pixiq.Frame) {
				if firstFrame {
					frame.Screen().SetColor(0, 0, color1)
					firstFrame = false
				} else {
					frame.Screen().SetColor(0, 0, color2)
					frame.CloseWindowEventually()
				}
			})
			// then
			require.Len(t, win.imagesDrawn, 2)
			assert.Equal(t, color1, win.imagesDrawn[0].WholeImageSelection().Color(0, 0))
			assert.Equal(t, color2, win.imagesDrawn[1].WholeImageSelection().Color(0, 0))
		})
	})

	t.Run("should swap images", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			openWindow := &windowMock{}
			windows := pixiq.NewWindows(images)
			// when
			windows.Loop(openWindow, func(frame *pixiq.Frame) {
				frame.CloseWindowEventually()
			})
			// then
			require.Len(t, openWindow.imagesDrawn, 1)
			assert.Same(t, openWindow.imagesDrawn[0], openWindow.visibleImage)
		})
		t.Run("after second frame", func(t *testing.T) {
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			openWindow := &windowMock{}
			windows := pixiq.NewWindows(images)
			firstFrame := true
			// when
			windows.Loop(openWindow, func(frame *pixiq.Frame) {
				if !firstFrame {
					frame.CloseWindowEventually()
				}
				firstFrame = false
			})
			// then
			require.Len(t, openWindow.imagesDrawn, 2)
			assert.Same(t, openWindow.imagesDrawn[1], openWindow.visibleImage)
		})
	})

	t.Run("should close the window after loop is finished", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		win := &windowMock{}
		windows := pixiq.NewWindows(images)
		frameNumber := 0
		// when
		windows.Loop(win, func(frame *pixiq.Frame) {
			if frameNumber == 2 {
				frame.CloseWindowEventually()
			} else {
				// then
				assert.False(t, win.closed)
				frameNumber += 1
			}
		})
		// then
		assert.True(t, win.closed)
	})

}

type windowMock struct {
	imagesDrawn  []*pixiq.Image
	width        int
	height       int
	visibleImage *pixiq.Image
	closed       bool
}

func (f *windowMock) Draw(image *pixiq.Image) {
	f.imagesDrawn = append(f.imagesDrawn, clone(image))
}

func (f *windowMock) SwapImages() {
	f.visibleImage = f.imagesDrawn[len(f.imagesDrawn)-1]
}

func (f *windowMock) Close() {
	f.closed = true
}

func (f *windowMock) Width() int {
	return f.width
}

func (f *windowMock) Height() int {
	return f.height
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
