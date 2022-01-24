package blend_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pixiq/blend"
	"github.com/elgopher/pixiq/image"
	"github.com/elgopher/pixiq/image/fake"
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
		// colorTx2 is a result of blending transparent color with color2
		colorTx2 image.Color
	}{
		"Tool": {
			tool:     blend.New(multiplyColors{}),
			color1x2: image.RGBA(5, 12, 21, 32),
			color1x3: image.RGBA(9, 20, 33, 48),
			color3x4: image.RGBA(54, 70, 88, 108),
			colorTx2: image.Transparent,
		},
		"Source": {
			tool:     blend.NewSource(),
			color1x2: color1,
			color1x3: color1,
			color3x4: color3,
			colorTx2: image.Transparent,
		},
		"SourceOver": {
			tool:     blend.NewSourceOver(),
			color1x2: image.RGBA(6, 8, 10, 12),
			color1x3: image.RGBA(10, 12, 14, 16),
			color3x4: image.RGBA(15, 17, 19, 21),
			colorTx2: color2,
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
						source := image.New(fake.NewAcceleratedImage(test.width, test.height)).
							WholeImageSelection()
						target := image.New(fake.NewAcceleratedImage(1, 1)).
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
				colorsTx2 := [][]image.Color{
					{
						blender.colorTx2,
					},
				}
				tests := map[string]struct {
					source, target image.Selection
					expectedPixels [][]image.Color
				}{
					"x=-1": {
						source:         newImage(colors1).Selection(-1, 0).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"x=1": {
						source:         newImage(colors1).Selection(1, 0).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"y=-1": {
						source:         newImage(colors1).Selection(0, -1).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"y=1": {
						source:         newImage(colors1).Selection(0, 1).WithSize(1, 1),
						target:         newImage(colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"y=-2 and target y=-1": {
						source:         newImage(colors1).Selection(0, -2).WithSize(1, 2),
						target:         newImage(colors2).Selection(0, -1),
						expectedPixels: colorsTx2,
					},
					"y=-1 and target x=-1": {
						source:         newImage(colors1).Selection(0, -1).WithSize(2, 1),
						target:         newImage(colors2).Selection(-1, 0),
						expectedPixels: colorsTx2,
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
					source, target image.Selection
					expectedPixels [][]image.Color
				}{
					"1x1 images": {
						source: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
						source: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
						source: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
						source: newImage([][]image.Color{
							{
								color1, color3,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
						source: newImage([][]image.Color{
							{
								color1, color3,
							},
						}).Selection(0, 0).WithSize(1, 1),
						target: newImage([][]image.Color{
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
						source: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).Selection(0, 0).WithSize(1, 1),
						target: newImage([][]image.Color{
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
						source: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
					"target out of boundaries x=-1": {
						source: newImage([][]image.Color{
							{
								color1, color3,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
					"target out of boundaries x=1": {
						source: newImage([][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
							{
								color4,
							},
						}).Selection(1, 0),
						expectedPixels: [][]image.Color{
							{
								color4,
							},
						},
					},
					"target out of boundaries y=-1": {
						source: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
					"target out of boundaries y=1": {
						source: newImage([][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
							{
								color4,
							},
						}).Selection(0, 1),
						expectedPixels: [][]image.Color{
							{
								color4,
							},
						},
					},
					"source wider than target": {
						source: newImage([][]image.Color{
							{
								color1, color2,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
					"source higher than target": {
						source: newImage([][]image.Color{
							{
								color1,
							},
							{
								color2,
							},
						}).WholeImageSelection(),
						target: newImage([][]image.Color{
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
						blender.tool.BlendSourceToTarget(test.source, test.target)
						// then
						assertColors(t, test.target.Image(), test.expectedPixels)
					})
				}
			})

		})
	}
}

func TestSourceOver_BlendSourceToTarget(t *testing.T) {
	t.Run("should blend color", func(t *testing.T) {
		tests := map[string]struct {
			source   image.Color
			target   image.Color
			expected image.Color
		}{
			"all transparent": {
				source:   image.Transparent,
				target:   image.Transparent,
				expected: image.Transparent,
			},
			"transparent source, fully opaque target": {
				source:   image.Transparent,
				target:   image.NRGBA(4, 5, 6, 255),
				expected: image.NRGBA(4, 5, 6, 255),
			},
			"fully opaque source": {
				source:   image.NRGBA(1, 2, 3, 255),
				target:   image.NRGBA(4, 5, 6, 100),
				expected: image.NRGBA(1, 2, 3, 255),
			},
			"semi-transparent white source, black opaque target": {
				source:   image.NRGBA(255, 255, 255, 127),
				target:   image.NRGBA(0, 0, 0, 255),
				expected: image.NRGBA(127, 127, 127, 255),
			},
			"semi-transparent violet source, black opaque target": {
				source:   image.NRGBA(118, 66, 138, 127),
				target:   image.NRGBA(0, 0, 0, 255),
				expected: image.NRGBA(58, 32, 68, 255),
			},
			"semi-transparent violet source, blue opaque target": {
				source:   image.NRGBA(118, 66, 138, 127),
				target:   image.NRGBA(99, 155, 255, 255),
				expected: image.NRGBA(108, 111, 197, 255),
			},
			"semi-transparent white source, semi-transparent black target": {
				source:   image.NRGBA(255, 255, 255, 127),
				target:   image.NRGBA(0, 0, 0, 127),
				expected: image.NRGBA(169, 169, 169, 191),
			},
			"semi-transparent violet source, semi-transparent blue target": {
				source:   image.NRGBA(118, 66, 138, 127),
				target:   image.NRGBA(99, 155, 255, 127),
				expected: image.NRGBA(111, 96, 178, 191),
			},
			"alpha 0": {
				source:   image.NRGBA(1, 3, 5, 0),
				target:   image.NRGBA(2, 4, 6, 0),
				expected: image.Transparent,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pre := newImage([][]image.Color{
					{
						image.Transparent,
					},
				}).WholeImageSelection()

				source := newImage([][]image.Color{
					{
						test.source,
					},
				}).WholeImageSelection()
				target := newImage([][]image.Color{
					{
						test.target,
					},
				}).WholeImageSelection()
				// when
				blend.NewSourceOver().BlendSourceToTarget(pre, source)
				blend.NewSourceOver().BlendSourceToTarget(source, target)
				// then
				result := target.Color(0, 0)
				const delta = 1
				assert.InDelta(t, test.expected.R(), result.R(), delta, "Red")
				assert.InDelta(t, test.expected.G(), result.G(), delta, "Green")
				assert.InDelta(t, test.expected.B(), result.B(), delta, "Blue")
				assert.InDelta(t, test.expected.A(), result.A(), delta, "Alpha")
			})
		}
	})
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
	img := image.New(fake.NewAcceleratedImage(width, height))
	selection := img.WholeImageSelection()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			selection.SetColor(x, y, pixels[y][x])
		}
	}
	return img
}
