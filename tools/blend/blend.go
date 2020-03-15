package blend

import "github.com/jacekolszak/pixiq/image"

// Override returns new instance of *blend.OverrideTool which will override target with source
// color
func Override() *OverrideTool {
	return &OverrideTool{}
}

type OverrideTool struct {
}

// BlendSource blends source into target.
//
// If target has 0x0 size then whole source is blended, otherwise source is clamped.
func (t *OverrideTool) BlendSource(source, target image.Selection) {
	lines := source.Lines()
	for y := 0; y < lines.Length(); y++ {
		line := lines.LineForRead(y)
		for x := 0; x < len(line); x++ {
			target.SetColor(x, y, line[x])
		}
	}
}
