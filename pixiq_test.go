package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

var transparent = pixiq.RGBA(0, 0, 0, 0)

func TestNewImages(t *testing.T) {
	t.Run("should create image factory", func(t *testing.T) {
		images := pixiq.NewImages()
		assert.NotNil(t, images)
	})
}

func TestImages_New(t *testing.T) {
	t.Run("should create an image of zero width when given width is less than 0", func(t *testing.T) {
		images := pixiq.NewImages()
		image := images.New(-1, 4)
		require.NotNil(t, image)
		assert.Equal(t, 0, image.Width())
	})
	t.Run("should create an image of zero height when given height is less than 0", func(t *testing.T) {
		images := pixiq.NewImages()
		image := images.New(2, -1)
		require.NotNil(t, image)
		assert.Equal(t, 0, image.Height())
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
		images := pixiq.NewImages()
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				image := images.New(test.width, test.height)
				// then
				require.NotNil(t, image)
				assert.Equal(t, test.width, image.Width())
				assert.Equal(t, test.height, image.Height())
			})
		}
	})
}

func TestImage_Selection(t *testing.T) {
	images := pixiq.NewImages()
	image := images.New(0, 0)

	t.Run("should create a selection for negative x", func(t *testing.T) {
		selection := image.Selection(-1, 0)
		assert.Equal(t, -1, selection.ImageX())
	})

	t.Run("should create a selection for negative y", func(t *testing.T) {
		selection := image.Selection(0, -1)
		assert.Equal(t, -1, selection.ImageY())
	})

	t.Run("should create a selection", func(t *testing.T) {
		selection := image.Selection(1, 2)
		assert.Equal(t, 1, selection.ImageX())
		assert.Equal(t, 2, selection.ImageY())
		assert.Equal(t, 0, selection.Width())
		assert.Equal(t, 0, selection.Height())
		assert.Same(t, image, selection.Image())
	})
}

func TestSelection_Selection(t *testing.T) {
	images := pixiq.NewImages()
	image := images.New(0, 0)

	t.Run("should create a selection for negative x", func(t *testing.T) {
		selection := image.Selection(2, 0)
		subject := selection.Selection(-1, 0)
		assert.Equal(t, 1, subject.ImageX())
	})

	t.Run("should create a selection for negative y", func(t *testing.T) {
		selection := image.Selection(0, 2)
		subject := selection.Selection(0, -1)
		assert.Equal(t, 1, subject.ImageY())
	})

	t.Run("should create a selection out of selection", func(t *testing.T) {
		selection := image.Selection(1, 2)
		subject := selection.Selection(2, 3)
		assert.Equal(t, 3, subject.ImageX())
		assert.Equal(t, 5, subject.ImageY())
		assert.Equal(t, 0, subject.Width())
		assert.Equal(t, 0, subject.Height())
		assert.Same(t, image, subject.Image())
	})
}

func TestImage_WholeImageSelection(t *testing.T) {
	t.Run("should create a selection of whole image", func(t *testing.T) {
		images := pixiq.NewImages()
		image := images.New(3, 2)
		// when
		selection := image.WholeImageSelection()
		// then
		assert.Equal(t, 0, selection.ImageX())
		assert.Equal(t, 0, selection.ImageY())
		assert.Equal(t, 3, selection.Width())
		assert.Equal(t, 2, selection.Height())
		assert.Same(t, image, selection.Image())
	})
}

func TestSelection_WithSize(t *testing.T) {
	images := pixiq.NewImages()
	image := images.New(0, 0)

	t.Run("should set selection width to zero if given width is negative", func(t *testing.T) {
		selection := image.Selection(1, 2)
		// when
		selection = selection.WithSize(-1, 4)
		assert.Equal(t, 0, selection.Width())
	})

	t.Run("should clamp width to zero if given width is negative and previously width was set to positive number", func(t *testing.T) {
		selection := image.Selection(1, 2).WithSize(5, 0)
		// when
		selection = selection.WithSize(-1, 4)
		assert.Equal(t, 0, selection.Width())
	})

	t.Run("should set selection height to zero if given height is negative", func(t *testing.T) {
		selection := image.Selection(1, 2)
		// when
		selection = selection.WithSize(3, -1)
		assert.Equal(t, 0, selection.Height())
	})

	t.Run("should clamp height to zero if given height is negative and previously height was set to positive number", func(t *testing.T) {
		selection := image.Selection(1, 2).WithSize(0, 5)
		// when
		selection = selection.WithSize(3, -1)
		assert.Equal(t, 0, selection.Height())
	})

	t.Run("should set selection size", func(t *testing.T) {
		selection := image.Selection(1, 2)
		// when
		selection = selection.WithSize(3, 4)
		assert.Equal(t, 3, selection.Width())
		assert.Equal(t, 4, selection.Height())
	})
}

func TestSelection_Color(t *testing.T) {
	images := pixiq.NewImages()

	t.Run("should return transparent color for pixel outside the image", func(t *testing.T) {
		tests := map[string]struct {
			image *pixiq.Image
			x, y  int
		}{
			"1": {
				image: images.New(0, 0),
				x:     0,
				y:     0,
			},
			"2": {
				image: images.New(1, 1),
				x:     1,
				y:     0,
			},
			"3": {
				image: images.New(1, 1),
				x:     0,
				y:     1,
			},
			"4": {
				image: images.New(1, 1),
				x:     -1,
				y:     0,
			},
			"5": {
				image: images.New(1, 1),
				x:     0,
				y:     -1,
			},
			"6": {
				image: images.New(2, 2),
				x:     0,
				y:     2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				selection := test.image.WholeImageSelection()
				// when
				color := selection.Color(test.x, test.y)
				assert.Equal(t, transparent, color)
			})
		}
	})

	t.Run("should return color for pixel outside the selection and inside the image", func(t *testing.T) {
		color := pixiq.RGBA(10, 20, 30, 40)

		t.Run("1", func(t *testing.T) {
			image := images.New(1, 1)
			image.WholeImageSelection().SetColor(0, 0, color)
			selection := image.Selection(1, 0)
			// expect
			assert.Equal(t, color, selection.Color(-1, 0))
		})

		t.Run("2", func(t *testing.T) {
			image := images.New(1, 1)
			image.WholeImageSelection().SetColor(0, 0, color)
			selection := image.Selection(-1, 0)
			// expect
			assert.Equal(t, color, selection.Color(1, 0))
		})

		t.Run("3", func(t *testing.T) {
			image := images.New(1, 1)
			image.WholeImageSelection().SetColor(0, 0, color)
			selection := image.Selection(0, 1)
			// expect
			assert.Equal(t, color, selection.Color(0, -1))
		})

		t.Run("4", func(t *testing.T) {
			image := images.New(1, 2)
			image.WholeImageSelection().SetColor(0, 1, color)
			selection := image.Selection(0, 1)
			// expect
			assert.Equal(t, color, selection.Color(0, 0))
		})
	})

}

func TestSelection_SetColor(t *testing.T) {
	images := pixiq.NewImages()
	color := pixiq.RGBA(10, 20, 30, 40)

	t.Run("should set pixel color inside the image", func(t *testing.T) {
		selection := images.New(1, 1).WholeImageSelection()
		// when
		selection.SetColor(0, 0, color)
		assert.Equal(t, color, selection.Color(0, 0))
	})

	t.Run("should set pixel color inside the image without modifying the other", func(t *testing.T) {
		selection := images.New(2, 1).WholeImageSelection()
		// when
		selection.SetColor(1, 0, color)
		assertColors(t, selection, [][]pixiq.Color{
			{transparent, color},
		})
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(0, 0).WholeImageSelection()
		// when
		selection.SetColor(0, 0, color)
		assert.Equal(t, transparent, selection.Color(0, 0))
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(1, 1).WholeImageSelection()
		// when
		selection.SetColor(1, 0, color)
		assert.Equal(t, transparent, selection.Color(0, 0))
	})

	t.Run("should set pixel color inside the image without modifying the other", func(t *testing.T) {
		selection := images.New(1, 2).WholeImageSelection()
		// when
		selection.SetColor(0, 1, color)
		assertColors(t, selection, [][]pixiq.Color{
			{transparent},
			{color},
		})
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(1, 1).WholeImageSelection()
		// when
		selection.SetColor(0, 1, color)
		assert.Equal(t, transparent, selection.Color(0, 0))
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(2, 1).WholeImageSelection()
		// when
		selection.SetColor(0, 1, color)
		assertColors(t, selection, [][]pixiq.Color{
			{transparent, transparent},
		})
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(1, 2).WholeImageSelection()
		// when
		selection.SetColor(1, 0, color)
		assertColors(t, selection, [][]pixiq.Color{
			{transparent},
			{transparent},
		})
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(1, 1).WholeImageSelection()
		// when
		selection.SetColor(-1, 0, color)
		assert.Equal(t, transparent, selection.Color(0, 0))
	})

	t.Run("setting pixel color outside the image does nothing", func(t *testing.T) {
		selection := images.New(1, 1).WholeImageSelection()
		// when
		selection.SetColor(0, -1, color)
		assert.Equal(t, transparent, selection.Color(0, 0))
	})

	t.Run("should set pixel color inside the image without modifying the other", func(t *testing.T) {
		selection := images.New(2, 2).WholeImageSelection()
		// when
		selection.SetColor(0, 1, color)
		assertColors(t, selection, [][]pixiq.Color{
			{transparent, transparent},
			{color, transparent},
		})
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {
		image := images.New(1, 1)
		selection := image.Selection(0, 0)
		// when
		selection.SetColor(0, 0, color)
		assert.Equal(t, color, image.WholeImageSelection().Color(0, 0))
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {
		image := images.New(2, 1)
		selection := image.Selection(1, 0)
		// when
		selection.SetColor(0, 0, color)
		assertColors(t, image.WholeImageSelection(), [][]pixiq.Color{
			{transparent, color},
		})
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {
		image := images.New(2, 1)
		selection := image.Selection(1, 0)
		// when
		selection.SetColor(-1, 0, color)
		assertColors(t, image.WholeImageSelection(), [][]pixiq.Color{
			{color, transparent},
		})
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {
		image := images.New(1, 1)
		selection := image.Selection(-1, 0)
		// when
		selection.SetColor(1, 0, color)
		assert.Equal(t, color, image.WholeImageSelection().Color(0, 0))
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {
		image := images.New(1, 2)
		selection := image.Selection(0, 1)
		// when
		selection.SetColor(0, 0, color)
		assertColors(t, image.WholeImageSelection(), [][]pixiq.Color{
			{transparent},
			{color},
		})
	})

	t.Run("should set pixel color outside the selection", func(t *testing.T) {
		image := images.New(1, 2)
		selection := image.Selection(0, 1)
		// when
		selection.SetColor(0, -1, color)
		assertColors(t, image.WholeImageSelection(), [][]pixiq.Color{
			{color},
			{transparent},
		})
	})

}

func assertColors(t *testing.T, selection pixiq.Selection, expectedColorLines [][]pixiq.Color) {
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
