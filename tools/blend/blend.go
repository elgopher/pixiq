package blend

import "github.com/jacekolszak/pixiq/image"

// TODO For optimization purposes there can be a dedicated tool which will override
// colors without using a blendFunc
var Source = func(source, target image.Color) image.Color {
	return source
}

// Aka Normal
var SourceOver = func(source, target image.Color) image.Color {
	return source
}

func SourceOverWithOpacity(opacity byte) func(source, target image.Color) image.Color {
	return func(source, target image.Color) image.Color {
		return source
	}
}

func New(blendFunc func(source, target image.Color) image.Color) *Tool {
	return &Tool{
		blendFunc: blendFunc,
	}
}

type Tool struct {
	blendFunc func(source, target image.Color) image.Color
}

// BlendSource blends source into target.
//
// If target has 0x0 size then whole source is blended, otherwise source is clamped.
func (t *Tool) BlendSource(source, target image.Selection) {
	lines := source.Lines()
	for y := 0; y < lines.Length(); y++ {
		line := lines.LineForRead(y)
		for x := 0; x < len(line); x++ {
			sourceColor := line[x]
			targetColor := target.Color(x, y)
			target.SetColor(x, y, t.blendFunc(sourceColor, targetColor))
		}
	}
}
