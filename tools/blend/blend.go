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

// BlendSourceToTarget blends source into target selection.
// Only position of the target Selection is used and the source is not clamped by
// the target size.
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
		width := len(sourceLine)
		if len(targetLine) < width-targetXOffset {
			width = len(targetLine)
		}
		for x := targetXOffset + sourceXOffset; x < width; x++ {
			targetLine[x-targetXOffset] = sourceLine[x-sourceXOffset]
		}
		for x := len(sourceLine); x < source.Width(); x++ {
			targetLine[x] = image.Transparent
		}
	}
}

// NewSourceOver creates a new blending tool which blends together source and target
// taking into account alpha channel of both. Source-over means that source will be
// painted on top of the target.
func NewSourceOver() *SourceOver {
	return &SourceOver{}
}

// SourceOver (aka Normal) is a blending tool which blends together source and target
// taking into account alpha channel of both. Source-over means that source will be
// painted on top of the target.
type SourceOver struct{}

// BlendSourceToTarget blends source into target selection. Results will be stored
// in the image pointed by target selection
// Only position of the target Selection is used and the source is not clamped by
// the target size.
func (s *SourceOver) BlendSourceToTarget(source, target image.Selection) {
	target = target.WithSize(source.Width(), source.Height())
	var (
		sourceLines   = source.Lines()
		targetLines   = target.Lines()
		sourceXOffset = sourceLines.XOffset()
		sourceYOffset = sourceLines.YOffset()
		targetXOffset = targetLines.XOffset()
		targetYOffset = targetLines.YOffset()
		height        = sourceLines.Length()
	)
	if height > targetLines.Length()+targetYOffset {
		height = targetLines.Length() + targetYOffset
	}
	startY := targetYOffset
	if startY < sourceYOffset {
		startY = sourceYOffset
	}
	for y := startY; y < height; y++ {
		sourceLine := sourceLines.LineForRead(y - sourceYOffset)
		targetLine := targetLines.LineForWrite(y - targetYOffset)
		width := len(sourceLine)
		if len(targetLine) < width-targetXOffset {
			width = len(targetLine)
		}
		for x := targetXOffset + sourceXOffset; x < width; x++ {
			// blend source with target color (following block of code is inlined to improve performance)
			source := sourceLine[x-sourceXOffset]
			target := targetLine[x-targetXOffset]
			srcR, srcG, srcB, srcA := source.RGBAi()
			dstR, dstG, dstB, dstA := target.RGBAi()
			dstFactor := 255 - srcA
			outR := srcR + mul(dstR, dstFactor)
			outG := srcG + mul(dstG, dstFactor)
			outB := srcB + mul(dstB, dstFactor)
			outA := srcA + mul(dstA, dstFactor)
			targetLine[x-targetXOffset] = image.RGBAi(outR, outG, outB, outA)
		}
	}
}

// mul is an optimized version of round(a * b / 255)
func mul(a, b int) int {
	t := a*b + 0x80
	return ((t >> 8) + t) >> 8
}

// Tool is a customizable blending tool which blends together two selections. It uses
// ColorBlender implementation for actual blending of two pixel colors.
type Tool struct {
	colorBlender ColorBlender
}

// BlendSourceToTarget blends source into target selection. Results will be stored
// in the image pointed by target selection.
// Only position of the target Selection is used and the source is not clamped by
// the target size.
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
