package image_test

import (
	"errors"
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
			_, _ = image.New(1, 1, nil)
		})
	})
	t.Run("should return error when width is less than 0", func(t *testing.T) {
		img, err := image.New(-1, 4, acceleratedImageStub{})
		assert.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should return error when height is less than 0", func(t *testing.T) {
		img, err := image.New(2, -1, acceleratedImageStub{})
		assert.Error(t, err)
		assert.Nil(t, img)
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
				img, err := image.New(test.width, test.height, acceleratedImageStub{})
				// then
				require.NoError(t, err)
				require.NotNil(t, img)
				assert.Equal(t, test.width, img.Width())
				assert.Equal(t, test.height, img.Height())
			})
		}
	})
}

func newImage(width, height int) *image.Image {
	img, err := image.New(width, height, acceleratedImageStub{})
	if err != nil {
		panic(err)
	}
	return img
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
}

func TestImage_Upload(t *testing.T) {
	t.Run("should upload pixels", func(t *testing.T) {
		t.Run("0x0", func(t *testing.T) {
			acceleratedImage := newFakeAcceleratedImage()
			img, _ := image.New(0, 0, acceleratedImage)
			// when
			img.Upload()
			// then
			require.NotNil(t, acceleratedImage.pixels)
			assert.Len(t, acceleratedImage.pixels, 0)
		})
		t.Run("1x1", func(t *testing.T) {
			acceleratedImage := newFakeAcceleratedImage()
			img, _ := image.New(1, 1, acceleratedImage)
			color := image.RGBA(10, 20, 30, 40)
			img.Selection(0, 0).SetColor(0, 0, color)
			// when
			img.Upload()
			// then
			require.Len(t, acceleratedImage.pixels, 1)
			assert.Equal(t, color, acceleratedImage.pixels[0])
		})
		t.Run("2x2", func(t *testing.T) {
			acceleratedImage := newFakeAcceleratedImage()
			img, _ := image.New(2, 2, acceleratedImage)
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
			require.Len(t, acceleratedImage.pixels, 4)
			assert.Equal(t, color1, acceleratedImage.pixels[0])
			assert.Equal(t, color2, acceleratedImage.pixels[1])
			assert.Equal(t, color3, acceleratedImage.pixels[2])
			assert.Equal(t, color4, acceleratedImage.pixels[3])
		})
	})
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

func TestSelection_Modify(t *testing.T) {
	t.Run("should return error when acceleratedProgram is not given", func(t *testing.T) {
		var (
			accelerator = fake.NewAccelerator()
			img, _      = image.New(1, 1, accelerator.NewImage(1, 1))
			selection   = img.WholeImageSelection()
		)
		// when
		err := selection.Modify(nil, func(drawer image.Drawer) {})
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when procedure is not given", func(t *testing.T) {
		var (
			accelerator      = fake.NewAccelerator()
			acceleratedImage = accelerator.NewImage(1, 1)
			program          = fake.NewProgram()
			img, _           = image.New(1, 1, acceleratedImage)
			selection        = img.WholeImageSelection()
		)
		// when
		err := selection.Modify(program, nil)
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when program is not a fake", func(t *testing.T) {
		var (
			accelerator      = fake.NewAccelerator()
			acceleratedImage = accelerator.NewImage(1, 1)
			img, _           = image.New(1, 1, acceleratedImage)
			selection        = img.WholeImageSelection()
		)
		// when
		err := selection.Modify(struct{}{}, func(drawer image.Drawer) {})
		// then
		assert.Error(t, err)
	})
	t.Run("should execute procedure", func(t *testing.T) {
		var (
			executed         = false
			accelerator      = fake.NewAccelerator()
			acceleratedImage = accelerator.NewImage(1, 1)
			program          = fake.NewProgram()
			img, _           = image.New(1, 1, acceleratedImage)
			selection        = img.WholeImageSelection()
		)
		// when
		err := selection.Modify(program, func(drawer image.Drawer) {
			executed = true
		})
		// then
		require.NoError(t, err)
		assert.True(t, executed)
	})
	t.Run("should run program", func(t *testing.T) {
		tests := map[string]struct {
			x, y, width, height int
		}{
			"selection 0,0 with size 1x1": {
				x: 0, y: 0, width: 1, height: 1,
			},
			"selection 1,2 with size 3x4": {
				x: 1, y: 2, width: 3, height: 4,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var (
					accelerator      = fake.NewAccelerator()
					acceleratedImage = accelerator.NewImage(4, 4)
					program          = fake.NewProgram()
					img, _           = image.New(1, 1, acceleratedImage)
					selection        = img.Selection(test.x, test.y).WithSize(test.width, test.height)
				)
				// when
				err := selection.Modify(program, func(drawer image.Drawer) {})
				// then
				require.NoError(t, err)
				assert.True(t, program.Executed())
				expectedLocation := image.AcceleratedImageLocation{
					X:      test.x,
					Y:      test.y,
					Width:  test.width,
					Height: test.height,
				}
				assert.Equal(t, expectedLocation, program.TargetLocation())
				assert.Equal(t, acceleratedImage, program.TargetImage())
			})
		}
	})

	t.Run("should draw Primitive", func(t *testing.T) {
		var (
			accelerator      = fake.NewAccelerator()
			acceleratedImage = accelerator.NewImage(1, 1)
			program          = fake.NewProgram()
			img, _           = image.New(1, 1, acceleratedImage)
			selection        = img.WholeImageSelection()
			primitive        = &fake.Primitive{}
		)
		// when
		err := selection.Modify(program, func(drawer image.Drawer) {
			err := drawer.Draw(primitive)
			require.NoError(t, err)
		})
		// then
		require.NoError(t, err)
		assert.True(t, primitive.Drawn())
	})

}

type fakeAcceleratedImage struct {
	pixels []image.Color
}

func newFakeAcceleratedImage() *fakeAcceleratedImage {
	return &fakeAcceleratedImage{}
}

func (i fakeAcceleratedImage) Modify(p image.AcceleratedProgram, location image.AcceleratedImageLocation, procedure func(drawer image.AcceleratedDrawer)) error {
	return errors.New("unknown program")
}

func (i *fakeAcceleratedImage) Upload(pixels []image.Color) {
	i.pixels = make([]image.Color, len(pixels))
	// copy pixels to ensure that Upload method has been called
	copy(i.pixels, pixels)
}

func (i *fakeAcceleratedImage) Download(output []image.Color) {
	for j := 0; j < len(output); j++ {
		output[j] = i.pixels[j]
	}
}

type acceleratedImageStub struct{}

func (i acceleratedImageStub) Upload([]image.Color)   {}
func (i acceleratedImageStub) Download([]image.Color) {}
func (i acceleratedImageStub) Modify(_ image.AcceleratedProgram, _ image.AcceleratedImageLocation, procedure func(drawer image.AcceleratedDrawer)) error {
	procedure(acceleratedDrawerStub{})
	return nil
}

type acceleratedDrawerStub struct{}

func (a acceleratedDrawerStub) Draw(image.Primitive, ...interface{}) error {
	return nil
}

func (a acceleratedDrawerStub) SetSelection(string, image.AcceleratedImageSelection) {}
