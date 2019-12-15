package pixiq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq"
)

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
		// given
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
	image := images.New(4, 2)

	t.Run("should create a selection for negative x", func(t *testing.T) {
		selection := image.Selection(-1, 0)
		assert.Equal(t, selection.X(), -1)
	})

	t.Run("should create a selection for negative y", func(t *testing.T) {
		selection := image.Selection(0, -1)
		assert.Equal(t, selection.Y(), -1)
	})

	t.Run("should create a selection", func(t *testing.T) {
		selection := image.Selection(1, 2)
		assert.Equal(t, selection.X(), 1)
		assert.Equal(t, selection.Y(), 2)
		assert.Equal(t, selection.Width(), 0)
		assert.Equal(t, selection.Height(), 0)
	})
}

func TestImage_WholeImageSelection(t *testing.T) {
	t.Run("should create a selection of whole image", func(t *testing.T) {
		// given
		images := pixiq.NewImages()
		image := images.New(3, 2)
		// when
		selection := image.WholeImageSelection()
		// then
		assert.Equal(t, selection.X(), 0)
		assert.Equal(t, selection.Y(), 0)
		assert.Equal(t, selection.Width(), 3)
		assert.Equal(t, selection.Height(), 2)
	})
}

func TestSelection_WithSize(t *testing.T) {
	images := pixiq.NewImages()
	image := images.New(0, 0)

	t.Run("should set selection width to zero if given width is negative", func(t *testing.T) {
		selection := image.Selection(1, 2)
		// when
		selection = selection.WithSize(-1, 4)
		assert.Equal(t, selection.Width(), 0)
	})

	t.Run("should clamp width to zero if given width is negative and previously width was set to positive number", func(t *testing.T) {
		selection := image.Selection(1, 2).WithSize(5, 0)
		// when
		selection = selection.WithSize(-1, 4)
		assert.Equal(t, selection.Width(), 0)
	})

	t.Run("should set selection height to zero if given height is negative", func(t *testing.T) {
		selection := image.Selection(1, 2)
		// when
		selection = selection.WithSize(3, -1)
		assert.Equal(t, selection.Height(), 0)
	})

	t.Run("should clamp height to zero if given height is negative and previously height was set to positive number", func(t *testing.T) {
		selection := image.Selection(1, 2).WithSize(0, 5)
		// when
		selection = selection.WithSize(3, -1)
		assert.Equal(t, selection.Height(), 0)
	})

	t.Run("should set selection size", func(t *testing.T) {
		selection := image.Selection(1, 2)
		// when
		selection = selection.WithSize(3, 4)
		assert.Equal(t, selection.Width(), 3)
		assert.Equal(t, selection.Height(), 4)
	})
}
