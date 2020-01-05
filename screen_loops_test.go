package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

func TestNewScreenLoops(t *testing.T) {
	t.Run("should return ScreenLoops object", func(t *testing.T) {
		loops := pixiq.NewScreenLoops(pixiq.NewImages(&fakeAcceleratedImages{}))
		assert.NotNil(t, loops)
	})
}

func TestScreenLoops_Loop(t *testing.T) {

	t.Run("should run callback function until loop is stopped", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		loops := pixiq.NewScreenLoops(images)
		frameNumber := 0
		// when
		loops.Loop(&screenMock{}, func(frame *pixiq.Frame) {
			if frameNumber == 2 {
				frame.StopLoopEventually()
			} else {
				frameNumber += 1
			}
		})
		// then
		assert.Equal(t, 2, frameNumber)
	})

	t.Run("frame should provide screen", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		loops := pixiq.NewScreenLoops(images)
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
				// when
				loops.Loop(&screenMock{width: test.width, height: test.height}, func(frame *pixiq.Frame) {
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
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		screen := &screenMock{}
		loops := pixiq.NewScreenLoops(images)
		firstFrame := true
		var recordedImages []*pixiq.Image
		// when
		loops.Loop(screen, func(frame *pixiq.Frame) {
			if !firstFrame {
				frame.StopLoopEventually()
			}
			firstFrame = false
			recordedImages = append(recordedImages, frame.Screen().Image())
		})
		// then
		assert.Equal(t, recordedImages, screen.imagesDrawn)
	})

	t.Run("initial screen is transparent", func(t *testing.T) {
		images := pixiq.NewImages(&fakeAcceleratedImages{})
		loops := pixiq.NewScreenLoops(images)
		var screen pixiq.Selection
		// when
		loops.Loop(&screenMock{width: 1, height: 1}, func(frame *pixiq.Frame) {
			screen = frame.Screen()
			frame.StopLoopEventually()
		})
		// then
		assert.Equal(t, transparent, screen.Color(0, 0))
	})

	t.Run("should draw modified screen", func(t *testing.T) {
		t.Run("after first frame", func(t *testing.T) {
			color := pixiq.RGBA(10, 20, 30, 40)
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			screen := &screenMock{width: 1, height: 1}
			loops := pixiq.NewScreenLoops(images)
			// when
			loops.Loop(screen, func(frame *pixiq.Frame) {
				frame.Screen().SetColor(0, 0, color)
				frame.StopLoopEventually()
			})
			// then
			require.Len(t, screen.imagesDrawn, 1)
			assert.Equal(t, color, screen.imagesDrawn[0].WholeImageSelection().Color(0, 0))
		})
		t.Run("after second frame", func(t *testing.T) {
			color1 := pixiq.RGBA(10, 20, 30, 40)
			color2 := pixiq.RGBA(10, 20, 30, 40)
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			screen := &screenMock{width: 1, height: 1}
			loops := pixiq.NewScreenLoops(images)
			firstFrame := true
			// when
			loops.Loop(screen, func(frame *pixiq.Frame) {
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
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			screen := &screenMock{}
			loops := pixiq.NewScreenLoops(images)
			// when
			loops.Loop(screen, func(frame *pixiq.Frame) {
				frame.StopLoopEventually()
			})
			// then
			require.Len(t, screen.imagesDrawn, 1)
			assert.Same(t, screen.imagesDrawn[0], screen.visibleImage)
		})
		t.Run("after second frame", func(t *testing.T) {
			images := pixiq.NewImages(&fakeAcceleratedImages{})
			screen := &screenMock{}
			loops := pixiq.NewScreenLoops(images)
			firstFrame := true
			// when
			loops.Loop(screen, func(frame *pixiq.Frame) {
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
	imagesDrawn  []*pixiq.Image
	width        int
	height       int
	visibleImage *pixiq.Image
}

func (f *screenMock) Draw(image *pixiq.Image) {
	f.imagesDrawn = append(f.imagesDrawn, clone(image))
}

func (f *screenMock) SwapImages() {
	f.visibleImage = f.imagesDrawn[len(f.imagesDrawn)-1]
}

func (f *screenMock) Width() int {
	return f.width
}

func (f *screenMock) Height() int {
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
