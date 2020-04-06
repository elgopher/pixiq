package blend_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/image/fake"
	"github.com/jacekolszak/pixiq/tools/blend"
)

func TestNew(t *testing.T) {
	t.Run("should panic when colorBlender is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			blend.New(nil)
		})
	})
	t.Run("should create tool", func(t *testing.T) {
		tool := blend.New(multiplyColors{})
		assert.NotNil(t, tool)
	})
}

func TestBlendSourceToTarget(t *testing.T) {
	var (
		color1 = image.RGBA(1, 2, 3, 4)
		color2 = image.RGBA(5, 6, 7, 8)
		color3 = image.RGBA(9, 10, 11, 12)
		color4 = image.RGBA(6, 7, 8, 9)
	)
	blenders := map[string]struct {
		tool interface {
			BlendSourceToTarget(source, target image.Selection)
		}
		color1x2, color1x3, color3x4 image.Color
	}{
		"Tool": {
			tool:     blend.New(multiplyColors{}),
			color1x2: image.RGBA(5, 12, 21, 32),
			color1x3: image.RGBA(9, 20, 33, 48),
			color3x4: image.RGBA(54, 70, 88, 108),
		},
		"Source": {
			tool:     blend.NewSource(),
			color1x2: color1,
			color1x3: color1,
			color3x4: color3,
		},
		// TODO SourceOver should behave differently than Source!
		"SourceOver": {
			tool:     blend.NewSourceOver(),
			color1x2: color1,
			color1x3: color1,
			color3x4: color3,
		},
	}
	for name, blender := range blenders {

		t.Run(name, func(t *testing.T) {

			t.Run("should skip blending when source image has 0 size", func(t *testing.T) {
				tests := map[string]struct {
					width, height int
				}{
					"0 height": {
						width: 1,
					},
					"0 width": {
						height: 1,
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						// when
						source := image.New(test.width, test.height, fake.NewAcceleratedImage(test.width, test.height)).
							WholeImageSelection()
						target := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).
							WholeImageSelection()
						originalColor := image.RGBA(1, 2, 3, 4)
						target.SetColor(0, 0, originalColor)
						// when
						blender.tool.BlendSourceToTarget(source, target)
						// then
						assert.Equal(t, originalColor, target.Color(0, 0))
					})
				}
			})

			t.Run("source selection out of boundaries", func(t *testing.T) {
				colors1 := [][]image.Color{
					{
						color1,
					},
				}
				colors2 := [][]image.Color{
					{
						color2,
					},
				}
				transparent := [][]image.Color{
					{
						image.Transparent,
					},
				}
				tests := map[string]struct {
					source, target image.Selection
					expectedPixels [][]image.Color
				}{
					"x=-1": {
						source:         newImage(colors1).Selection(-1, 0).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: transparent,
					},
					"x=1": {
						source:         newImage(colors1).Selection(1, 0).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: transparent,
					},
					"y=-1": {
						source:         newImage(colors1).Selection(0, -1).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: transparent,
					},
					"y=1": {
						source:         newImage(colors1).Selection(0, 1).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: transparent,
					},
					"y=-2 and target y=-1": {
						source:         newImage(colors1).Selection(0, -2).WithSize(1, 2),
						target:         newImage(colors2).Selection(0, -1),
						expectedPixels: transparent,
					},
					"y=-1 and target x=-1": {
						source:         newImage(colors1).Selection(0, -1).WithSize(2, 1),
						target:         newImage(colors2).Selection(-1, 0),
						expectedPixels: transparent,
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						// when
						blender.tool.BlendSourceToTarget(test.source, test.target)
						// then
						assertColors(t, test.target.Image(), test.expectedPixels)
					})
				}
			})

			t.Run("source is not modified", func(t *testing.T) {
				sourceOriginalColors := [][]image.Color{
					{image.RGBA(1, 2, 3, 4)},
				}
				source := newImage(sourceOriginalColors).WholeImageSelection()
				target := newImage([][]image.Color{
					{image.RGBA(5, 6, 7, 8)},
				}).WholeImageSelection()
				// when
				blender.tool.BlendSourceToTarget(source, target)
				// then
				assertColors(t, source.Image(), sourceOriginalColors)
			})

			t.Run("should blend selections", func(t *testing.T) {
				tests := map[string]struct {
					sourceSelection, targetSelection image.Selection
					expectedPixels                   [][]image.Color
				}{
					"1x1 images": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color2,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2,
							},
						},
					},
					"target bigger than source 1": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color2, color3,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2, color3,
							},
						},
					},
					"target bigger than source 2": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color2,
							},
							{
								color3,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2,
							},
							{
								color3,
							},
						},
					},
					"2x1 images": {
						sourceSelection: newImage([][]image.Color{
							{
								color1, color3,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color2, color4,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2, blender.color3x4,
							},
						},
					},
					"source clamped x": {
						sourceSelection: newImage([][]image.Color{
							{
								color1, color3,
							},
						}).Selection(0, 0).WithSize(1, 1),
						targetSelection: newImage([][]image.Color{
							{
								color2, color4,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2, color4,
							},
						},
					},
					"source clamped y": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).Selection(0, 0).WithSize(1, 1),
						targetSelection: newImage([][]image.Color{
							{
								color2,
							},
							{
								color4,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2,
							},
							{
								color4,
							},
						},
					},
					"1x2 images": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color2,
							},
							{
								color4,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x2,
							},
							{
								blender.color3x4,
							},
						},
					},
					"target out boundaries x=-1": {
						sourceSelection: newImage([][]image.Color{
							{
								color1, color3,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color4,
							},
						}).Selection(-1, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color3x4,
							},
						},
					},
					"target out boundaries y=-1": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color4,
							},
						}).Selection(0, -1),
						expectedPixels: [][]image.Color{
							{
								blender.color3x4,
							},
						},
					},
					"source higher than target": {
						sourceSelection: newImage([][]image.Color{
							{
								color1,
							},
							{
								color2,
							},
						}).WholeImageSelection(),
						targetSelection: newImage([][]image.Color{
							{
								color3,
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								blender.color1x3,
							},
						},
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						// when
						blender.tool.BlendSourceToTarget(test.sourceSelection, test.targetSelection)
						// then
						assertColors(t, test.targetSelection.Image(), test.expectedPixels)
					})
				}
			})

		})
	}
}

type multiplyColors struct{}

func (c multiplyColors) BlendSourceToTargetColor(source, target image.Color) image.Color {
	return image.RGBA(
		source.R()*target.R(),
		source.G()*target.G(),
		source.B()*target.B(),
		source.A()*target.A())
}

func assertColors(t *testing.T, img *image.Image, expectedColorLines [][]image.Color) {
	selection := img.WholeImageSelection()
	for y := 0; y < selection.Height(); y++ {
		expectedColorLine := expectedColorLines[y]
		for x := 0; x < selection.Width(); x++ {
			color := selection.Color(x, y)
			assert.Equal(t, expectedColorLine[x], color, "position (%d,%d)", x, y)
		}
	}
}

func newImage(pixels [][]image.Color) *image.Image {
	width := len(pixels[0])
	height := len(pixels)
	img := image.New(width, height, fake.NewAcceleratedImage(width, height))
	selection := img.WholeImageSelection()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			selection.SetColor(x, y, pixels[y][x])
		}
	}
	return img
}
