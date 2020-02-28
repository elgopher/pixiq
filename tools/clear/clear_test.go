package clear_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
	"github.com/jacekolszak/pixiq/tools/clear"
)

func TestNew(t *testing.T) {
	t.Run("should create tool", func(t *testing.T) {
		tool := clear.New()
		assert.NotNil(t, tool)
	})
}

func TestClear(t *testing.T) {
	t.Run("should clear selection", func(t *testing.T) {
		colors := []image.Color{
			image.RGBA(10, 20, 30, 40),
			image.RGBA(50, 60, 70, 80),
		}
		for _, color := range colors {
			tests := map[string]struct {
				width, height  int
				expectedColors [][]image.Color
			}{
				"top left corner": {
					width: 1, height: 1,
					expectedColors: [][]image.Color{
						{color, image.Transparent},
						{image.Transparent, image.Transparent},
					},
				},
				"top line": {
					width: 2, height: 1,
					expectedColors: [][]image.Color{
						{color, color},
						{image.Transparent, image.Transparent},
					},
				},
				"left line": {
					width: 1, height: 2,
					expectedColors: [][]image.Color{
						{color, image.Transparent},
						{color, image.Transparent},
					},
				},
			}
			for name, test := range tests {
				testName := fmt.Sprintf("%s %v", name, color)
				t.Run(testName, func(t *testing.T) {
					var (
						img       = image.New(2, 2, fake.NewAcceleratedImage(2, 2))
						selection = img.Selection(0, 0).WithSize(test.width, test.height)
						tool      = clear.New()
					)
					tool.SetColor(color)
					// when
					tool.Clear(selection)
					// then
					assertColors(t, img.WholeImageSelection(), test.expectedColors)
				})
			}
		}
	})
}

// TODO Make this function reusable
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
