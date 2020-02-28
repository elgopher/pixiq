package clear

import (
	"github.com/jacekolszak/pixiq/image"
)

func New() *Tool {
	return &Tool{}
}

type Tool struct {
	color image.Color
}

func (t *Tool) SetColor(color image.Color) {
	t.color = color
}

func (t *Tool) Clear(selection image.Selection) {
	for y := 0; y < selection.Height(); y++ {
		for x := 0; x < selection.Width(); x++ {
			selection.SetColor(x, y, t.color)
		}
	}
}
