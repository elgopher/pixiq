package loop_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
)

func TestRun(t *testing.T) {

	t.Run("should run callback function until loop is stopped", func(t *testing.T) {
		var frameNumber = 0
		// when
		screen := newScreenMock(1, 1)
		loop.Run(screen, func(frame *loop.Frame) {
			if frameNumber == 2 {
				frame.StopLoopEventually()
			} else {
				frameNumber++
			}
		})
		// then
		assert.Equal(t, 2, frameNumber)
	})

	t.Run("frame should provide screen", func(t *testing.T) {
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
				var screen image.Selection
				// when
				loop.Run(newScreenMock(test.width, test.height), func(frame *loop.Frame) {
					screen = frame.Screen()
					frame.StopLoopEventually()
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
		var (
			screen         = newScreenMock(1, 1)
			firstFrame     = true
			recordedImages []*image.Image
		)
		// when
		loop.Run(screen, func(frame *loop.Frame) {
			if !firstFrame {
				frame.StopLoopEventually()
			}
			firstFrame = false
			recordedImages = append(recordedImages, frame.Screen().Image())
		})
		// then
		assert.Equal(t, recordedImages, screen.imagesDrawn)
	})

	t.Run("should draw modified screen", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			var (
				color  = image.RGBA(10, 20, 30, 40)
				screen = newScreenMock(1, 1)
			)
			// when
			loop.Run(screen, func(frame *loop.Frame) {
				frame.Screen().SetColor(0, 0, color)
				frame.StopLoopEventually()
			})
			// then
			require.Len(t, screen.imagesDrawn, 1)
			assert.Equal(t, color, screen.imagesDrawn[0].WholeImageSelection().Color(0, 0))
		})
		t.Run("after second frame", func(t *testing.T) {
			var (
				color1     = image.RGBA(10, 20, 30, 40)
				color2     = image.RGBA(10, 20, 30, 40)
				screen     = newScreenMock(1, 1)
				firstFrame = true
			)
			// when
			loop.Run(screen, func(frame *loop.Frame) {
				if firstFrame {
					frame.Screen().SetColor(0, 0, color1)
					firstFrame = false
				} else {
					frame.Screen().SetColor(0, 0, color2)
					frame.StopLoopEventually()
				}
			})
			// then
			require.Len(t, screen.imagesDrawn, 2)
			assert.Equal(t, color1, screen.imagesDrawn[0].WholeImageSelection().Color(0, 0))
			assert.Equal(t, color2, screen.imagesDrawn[1].WholeImageSelection().Color(0, 0))
		})
	})

	t.Run("should swap images", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			screen := newScreenMock(1, 1)
			// when
			loop.Run(screen, func(frame *loop.Frame) {
				frame.StopLoopEventually()
			})
			// then
			require.Len(t, screen.imagesDrawn, 1)
			assert.Same(t, screen.imagesDrawn[0], screen.visibleImage)
		})
		t.Run("after second frame", func(t *testing.T) {
			var (
				screen     = newScreenMock(1, 1)
				firstFrame = true
			)
			// when
			loop.Run(screen, func(frame *loop.Frame) {
				if !firstFrame {
					frame.StopLoopEventually()
				}
				firstFrame = false
			})
			// then
			require.Len(t, screen.imagesDrawn, 2)
			assert.Same(t, screen.imagesDrawn[1], screen.visibleImage)
		})
	})

}

type screenMock struct {
	currentImage *image.Image
	imagesDrawn  []*image.Image
	width        int
	height       int
	visibleImage *image.Image
}

func newScreenMock(width, height int) *screenMock {
	return &screenMock{
		currentImage: image.New(width, height, &acceleratedImageStub{}),
		width:        width,
		height:       height,
	}
}

func (f *screenMock) Image() *image.Image {
	return f.currentImage
}

func (f *screenMock) Draw() {
	f.currentImage = clone(f.currentImage)
	f.imagesDrawn = append(f.imagesDrawn, f.currentImage)
}

func (f *screenMock) SwapImages() {
	f.visibleImage = f.currentImage
	f.currentImage = image.New(f.width, f.height, &acceleratedImageStub{})
}

func clone(original *image.Image) *image.Image {
	var (
		clone             = image.New(original.Width(), original.Height(), &acceleratedImageStub{})
		originalSelection = original.WholeImageSelection()
		cloneSelection    = clone.WholeImageSelection()
	)
	for y := 0; y < originalSelection.Height(); y++ {
		for x := 0; x < originalSelection.Width(); x++ {
			cloneSelection.SetColor(x, y, originalSelection.Color(x, y))
		}
	}
	return clone
}

type acceleratedImageStub struct{}

func (a acceleratedImageStub) Upload(selection image.AcceleratedSelection, pixels image.PixelSlice) {
}
func (a acceleratedImageStub) Download(selection image.AcceleratedSelection, pixels image.PixelSlice) {
}
func (a acceleratedImageStub) Modify(selection image.AcceleratedSelection, call image.AcceleratedCall) {
}
