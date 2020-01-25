package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/image"
)

var transparent = image.RGBA(0, 0, 0, 0)

func TestNew(t *testing.T) {
	t.Run("should return error when AcceleratedImage is nil", func(t *testing.T) {
		img, err := image.New(1, 1, nil)
		assert.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should return error when width is less than 0", func(t *testing.T) {
		img, err := image.New(-1, 4, newFakeAcceleratedImage(0, 0))
		assert.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should return error when height is less than 0", func(t *testing.T) {
		img, err := image.New(2, -1, newFakeAcceleratedImage(0, 0))
		assert.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should create an image of any size", func(t *testing.T) {
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
				img, err := image.New(test.width, test.height, newFakeAcceleratedImage(0, 0))
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
	img, err := image.New(width, height, newFakeAcceleratedImage(0, 0))
	if err != nil {
		panic(err)
	}
	return img
}

func TestImage_Selection(t *testing.T) {
	img := newImage(0, 0)

	t.Run("should create a selection for negative x", func(t *testing.T) {
		selection := img.Selection(-1, 0)
		assert.Equal(t, -1, selection.ImageX())
	})

	t.Run("should create a selection for negative y", func(t *testing.T) {
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

//func TestImage_Upload(t *testing.T) {
//	t.Run("should upload pixels", func(t *testing.T) {
//		t.Run("0x0", func(t *testing.T) {
//			acceleratedImage := fakeAcceleratedImage()
//			img := image.New(0, 0, acceleratedImage)
//			// when
//			img.Upload()
//			// then
//			require.NotNil(t, acceleratedImage.pixels)
//			assert.Len(t, acceleratedImage.pixels, 0)
//		})
//		t.Run("1x1", func(t *testing.T) {
//			acceleratedImage := fakeAcceleratedImage()
//			img := image.New(1, 1, acceleratedImage)
//			color := image.RGBA(10, 20, 30, 40)
//			img.Selection(0, 0).SetColor(0, 0, color)
//			// when
//			img.Upload()
//			// then
//			require.Len(t, acceleratedImage.pixels, 1)
//			assert.Equal(t, color, acceleratedImage.pixels[0])
//		})
//		t.Run("2x2", func(t *testing.T) {
//			acceleratedImage := fakeAcceleratedImage()
//			img := image.New(2, 2, acceleratedImage)
//			color1 := image.RGBA(10, 20, 30, 40)
//			color2 := image.RGBA(50, 50, 60, 70)
//			color3 := image.RGBA(80, 90, 100, 110)
//			color4 := image.RGBA(120, 130, 140, 150)
//			selection := img.Selection(0, 0)
//			selection.SetColor(0, 0, color1)
//			selection.SetColor(1, 0, color2)
//			selection.SetColor(0, 1, color3)
//			selection.SetColor(1, 1, color4)
//			// when
//			img.Upload()
//			// then
//			require.Len(t, acceleratedImage.pixels, 4)
//			assert.Equal(t, color1, acceleratedImage.pixels[0])
//			assert.Equal(t, color2, acceleratedImage.pixels[1])
//			assert.Equal(t, color3, acceleratedImage.pixels[2])
//			assert.Equal(t, color4, acceleratedImage.pixels[3])
//		})
//	})
//}

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
	t.Run("should run AcceleratedCall on AcceleratedFragment", func(t *testing.T) {
		var (
			colorToAdd = image.RGBA(1, 2, 3, 4)
			color1     = image.RGBA(10, 20, 30, 40)
			modif1     = image.RGBA(11, 22, 33, 44)
			color2     = image.RGBA(50, 60, 70, 80)
			modif2     = image.RGBA(51, 62, 73, 84)
			color3     = image.RGBA(90, 100, 110, 120)
			//modif3     = image.RGBA(91, 102, 113, 124)
			color4 = image.RGBA(130, 140, 150, 160)
			modif4 = image.RGBA(131, 142, 153, 164)
		)
		tests := map[string]struct {
			imageWidth, imageHeight         int
			selectionX, selectionY          int
			selectionWidth, selectionHeight int
			givenColors                     [][]image.Color
			expectedColors                  [][]image.Color
		}{
			"image 1x1, selection 0,0 with size 1x1": {
				imageWidth: 1, imageHeight: 1,
				selectionWidth: 1, selectionHeight: 1,
				givenColors:    [][]image.Color{{color1}},
				expectedColors: [][]image.Color{{modif1}},
			},
			"image 1x1, selection 1,0 with size 0x1": {
				imageWidth: 1, imageHeight: 1,
				selectionX: 1, selectionHeight: 1,
				givenColors:    [][]image.Color{{color1}},
				expectedColors: [][]image.Color{{color1}},
			},
			"image 1x1, selection 0,1 with size 1x0": {
				imageWidth: 1, imageHeight: 1,
				selectionY: 1, selectionWidth: 1,
				givenColors:    [][]image.Color{{color1}},
				expectedColors: [][]image.Color{{color1}},
			},
			"image 1x1, selection 0,0 with size 1x0": {
				imageWidth: 1, imageHeight: 1,
				selectionWidth: 1,
				givenColors:    [][]image.Color{{color1}},
				expectedColors: [][]image.Color{{color1}},
			},
			"image 1x1, selection 0,0 with size 0x1": {
				imageWidth: 1, imageHeight: 1,
				selectionHeight: 1,
				givenColors:     [][]image.Color{{color1}},
				expectedColors:  [][]image.Color{{color1}},
			},
			"image 1x1, selection 0,0 with size 2x1": {
				imageWidth: 1, imageHeight: 1,
				selectionWidth: 2, selectionHeight: 1,
				givenColors:    [][]image.Color{{color1}},
				expectedColors: [][]image.Color{{modif1}},
			},
			"image 1x1, selection 0,0 with size 1x2": {
				imageWidth: 1, imageHeight: 1,
				selectionWidth: 1, selectionHeight: 2,
				givenColors:    [][]image.Color{{color1}},
				expectedColors: [][]image.Color{{modif1}},
			},
			"image 2x1, selection 1,0 with size 2x1": {
				imageWidth: 2, imageHeight: 1,
				selectionX:     1,
				selectionWidth: 2, selectionHeight: 1,
				givenColors:    [][]image.Color{{color1, color2}},
				expectedColors: [][]image.Color{{color1, modif2}},
			},
			"image 1x2, selection 0,1 with size 1x2": {
				imageWidth: 1, imageHeight: 2,
				selectionY:     1,
				selectionWidth: 1, selectionHeight: 2,
				givenColors: [][]image.Color{
					{color1},
					{color2},
				},
				expectedColors: [][]image.Color{
					{color1},
					{modif2},
				},
			},
			"image 2x2, selection 1,1 with size 1x1": {
				imageWidth: 2, imageHeight: 2,
				selectionX: 1, selectionY: 1,
				selectionWidth: 1, selectionHeight: 1,
				givenColors: [][]image.Color{
					{color1, color2},
					{color3, color4},
				},
				expectedColors: [][]image.Color{
					{color1, color2},
					{color3, modif4},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var (
					images      = image.NewFakeImages()
					accImage, _ = images.NewAcceleratedImage(test.imageWidth, test.imageHeight)
					call        = images.AddColor(colorToAdd)
					img, _      = image.New(test.imageWidth, test.imageHeight, accImage)
					selection   = img.Selection(test.selectionX, test.selectionY).
							WithSize(test.selectionWidth, test.selectionHeight)
				)
				wholeImage := img.WholeImageSelection()
				for y, row := range test.givenColors {
					for x, color := range row {
						wholeImage.SetColor(x, y, color)
					}
				}
				// when
				selection.Modify(call)
				// then
				actual := make([][]image.Color, len(test.givenColors))
				for y, row := range test.givenColors {
					actual[y] = make([]image.Color, len(row))
					for x := range row {
						actual[y][x] = wholeImage.Color(x, y)
					}
				}
				assert.Equal(t, test.expectedColors, actual)
			})
		}
	})
	t.Run("location should be a selection clamped to image boundaries", func(t *testing.T) {
		tests := map[string]struct {
			imageWidth, imageHeight         int
			selectionX, selectionY          int
			selectionWidth, selectionHeight int
		}{
			"image 0x0, selection -1,0 with size 0x0": {
				selectionX: -1,
			},
			"image 0x0, selection -2,0 with size 0x0": {
				selectionX: -2,
			},
			"image 1x1, selection -1,0 with size 0x0": {
				imageWidth: 1, imageHeight: 1,
				selectionX: -1,
			},
			"image 0x0, selection 0,-1 with size 0x0": {
				selectionY: -1,
			},
			"image 0x0, selection 0,-2 with size 0x0": {
				selectionY: -2,
			},
			"image 1x1, selection 0,-1 with size 0x0": {
				imageWidth: 1, imageHeight: 1,
				selectionY: -1,
			},
			"image 0x0, selection 0,0 with size 1x0": {
				selectionWidth: 1,
			},
			"image 0x0, selection 0,0 with size 2x0": {
				selectionWidth: 2,
			},
			"image 1x1, selection 0,0 with size 2x0": {
				imageWidth: 1, imageHeight: 1,
				selectionWidth: 2,
			},
			"image 0x0, selection 0,0 with size 0x1": {
				selectionHeight: 1,
			},
			"image 0x0, selection 0,0 with size 0x2": {
				selectionHeight: 2,
			},
			"image 1x1, selection 0,0 with size 0x2": {
				imageWidth: 1, imageHeight: 1,
				selectionHeight: 2,
			},
			"image 1x1, selection 1,0 with size 1x0": {
				imageWidth: 1, imageHeight: 1,
				selectionX:     1,
				selectionWidth: 1,
			},
			"image 2x1, selection 1,0 with size 2x0": {
				imageWidth: 2, imageHeight: 1,
				selectionX:     1,
				selectionWidth: 2,
			},
			"image 1x1, selection 2,0 with size 0x0": {
				imageWidth: 1, imageHeight: 1,
				selectionX: 2,
			},
			"image 1x1, selection 0,1 with size 0x1": {
				imageWidth: 1, imageHeight: 1,
				selectionY:      1,
				selectionHeight: 1,
			},
			"image 1x2, selection 0,1 with size 0x2": {
				imageWidth: 1, imageHeight: 2,
				selectionY:      1,
				selectionHeight: 2,
			},
			"image 1x1, selection 0,2 with size 0x0": {
				imageWidth: 1, imageHeight: 1,
				selectionY: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var (
					fakes       = image.NewFakeImages()
					accImage, _ = fakes.NewAcceleratedImage(test.imageWidth, test.imageHeight)
				)
				img, _ := image.New(test.imageWidth, test.imageHeight, accImage)
				var (
					call      = &assertLocationClampedCall{t: t, image: img}
					selection = img.Selection(test.selectionX, test.selectionY).
							WithSize(test.selectionWidth, test.selectionHeight)
				)
				fakes.RegisterCall(call)
				// when
				selection.Modify(call)
			})
		}
	})

}

type assertLocationClampedCall struct {
	image *image.Image
	t     *testing.T
}

func (a *assertLocationClampedCall) Run(selection image.AcceleratedFragmentLocation, _ *image.FakeAcceleratedImage) {
	assert.True(a.t, selection.X >= 0, "x>=0")
	assert.True(a.t, selection.Y >= 0, "y>=0")
	assert.True(a.t, selection.Width >= 0, "width>=0")
	assert.True(a.t, selection.Height >= 0, "height>=0")
	assert.True(a.t, selection.X+selection.Width <= a.image.Width(), "x+width<image.width")
	assert.True(a.t, selection.Y+selection.Height <= a.image.Height(), "y+height<image.height")
}
