// Package clear provides CPU tools for clearing selections
// using specific color
package clear

import (
	"github.com/elgopher/pixiq/image"
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
	lines := selection.Lines()
	for y := 0; y < lines.Length(); y++ {
		line := lines.LineForWrite(y)
		for x := 0; x < len(line); x++ {
			line[x] = t.color
		}
	}
}
