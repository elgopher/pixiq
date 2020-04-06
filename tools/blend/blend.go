package blend

import (
	"github.com/jacekolszak/pixiq/image"
)

// ColorBlender blends source and target colors together. It is executed by Tool
// for each pixel in source and target selection.
type ColorBlender interface {
	BlendSourceToTargetColor(source, target image.Color) image.Color
}

// New creates a blending Tool with given ColorBlender implementation.
func New(colorBlender ColorBlender) *Tool {
	if colorBlender == nil {
		panic("nil colorBlender")
	}
	return &Tool{
		colorBlender: colorBlender,
	}
}

// NewSource creates a new blending tool which replaces target selection with source
// colors. It is like coping of source selection colors into target.
func NewSource() *Source {
	return &Source{}
}

// Source is a blending tool which replaces target selection with source
// colors. It is like coping of source selection colors into target.
type Source struct{}

// BlendSourceToTarget blends source into target. Only position of the target Selection
// is used and the source is not clamped by the target size.
func (s *Source) BlendSourceToTarget(source, target image.Selection) {
	target = target.WithSize(source.Width(), source.Height())
	var (
		sourceLines   = source.Lines()
		targetLines   = target.Lines()
		sourceXOffset = sourceLines.XOffset()
		sourceYOffset = sourceLines.YOffset()
		targetXOffset = targetLines.XOffset()
		targetYOffset = targetLines.YOffset()
		height        = source.Height()
	)
	if height > targetLines.Length()+targetYOffset {
		height = targetLines.Length() + targetYOffset
	}
	for y := targetYOffset; y < height; y++ {
		sourceOutOfBounds := y < sourceYOffset || y-sourceYOffset >= sourceLines.Length()
		if sourceOutOfBounds {
			targetLine := targetLines.LineForWrite(y - targetYOffset)
			for x := 0; x < len(targetLine); x++ {
				targetLine[x] = image.Transparent
			}
			continue
		}
		sourceLine := sourceLines.LineForRead(y - sourceYOffset)
		targetLine := targetLines.LineForWrite(y - targetYOffset)
		for x := 0; x < sourceXOffset-targetXOffset; x++ {
			targetLine[x] = image.Transparent
		}
		width := source.Width()
		if width > len(sourceLine) {
			width = len(sourceLine)
		}
		for x := targetXOffset + sourceXOffset; x < width; x++ {
			targetLine[x-targetXOffset] = sourceLine[x-sourceXOffset]
		}
		for x := width; x < source.Width(); x++ {
			targetLine[x] = image.Transparent
		}
	}
}

// NewSourceOver creates a new blending tool which blends together source and target
// taking into account alpha channel of both.
func NewSourceOver() *SourceOver {
	return NewSourceOverWithOpacity(100)
}

// NewSourceOver creates a new blending tool which blends together source and target
// taking into account alpha channel of both plus the additional opacity of the source.
func NewSourceOverWithOpacity(opacity int) *SourceOver {
	tool := &SourceOver{opacity: opacity}
	tool.Tool = New(tool)
	return tool
}

// SourceOver (aka Normal) is a blending tool which blends together source and target
// taking into account alpha channel of both plus the additional opacity of the source.
type SourceOver struct {
	*Tool
	opacity int
}

// BlendSourceToTarget blends source into target. Only position of the target Selection
// is used and the source is not clamped by the target size.
func (s *SourceOver) BlendSourceToTargetColor(source, target image.Color) image.Color {
	return source
}

// SetOpacity specifies the opacity of source selection.
func (s *SourceOver) SetOpacity(opacity int) {
	s.opacity = opacity
}

// Tool is a customizable blending tool which blends together two selections. It uses
// ColorBlender implementation for actual blending of two pixel colors.
type Tool struct {
	colorBlender ColorBlender
}

// BlendSourceToTarget blends source into target. Only position of the target Selection
// is used and the source is not clamped by the target size.
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
