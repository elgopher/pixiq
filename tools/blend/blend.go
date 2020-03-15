package blend

import "github.com/jacekolszak/pixiq/image"

// Override returns new instance of *blend.Tool which will override target with source
// color
func Override() *Tool {
	return &Tool{}
}

type Tool struct {
}

func (t *Tool) BlendSource(source, target image.Selection) {
	lines := source.Lines()
	for y := 0; y < lines.Length(); y++ {
		line := lines.LineForRead(y)
		for x := 0; x < len(line); x++ {
			target.SetColor(x, y, line[x])
		}
	}
}
