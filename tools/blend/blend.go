package blend

import (
	"github.com/jacekolszak/pixiq/image"
)

type ColorBlender interface {
	BlendSourceToTargetColor(source, target image.Color) image.Color
}

func New(colorBlender ColorBlender) *Tool {
	if colorBlender == nil {
		panic("nil colorBlender")
	}
	return &Tool{
		colorBlender: colorBlender,
	}
}

func NewSource() *Source {
	return &Source{}
}

type Source struct{}

func (s *Source) BlendSourceToTarget(source, target image.Selection) {
	for y := 0; y < source.Height(); y++ {
		for x := 0; x < source.Width(); x++ {
			target.SetColor(x, y, source.Color(x, y))
		}
	}
	//if source.Width() == 0 || source.Height() == 0 {
	//	return
	//}
	//sourceLines := source.Lines()
	//if target.Width() == 0 && target.Height() == 0 {
	//	target = target.WithSize(source.Width(), source.Height())
	//}
	//targetLines := target.Lines()
	//for y := 0; y < targetLines.Length(); y++ {
	//	targetLine := targetLines.LineForWrite(y)
	//	sourceLine := sourceLines.LineForRead(y)
	//	for x := 0; x < len(targetLine); x++ {
	//		targetLine[x] = sourceLine[x]
	//	}
	//}
}

func NewSourceOver() *SourceOver {
	return NewSourceOverWithOpacity(100)
}

func NewSourceOverWithOpacity(opacity int) *SourceOver {
	tool := &SourceOver{opacity: opacity}
	tool.Tool = New(tool)
	return tool
}

// Aka Normal
type SourceOver struct {
	*Tool
	opacity int
}

func (s *SourceOver) BlendSourceToTargetColor(source, target image.Color) image.Color {
	return source
}

func (s *SourceOver) SetOpacity(opacity int) {
	s.opacity = opacity
}

type Tool struct {
	colorBlender ColorBlender
}

// BlendSourceToTarget blends source into target. The source is not clamped by target
// size.
func (t *Tool) BlendSourceToTarget(source, target image.Selection) {
	for y := 0; y < source.Height(); y++ {
		for x := 0; x < source.Width(); x++ {
			sourceColor := source.Color(x, y)
			targetColor := target.Color(x, y)
			color := t.colorBlender.BlendSourceToTargetColor(sourceColor, targetColor)
			target.SetColor(x, y, color)
		}
	}
}
