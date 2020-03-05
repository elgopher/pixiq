package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
)

var transparent = image.RGBA(0, 0, 0, 0)

func TestNew(t *testing.T) {
	t.Run("should panic when AcceleratedImage is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			image.New(1, 1, nil)
		})
	})
	t.Run("should panic when width is less than 0", func(t *testing.T) {
		assert.Panics(t, func() {
			image.New(-1, 4, acceleratedImageStub{})
		})
	})
	t.Run("should panic when height is less than 0", func(t *testing.T) {
		assert.Panics(t, func() {
			image.New(2, -1, acceleratedImageStub{})
		})
	})
	t.Run("should create an image of given size", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
		}{
			"0x0": {
				width:  0,
				height: 0,
			},
			"0x1": {
				width:  0,
				height: 1,
			},
			"1x0": {
				width:  1,
				height: 0,
			},
			"1x1": {
				width:  1,
				height: 1,
			},
			"2x3": {
				width:  2,
				height: 3,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				img := image.New(test.width, test.height, acceleratedImageStub{})
				// then
				require.NotNil(t, img)
				assert.Equal(t, test.width, img.Width())
				assert.Equal(t, test.height, img.Height())
			})
		}
	})
}

func newImage(width, height int) *image.Image {
	return image.New(width, height, acceleratedImageStub{})
}

func TestImage_Selection(t *testing.T) {
	img := newImage(0, 0)

	t.Run("should create a selection for negative X", func(t *testing.T) {
		selection := img.Selection(-1, 0)
		assert.Equal(t, -1, selection.ImageX())
	})

	t.Run("should create a selection for negative Y", func(t *testing.T) {
		selection := img.Selection(0, -1)
		assert.Equal(t, -1, selection.ImageY())
	})

	t.Run("should create a selection", func(t *testing.T) {
		selection := img.Selection(1, 2)
		assert.Equal(t, 1, selection.ImageX())
		assert.Equal(t, 2, selection.ImageY())
		assert.Equal(t, 0, selection.Width())
		assert.Equal(t, 0, selection.Height())
		assert.Same(t, img, selection.Image())
	})
}

func TestSelection_Selection(t *testing.T) {
	img := newImage(0, 0)

	t.Run("should create a selection for negative x", func(t *testing.T) {
		selection := img.Selection(2, 0)
		subject := selection.Selection(-1, 0)
		assert.Equal(t, 1, subject.ImageX())
	})

	t.Run("should create a selection for negative y", func(t *testing.T) {
		selection := img.Selection(0, 2)
		subject := selection.Selection(0, -1)
		assert.Equal(t, 1, subject.ImageY())
	})

	t.Run("should create a selection out of selection", func(t *testing.T) {
		selection := img.Selection(1, 2)
		subject := selection.Selection(2, 3)
		assert.Equal(t, 3, subject.ImageX())
		assert.Equal(t, 5, subject.ImageY())
		assert.Equal(t, 0, subject.Width())
		assert.Equal(t, 0, subject.Height())
		assert.Same(t, img, subject.Image())
	})
}

func TestImage_WholeImageSelection(t *testing.T) {
	t.Run("should create a selection of whole image", func(t *testing.T) {
		tests := map[string]struct {
			image          *image.Image
			expectedWidth  int
			expectedHeight int
		}{
			"1": {
				image:          newImage(0, 0),
				expectedWidth:  0,
				expectedHeight: 0,
			},
			"2": {
				image:          newImage(3, 2),
				expectedWidth:  3,
				expectedHeight: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				selection := test.image.WholeImageSelection()
				// then
				assert.Equal(t, 0, selection.ImageX())
				assert.Equal(t, 0, selection.ImageY())
				assert.Equal(t, test.expectedWidth, selection.Width())
				assert.Equal(t, test.expectedHeight, selection.Height())
				assert.Same(t, test.image, selection.Image())
			})
		}
	})
}

func TestSelection_WithSize(t *testing.T) {
	img := newImage(0, 0)

	t.Run("should set selection width to zero if given width is negative", func(t *testing.T) {
		selection := img.Selection(1, 2)
		// when
		selection = selection.WithSize(-1, 4)
		assert.Equal(t, 0, selection.Width())
	})

	t.Run("should constrain width to zero if given width is negative and previously width was set to positive number", func(t *testing.T) {
		selection := img.Selection(1, 2).WithSize(5, 0)
		// when
		selection = selection.WithSize(-1, 4)
		assert.Equal(t, 0, selection.Width())
	})

	t.Run("should set selection height to zero if given height is negative", func(t *testing.T) {
		selection := img.Selection(1, 2)
		// when
		selection = selection.WithSize(3, -1)
		assert.Equal(t, 0, selection.Height())
	})

	t.Run("should constrain height to zero if given height is negative and previously height was set to positive number", func(t *testing.T) {
		selection := img.Selection(1, 2).WithSize(0, 5)
		// when
		selection = selection.WithSize(3, -1)
		assert.Equal(t, 0, selection.Height())
	})

	t.Run("should set selection size", func(t *testing.T) {
		selection := img.Selection(1, 2)
		// when
		selection = selection.WithSize(3, 4)
		assert.Equal(t, 3, selection.Width())
		assert.Equal(t, 4, selection.Height())
	})
}

func TestSelection_Color(t *testing.T) {
	t.Run("should return transparent color for pixel outside the image", func(t *testing.T) {
		tests := map[string]struct {
			image *image.Image
			x, y  int
		}{
			"1": {
				image: newImage(0, 0),
				x:     0,
				y:     0,
			},
			"2": {
				image: newImage(1, 1),
				x:     1,
				y:     0,
			},
			"3": {
				image: newImage(1, 1),
				x:     0,
				y:     1,
			},
			"4": {
				image: newImage(1, 1),
				x:     -1,
				y:     0,
			},
			"5": {
				image: newImage(1, 1),
				x:     0,
				y:     -1,
			},
			"6": {
				image: newImage(2, 2),
				x:     0,
				y:     2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				selection := test.image.Selection(0, 0)
				// when
				color := selection.Color(test.x, test.y)
				assert.Equal(t, transparent, color)
			})
		}
	})

	t.Run("should return color for pixel outside the selection and inside the image", func(t *testing.T) {
		color := image.RGBA(10, 20, 30, 40)

		t.Run("1", func(t *testing.T) {
			img := newImage(1, 1)
			img.Selection(0, 0).SetColor(0, 0, color)
			selection := img.Selection(1, 0)
			// expect
			assert.Equal(t, color, selection.Color(-1, 0))
		})

		t.Run("2", func(t *testing.T) {
			img := newImage(1, 1)
			img.Selection(0, 0).SetColor(0, 0, color)
			selection := img.Selection(-1, 0)
			// expect
			assert.Equal(t, color, selection.Color(1, 0))
		})

		t.Run("3", func(t *testing.T) {
			img := newImage(1, 1)
			img.Selection(0, 0).SetColor(0, 0, color)
			selection := img.Selection(0, 1)
			// expect
			assert.Equal(t, color, selection.Color(0, -1))
		})

		t.Run("4", func(t *testing.T) {
			img := newImage(1, 2)
			img.Selection(0, 0).SetColor(0, 1, color)
			selection := img.Selection(0, 1)
			// expect
			assert.Equal(t, color, selection.Color(0, 0))
		})
	})

}

func TestSelection_SetColor(t *testing.T) {
	color := image.RGBA(10, 20, 30, 40)

	t.Run("should set pixel color inside the image", func(t *testing.T) {

		t.Run("1", func(t *testing.T) {
			selection := newImage(1, 1).Selection(0, 0)
			// when
			selection.SetColor(0, 0, color)
			assert.Equal(t, color, selection.Color(0, 0))
		})

		t.Run("2", func(t *testing.T) {
			selection := newImage(2, 1).WholeImageSelection()
			// when
			selection.SetColor(1, 0, color)
			assertColors(t, selection, [][]image.Color{
				{transparent, color},
			})
		})

		t.Run("3", func(t *testing.T) {
			selection := newImage(1, 2).WholeImageSelection()
			// when
			selection.SetColor(0, 1, color)
			assertColors(t, selection, [][]image.Color{
				{transparent},
				{color},
			})
		})

		t.Run("4", func(t *testing.T) {
			selection := newImage(2, 2).WholeImageSelection()
			// when
			selection.SetColor(0, 1, color)
			assertColors(t, selection, [][]image.Color{
				{transparent, transparent},
				{color, transparent},
			})
		})
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {

		t.Run("1", func(t *testing.T) {
			selection := newImage(0, 0).Selection(0, 0)
			// when
			selection.SetColor(0, 0, color)
			assert.Equal(t, transparent, selection.Color(0, 0))
		})

		t.Run("2", func(t *testing.T) {
			selection := newImage(1, 1).Selection(0, 0)
			// when
			selection.SetColor(1, 0, color)
			assert.Equal(t, transparent, selection.Color(0, 0))
		})

		t.Run("3", func(t *testing.T) {
			selection := newImage(1, 1).Selection(0, 0)
			// when
			selection.SetColor(0, 1, color)
			assert.Equal(t, transparent, selection.Color(0, 0))
		})

		t.Run("4", func(t *testing.T) {
			selection := newImage(2, 1).WholeImageSelection()
			// when
			selection.SetColor(0, 1, color)
			assertColors(t, selection, [][]image.Color{
				{transparent, transparent},
			})
		})

		t.Run("5", func(t *testing.T) {
			selection := newImage(1, 2).WholeImageSelection()
			// when
			selection.SetColor(1, 0, color)
			assertColors(t, selection, [][]image.Color{
				{transparent},
				{transparent},
			})
		})

		t.Run("6", func(t *testing.T) {
			selection := newImage(1, 1).Selection(0, 0)
			// when
			selection.SetColor(-1, 0, color)
			assert.Equal(t, transparent, selection.Color(0, 0))
		})

		t.Run("7", func(t *testing.T) {
			selection := newImage(1, 1).Selection(0, 0)
			// when
			selection.SetColor(0, -1, color)
			assert.Equal(t, transparent, selection.Color(0, 0))
		})
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {

		t.Run("1", func(t *testing.T) {
			img := newImage(1, 1)
			selection := img.Selection(0, 0)
			// when
			selection.SetColor(0, 0, color)
			assert.Equal(t, color, img.Selection(0, 0).Color(0, 0))
		})

		t.Run("2", func(t *testing.T) {
			img := newImage(2, 1)
			selection := img.Selection(1, 0)
			// when
			selection.SetColor(0, 0, color)
			assertColors(t, img.WholeImageSelection(), [][]image.Color{
				{transparent, color},
			})
		})

		t.Run("3", func(t *testing.T) {
			img := newImage(2, 1)
			selection := img.Selection(1, 0)
			// when
			selection.SetColor(-1, 0, color)
			assertColors(t, img.WholeImageSelection(), [][]image.Color{
				{color, transparent},
			})
		})

		t.Run("4", func(t *testing.T) {
			img := newImage(1, 1)
			selection := img.Selection(-1, 0)
			// when
			selection.SetColor(1, 0, color)
			assert.Equal(t, color, img.Selection(0, 0).Color(0, 0))
		})

		t.Run("5", func(t *testing.T) {
			img := newImage(1, 2)
			selection := img.Selection(0, 1)
			// when
			selection.SetColor(0, 0, color)
			assertColors(t, img.WholeImageSelection(), [][]image.Color{
				{transparent},
				{color},
			})
		})

		t.Run("6", func(t *testing.T) {
			img := newImage(1, 2)
			selection := img.Selection(0, 1)
			// when
			selection.SetColor(0, -1, color)
			assertColors(t, img.WholeImageSelection(), [][]image.Color{
				{color},
				{transparent},
			})
		})
	})

	t.Run("SetColor overrides color set by AcceleratedCommand", func(t *testing.T) {
		var (
			color        = image.RGBA(10, 20, 30, 40)
			commandColor = image.RGBA(50, 60, 70, 80)
			accImg       = fake.NewAcceleratedImage(2, 1)
			img          = image.New(2, 1, accImg)
			selection    = img.Selection(0, 0)
		)
		selection.Modify(&acceleratedCommandMock{
			command: func(image.AcceleratedImageSelection, []image.AcceleratedImageSelection) {
				accImg.Upload([]image.Color{commandColor, commandColor})
			},
		})
		// when
		selection.SetColor(0, 0, color)
		// then
		assert.Equal(t, color, selection.Color(0, 0))
		// and
		assert.Equal(t, commandColor, selection.Color(1, 0))
	})

}

func TestImage_Upload(t *testing.T) {
	t.Run("should upload pixels", func(t *testing.T) {
		t.Run("0x0", func(t *testing.T) {
			acceleratedImage := fake.NewAcceleratedImage(0, 0)
			img := image.New(0, 0, acceleratedImage)
			// when
			img.Upload()
			// then
			assert.Equal(t, [][]image.Color{}, acceleratedImage.PixelsTable())
		})
		t.Run("1x1", func(t *testing.T) {
			acceleratedImage := fake.NewAcceleratedImage(1, 1)
			img := image.New(1, 1, acceleratedImage)
			color := image.RGBA(10, 20, 30, 40)
			img.Selection(0, 0).SetColor(0, 0, color)
			// when
			img.Upload()
			// then
			assert.Equal(t, [][]image.Color{{color}}, acceleratedImage.PixelsTable())
		})
		t.Run("2x2", func(t *testing.T) {
			acceleratedImage := fake.NewAcceleratedImage(2, 2)
			img := image.New(2, 2, acceleratedImage)
			color1 := image.RGBA(10, 20, 30, 40)
			color2 := image.RGBA(50, 50, 60, 70)
			color3 := image.RGBA(80, 90, 100, 110)
			color4 := image.RGBA(120, 130, 140, 150)
			selection := img.Selection(0, 0)
			selection.SetColor(0, 0, color1)
			selection.SetColor(1, 0, color2)
			selection.SetColor(0, 1, color3)
			selection.SetColor(1, 1, color4)
			// when
			img.Upload()
			// then
			table := [][]image.Color{
				{color3, color4},
				{color1, color2},
			}
			assert.Equal(t, table, acceleratedImage.PixelsTable())
		})
	})
	t.Run("Upload should not override colors set by AcceleratedCommand", func(t *testing.T) {
		var (
			color            = image.RGBA(50, 60, 70, 80)
			acceleratedImage = fake.NewAcceleratedImage(1, 1)
			img              = image.New(1, 1, acceleratedImage)
			selection        = img.Selection(0, 0)
		)
		selection.Modify(&acceleratedCommandMock{
			command: func(image.AcceleratedImageSelection, []image.AcceleratedImageSelection) {
				acceleratedImage.Upload([]image.Color{color})
			},
		})
		// when
		img.Upload()
		// then
		assert.Equal(t, [][]image.Color{{color}}, acceleratedImage.PixelsTable())
	})
}

func TestSelection_Modify(t *testing.T) {
	t.Run("should execute command", func(t *testing.T) {
		acceleratedImage := fake.NewAcceleratedImage(1, 1)
		img := image.New(1, 1, acceleratedImage)
		selection := img.WholeImageSelection()
		command := &acceleratedCommandMock{}
		// when
		selection.Modify(command)
		assert.Equal(t, 1, command.timesExecuted)
	})
	t.Run("should not do anything when command nil", func(t *testing.T) {
		acceleratedImage := fake.NewAcceleratedImage(1, 1)
		img := image.New(1, 1, acceleratedImage)
		selection := img.WholeImageSelection()
		selection.Modify(nil)
	})
	t.Run("should pass AcceleratedImageSelection to command.Run", func(t *testing.T) {
		tests := map[string]struct {
			x, y, width, height int
		}{
			"selection 0,1 with size 2x3":   {x: 0, y: 1, width: 2, height: 3},
			"selection -1,-2 with size 1x2": {x: -1, y: -2, width: 1, height: 2},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var (
					acceleratedImage = fake.NewAcceleratedImage(0, 0)
					img              = image.New(0, 0, acceleratedImage)
					selection        = img.
								Selection(test.x, test.y).
								WithSize(test.width, test.height)
					command = &acceleratedCommandMock{}
				)
				// when
				selection.Modify(command)
				// then
				assert.Equal(t, image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{
						X:      test.x,
						Y:      test.y,
						Width:  test.width,
						Height: test.height,
					},
					Image: acceleratedImage,
				}, command.output)
			})
		}
	})
	t.Run("should convert passed selections", func(t *testing.T) {
		var (
			acceleratedImage1 = fake.NewAcceleratedImage(0, 0)
			acceleratedImage2 = fake.NewAcceleratedImage(0, 0)
			img1              = image.New(0, 0, acceleratedImage1)
			img2              = image.New(0, 0, acceleratedImage2)
			command           = &acceleratedCommandMock{}
			output            = img1.WholeImageSelection()
		)
		tests := map[string]struct {
			selections []image.Selection
			expected   []image.AcceleratedImageSelection
		}{
			"no selections": {
				expected: []image.AcceleratedImageSelection{},
			},
			"1 selection": {
				selections: []image.Selection{
					img2.Selection(0, 1).WithSize(2, 3),
				},
				expected: []image.AcceleratedImageSelection{
					{
						Location: image.AcceleratedImageLocation{
							X:      0,
							Y:      1,
							Width:  2,
							Height: 3,
						},
						Image: acceleratedImage2,
					},
				},
			},
			"2 selections": {
				selections: []image.Selection{
					img1.Selection(1, 2).WithSize(3, 4),
					img2.Selection(5, 6).WithSize(7, 8),
				},
				expected: []image.AcceleratedImageSelection{
					{
						Location: image.AcceleratedImageLocation{
							X:      1,
							Y:      2,
							Width:  3,
							Height: 4,
						},
						Image: acceleratedImage1,
					},
					{
						Location: image.AcceleratedImageLocation{
							X:      5,
							Y:      6,
							Width:  7,
							Height: 8,
						},
						Image: acceleratedImage2,
					},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				output.Modify(command, test.selections...)
				// then
				assert.Equal(t, test.expected, command.selections)
			})
		}
	})
	t.Run("should convert selections in next call with different number of arguments", func(t *testing.T) {
		var (
			acceleratedImage1 = fake.NewAcceleratedImage(0, 0)
			acceleratedImage2 = fake.NewAcceleratedImage(0, 0)
			img1              = image.New(0, 0, acceleratedImage1)
			img2              = image.New(0, 0, acceleratedImage2)
			selection1        = img1.WholeImageSelection()
			selection2        = img2.WholeImageSelection()
			command           = &acceleratedCommandMock{}
			output            = img1.WholeImageSelection()
		)
		tests := map[string]struct {
			selectionsFirst  []image.Selection
			selectionsSecond []image.Selection
			expected         []image.AcceleratedImageSelection
		}{
			"1, then 0": {
				selectionsFirst:  []image.Selection{selection1},
				selectionsSecond: []image.Selection{},
				expected:         []image.AcceleratedImageSelection{},
			},
			"0, then 1": {
				selectionsFirst:  []image.Selection{},
				selectionsSecond: []image.Selection{selection1},
				expected:         []image.AcceleratedImageSelection{{Image: acceleratedImage1}},
			},
			"2, then 1": {
				selectionsFirst:  []image.Selection{selection1, selection2},
				selectionsSecond: []image.Selection{selection2},
				expected:         []image.AcceleratedImageSelection{{Image: acceleratedImage2}},
			},
			"1, then 2": {
				selectionsFirst:  []image.Selection{selection1},
				selectionsSecond: []image.Selection{selection2, selection1},
				expected: []image.AcceleratedImageSelection{
					{Image: acceleratedImage2},
					{Image: acceleratedImage1},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				output.Modify(command, test.selectionsFirst...)
				// when
				output.Modify(command, test.selectionsSecond...)
				// then
				assert.Equal(t, test.expected, command.selections)
			})
		}
	})
	t.Run("should upload passed selection", func(t *testing.T) {
		var (
			color00                = image.RGB(0, 0, 255)
			color10                = image.RGB(255, 0, 255)
			color01                = image.RGB(0, 255, 255)
			color11                = image.RGB(255, 255, 255)
			targetAcceleratedImage = fake.NewAcceleratedImage(1, 1)
			targetImage            = image.New(1, 1, targetAcceleratedImage)
			sourceAcceleratedImage = fake.NewAcceleratedImage(2, 2)
			sourceImage            = image.New(2, 2, sourceAcceleratedImage)
			uploadedPixels         = make([]image.Color, 4)
			command                = &acceleratedCommandMock{
				command: func(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) {
					source := selections[0].Image
					source.Download(uploadedPixels)
				},
			}
			outputSelection = targetImage.WholeImageSelection()
			sourceSelection = sourceImage.WholeImageSelection()
		)
		sourceSelection.SetColor(0, 0, color00)
		sourceSelection.SetColor(1, 0, color10)
		sourceSelection.SetColor(0, 1, color01)
		sourceSelection.SetColor(1, 1, color11)
		// when
		outputSelection.Modify(command, sourceSelection)
		// then
		assert.Equal(t, []image.Color{color01, color11, color00, color10}, uploadedPixels)
	})

	t.Run("should upload passed selections", func(t *testing.T) {
		var (
			color   = image.RGB(100, 100, 100)
			color00 = image.RGB(0, 0, 255)
			color10 = image.RGB(255, 0, 255)
			color01 = image.RGB(0, 255, 255)
			color11 = image.RGB(255, 255, 255)
			//
			targetAcceleratedImage = fake.NewAcceleratedImage(1, 1)
			targetImage            = image.New(1, 1, targetAcceleratedImage)
			//
			sourceAcceleratedImage0 = fake.NewAcceleratedImage(1, 1)
			sourceImage0            = image.New(1, 1, sourceAcceleratedImage0)
			uploadedPixels0         = make([]image.Color, 1)
			//
			sourceAcceleratedImage1 = fake.NewAcceleratedImage(2, 2)
			sourceImage1            = image.New(2, 2, sourceAcceleratedImage1)
			uploadedPixels1         = make([]image.Color, 4)
			//
			command = &acceleratedCommandMock{
				command: func(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) {
					selections[0].Image.Download(uploadedPixels0)
					selections[1].Image.Download(uploadedPixels1)
				},
			}
			outputSelection  = targetImage.WholeImageSelection()
			source0Selection = sourceImage0.WholeImageSelection()
			source1Selection = sourceImage1.WholeImageSelection()
		)
		source0Selection.SetColor(0, 0, color)
		source1Selection.SetColor(0, 0, color00)
		source1Selection.SetColor(1, 0, color10)
		source1Selection.SetColor(0, 1, color01)
		source1Selection.SetColor(1, 1, color11)
		// when
		outputSelection.Modify(command, source0Selection, source1Selection)
		// then
		assert.Equal(t, []image.Color{color}, uploadedPixels0)
		assert.Equal(t, []image.Color{color01, color11, color00, color10}, uploadedPixels1)
	})

	t.Run("should download pixels", func(t *testing.T) {
		var (
			color00 = image.RGB(0, 0, 255)
			color10 = image.RGB(255, 0, 255)
			color01 = image.RGB(0, 255, 255)
			color11 = image.RGB(255, 255, 255)
			//
			targetAcceleratedImage = fake.NewAcceleratedImage(2, 2)
			targetImage            = image.New(2, 2, targetAcceleratedImage)
			//
			command = &acceleratedCommandMock{
				command: func(image.AcceleratedImageSelection, []image.AcceleratedImageSelection) {
					targetAcceleratedImage.Upload([]image.Color{color01, color11, color00, color10})
				},
			}
			outputSelection = targetImage.WholeImageSelection()
		)
		// when
		outputSelection.Modify(command)
		//then
		assert.Equal(t, color00, outputSelection.Color(0, 0))
		assert.Equal(t, color10, outputSelection.Color(1, 0))
		assert.Equal(t, color01, outputSelection.Color(0, 1))
		assert.Equal(t, color11, outputSelection.Color(1, 1))
	})

	t.Run("should use results from last Modify", func(t *testing.T) {
		var (
			color     = image.RGBA(10, 20, 30, 40)
			accImg    = fake.NewAcceleratedImage(1, 1)
			img       = image.New(1, 1, accImg)
			selection = img.Selection(0, 0)
		)
		selection.Modify(&acceleratedCommandMock{
			command: func(image.AcceleratedImageSelection, []image.AcceleratedImageSelection) {
				accImg.Upload([]image.Color{color})
			},
		})
		// when
		selection.Modify(&acceleratedCommandStub{}, selection)
		// then
		assert.Equal(t, [][]image.Color{{color}}, accImg.PixelsTable())
	})
}

func TestLines_Length(t *testing.T) {
	t.Run("should return lines length", func(t *testing.T) {
		image0x0 := newImage(0, 0)
		image1x1 := newImage(1, 1)
		image1x2 := newImage(1, 2)
		tests := map[string]struct {
			image          *image.Image
			selection      image.Selection
			expectedLength int
		}{
			"image height 0": {
				image:          image0x0,
				selection:      image0x0.Selection(0, 0).WithSize(0, 1),
				expectedLength: 0,
			},
			"selection height 0": {
				image:          image1x1,
				selection:      image1x1.Selection(0, 0).WithSize(0, 0),
				expectedLength: 0,
			},
			"selection y 1, height 1": {
				image:          image1x1,
				selection:      image1x1.Selection(0, 1).WithSize(0, 1),
				expectedLength: 0,
			},
			"selection y -1, height 1": {
				image:          image1x1,
				selection:      image1x1.Selection(0, -1).WithSize(0, 1),
				expectedLength: 0,
			},
			"selection y -2, height 1": {
				image:          image1x1,
				selection:      image1x1.Selection(0, -2).WithSize(0, 1),
				expectedLength: 0,
			},
			"selection height 1": {
				image:          image1x1,
				selection:      image1x1.Selection(0, 0).WithSize(0, 1),
				expectedLength: 1,
			},
			"selection height 2": {
				image:          image1x1,
				selection:      image1x1.Selection(0, 0).WithSize(0, 2),
				expectedLength: 1,
			},
			"selection y -1, height 2": {
				image:          image1x1,
				selection:      image1x1.Selection(0, -1).WithSize(0, 2),
				expectedLength: 1,
			},
			"image height 2, selection y 1, height 1": {
				image:          image1x2,
				selection:      image1x2.Selection(0, 1).WithSize(0, 1),
				expectedLength: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				lines := test.selection.Lines()
				// when
				length := lines.Length()
				assert.Equal(t, test.expectedLength, length)
			})
		}
	})
}

func TestLines_YOffset(t *testing.T) {
	img := newImage(0, 0)
	tests := map[string]struct {
		selection       image.Selection
		expectedYOffset int
	}{
		"-1": {
			selection:       img.Selection(0, -1),
			expectedYOffset: 1,
		},
		"-2": {
			selection:       img.Selection(0, -2),
			expectedYOffset: 2,
		},
		"0": {
			selection:       img.Selection(0, 0),
			expectedYOffset: 0,
		},
		"1": {
			selection:       img.Selection(0, 1),
			expectedYOffset: 0,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lines := test.selection.Lines()
			// when
			offset := lines.YOffset()
			assert.Equal(t, test.expectedYOffset, offset)
		})
	}
}

func TestLines_XOffset(t *testing.T) {
	img := newImage(0, 0)
	tests := map[string]struct {
		selection       image.Selection
		expectedXOffset int
	}{
		"-1": {
			selection:       img.Selection(-1, 0),
			expectedXOffset: 1,
		},
		"-2": {
			selection:       img.Selection(-2, 0),
			expectedXOffset: 2,
		},
		"0": {
			selection:       img.Selection(0, 0),
			expectedXOffset: 0,
		},
		"1": {
			selection:       img.Selection(1, 0),
			expectedXOffset: 0,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lines := test.selection.Lines()
			// when
			offset := lines.XOffset()
			assert.Equal(t, test.expectedXOffset, offset)
		})
	}
}

func TestSelection_Lines(t *testing.T) {
	functions := map[string]func(image.Lines, int) []image.Color{
		"LineForRead":  image.Lines.LineForRead,
		"LineForWrite": image.Lines.LineForWrite,
	}
	for name, function := range functions {
		t.Run(name, func(t *testing.T) {
			t.Run("should panic when line is out-of-bounds the image", func(t *testing.T) {
				image0x0 := newImage(0, 0)
				image1x1 := newImage(1, 1)
				image1x2 := newImage(1, 2)
				tests := map[string]struct {
					line      int
					image     *image.Image
					selection image.Selection
				}{
					"selection height 0": {
						image:     image1x1,
						selection: image1x1.Selection(0, 0).WithSize(0, 0),
						line:      0,
					},
					"image height 0": {
						image:     image0x0,
						selection: image0x0.Selection(0, 0).WithSize(0, 1),
						line:      0,
					},
					"line negative": {
						image:     image1x1,
						selection: image1x1.Selection(0, 0).WithSize(0, 1),
						line:      -1,
					},
					"line equal to selection height": {
						image:     image1x1,
						selection: image1x1.Selection(0, 0).WithSize(0, 1),
						line:      1,
					},
					"line higher than selection height": {
						image:     image1x1,
						selection: image1x1.Selection(0, 0).WithSize(0, 1),
						line:      2,
					},
					"line higher than image height": {
						image:     image1x1,
						selection: image1x1.Selection(0, 1).WithSize(0, 1),
						line:      0,
					},
					"line above the image": {
						image:     image1x1,
						selection: image1x1.Selection(0, -1).WithSize(0, 1),
						line:      0,
					},
					"line higher than selection height but inside image": {
						image:     image1x2,
						selection: image1x2.Selection(0, 0).WithSize(0, 1),
						line:      1,
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						lines := test.selection.Lines()
						assert.Panics(t, func() {
							// when
							function(lines, test.line)
						})
					})
				}
			})
			t.Run("should return line", func(t *testing.T) {
				color1 := image.RGBA(10, 20, 30, 40)
				color2 := image.RGBA(50, 50, 60, 70)

				image1x1 := newImage(1, 1)
				image1x1.Selection(0, 0).SetColor(0, 0, color1)

				image1x2 := newImage(1, 2)
				image1x2Selection := image1x2.WholeImageSelection()
				image1x2Selection.SetColor(0, 0, color1)
				image1x2Selection.SetColor(0, 1, color2)

				image2x1 := newImage(2, 1)
				image2x1Selection := image2x1.WholeImageSelection()
				image2x1Selection.SetColor(0, 0, color1)
				image2x1Selection.SetColor(1, 0, color2)

				tests := map[string]struct {
					image     *image.Image
					selection image.Selection
					line      int
					expected  []image.Color
				}{
					"1": {
						image:     image1x1,
						selection: image1x1.Selection(0, 0).WithSize(1, 1),
						line:      0,
						expected:  []image.Color{color1},
					},
					"2": {
						image:     image1x2,
						selection: image1x2.Selection(0, 0).WithSize(1, 2),
						line:      1,
						expected:  []image.Color{color2},
					},
					"3": {
						image:     image1x2,
						selection: image1x2.Selection(0, 1).WithSize(1, 1),
						line:      0,
						expected:  []image.Color{color2},
					},
					"4": {
						image:     image1x2,
						selection: image1x2.Selection(0, 0).WithSize(1, 2),
						line:      0,
						expected:  []image.Color{color1},
					},
					"5": {
						image:     image1x1,
						selection: image1x1.Selection(1, 0).WithSize(1, 1),
						line:      0,
						expected:  []image.Color{},
					},
					"6": {
						image:     image1x1,
						selection: image1x1.Selection(-1, 0).WithSize(1, 1),
						line:      0,
						expected:  []image.Color{},
					},
					"7": {
						image:     image2x1,
						selection: image2x1.Selection(0, 0).WithSize(1, 1),
						line:      0,
						expected:  []image.Color{color1},
					},
					"8": {
						image:     image1x2,
						selection: image1x2.Selection(0, 0).WithSize(1, 2),
						line:      1,
						expected:  []image.Color{color2},
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						lines := test.selection.Lines()
						// when
						line := function(lines, test.line)
						// then
						require.NotNil(t, line)
						assert.Equal(t, test.expected, line)
					})
				}

			})
		})
	}

}

func assertColors(t *testing.T, selection image.Selection, expectedColorLines [][]image.Color) {
	assert.Equal(t, len(expectedColorLines), selection.Height(), "number of lines should be equal to selection height")
	for y := 0; y < selection.Height(); y++ {
		expectedColorLine := expectedColorLines[y]
		assert.Equal(t, len(expectedColorLine), selection.Width(), "number of pixels in a row should be equal to selection width")
		for x := 0; x < selection.Width(); x++ {
			color := selection.Color(x, y)
			assert.Equal(t, expectedColorLine[x], color, "position (%d,%d)", x, y)
		}
	}
}

type acceleratedImageStub struct{}

func (i acceleratedImageStub) Upload([]image.Color)   {}
func (i acceleratedImageStub) Download([]image.Color) {}

type acceleratedCommandMock struct {
	timesExecuted int
	output        image.AcceleratedImageSelection
	selections    []image.AcceleratedImageSelection
	command       func(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection)
}

func (a *acceleratedCommandMock) Run(output image.AcceleratedImageSelection, selections []image.AcceleratedImageSelection) {
	a.timesExecuted += 1
	a.output = output
	a.selections = selections
	if a.command != nil {
		a.command(output, selections)
	}
}
