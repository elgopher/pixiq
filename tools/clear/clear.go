package clear

import (
	"github.com/jacekolszak/pixiq/image"
)

// New returns new instance of *clear.Tool
func New() *Tool {
	return &Tool{}
}

// Tool is a clearing tool. It clears the image.Selection with specific color
// which basically means setting the color for each pixel in the selection.
//
// Tool uses CPU.
type Tool struct {
	color image.Color
}

// SetColor sets color which will be used by Clear method
func (t *Tool) SetColor(color image.Color) {
	t.color = color
}

// Clear clears selection with previously set color
func (t *Tool) Clear(selection image.Selection) {
	for y := 0; y < selection.Height(); y++ {
		for x := 0; x < selection.Width(); x++ {
			selection.SetColor(x, y, t.color)
		}
	}
}
