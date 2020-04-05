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
	blenders := map[string]interface {
		BlendSourceToTarget(source, target image.Selection)
	}{
		"Tool":       blend.New(multiplyColors{}),
		"Source":     blend.NewSource(),
		"SourceOver": blend.NewSourceOver(),
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
						source := image.New(test.width, test.height, fake.NewAcceleratedImage(test.width, test.height)).WholeImageSelection()
						target := image.New(1, 1, fake.NewAcceleratedImage(1, 1)).WholeImageSelection()
						originalColor := image.RGBA(1, 2, 3, 4)
						target.SetColor(0, 0, originalColor)
						// when
						blender.BlendSourceToTarget(source, target)
						// then
						assert.Equal(t, originalColor, target.Color(0, 0))
					})
				}
			})

			t.Run("source selection of boundaries", func(t *testing.T) {
				tests := map[string]struct {
					sourceSelection, targetSelection image.Selection
					expectedPixels                   [][]image.Color
				}{
					"x=-1": {
						sourceSelection: newImage([][]image.Color{
							{
								image.RGBA(5, 6, 7, 8),
							},
						}).Selection(-1, 0).WithSize(1, 1),
						targetSelection: newImage([][]image.Color{
							{
								image.RGBA(1, 2, 3, 4),
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								image.RGBA(0, 0, 0, 0),
							},
						},
					},
					"x=1": {
						sourceSelection: newImage([][]image.Color{
							{
								image.RGBA(5, 6, 7, 8),
							},
						}).Selection(1, 0).WithSize(1, 1),
						targetSelection: newImage([][]image.Color{
							{
								image.RGBA(1, 2, 3, 4),
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								image.RGBA(0, 0, 0, 0),
							},
						},
					},
					"y=-1": {
						sourceSelection: newImage([][]image.Color{
							{
								image.RGBA(5, 6, 7, 8),
							},
						}).Selection(0, -1).WithSize(1, 1),
						targetSelection: newImage([][]image.Color{
							{
								image.RGBA(1, 2, 3, 4),
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								image.RGBA(0, 0, 0, 0),
							},
						},
					},
					"y=1": {
						sourceSelection: newImage([][]image.Color{
							{
								image.RGBA(5, 6, 7, 8),
							},
						}).Selection(0, 1).WithSize(1, 1),
						targetSelection: newImage([][]image.Color{
							{
								image.RGBA(1, 2, 3, 4),
							},
						}).Selection(0, 0),
						expectedPixels: [][]image.Color{
							{
								image.RGBA(0, 0, 0, 0),
							},
						},
					},
					"y=-2 and target y=-1": {
						sourceSelection: newImage([][]image.Color{
							{
								image.RGBA(1, 2, 3, 4),
							},
						}).Selection(0, -2).WithSize(1, 2),
						targetSelection: newImage([][]image.Color{
							{
								image.RGBA(5, 6, 7, 8),
							},
						}).Selection(0, -1),
						expectedPixels: [][]image.Color{
							{
								image.RGBA(0, 0, 0, 0),
							},
						},
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						// when
						blender.BlendSourceToTarget(test.sourceSelection, test.targetSelection)
						// then
						assertColors(t, test.targetSelection.Image(), test.expectedPixels)
					})
				}
			})

		})
	}
	// TODO Test if source is left unmodified
}

func TestTool_BlendSourceToTarget(t *testing.T) {
	t.Run("should blend selections", func(t *testing.T) {
		tests := map[string]struct {
			sourceSelection, targetSelection image.Selection
			expectedPixels                   [][]image.Color
		}{
			"1x1 images": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32),
					},
				},
			},
			"target bigger than source 1": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32), image.RGBA(9, 10, 11, 12),
					},
				},
			},
			"target bigger than source 2": {
				sourceSelection: newImage([][]image.Color{
					{image.RGBA(1, 2, 3, 4)},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{image.RGBA(5, 6, 7, 8)},
					{image.RGBA(9, 10, 11, 12)},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{image.RGBA(5, 12, 21, 32)},
					{image.RGBA(9, 10, 11, 12)},
				},
			},
			"2x1 images": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32), image.RGBA(54, 70, 88, 108),
					},
				},
			},
			"source clamped x": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0).WithSize(1, 1),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32), image.RGBA(6, 7, 8, 9),
					},
				},
			},
			"source clamped y": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
					{
						image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0).WithSize(1, 1),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32),
					},
					{
						image.RGBA(6, 7, 8, 9),
					},
				},
			},
			"1x2 images": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
					{
						image.RGBA(9, 10, 11, 12),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32),
					},
					{
						image.RGBA(54, 70, 88, 108),
					},
				},
			},
			"target out boundaries x=-1": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(5, 6, 7, 8),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).Selection(-1, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32),
					},
				},
			},
			"target out boundaries y=-1": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(5, 6, 7, 8),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).Selection(0, -1),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 12, 21, 32),
					},
				},
			},
			"source higher than target": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(5, 6, 7, 8),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(9, 20, 33, 48),
					},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				tool := blend.New(multiplyColors{})
				// when
				tool.BlendSourceToTarget(test.sourceSelection, test.targetSelection)
				// then
				assertColors(t, test.targetSelection.Image(), test.expectedPixels)
			})
		}
	})
}

func TestSource_BlendSourceToTarget(t *testing.T) {
	t.Run("should blend selections", func(t *testing.T) {
		tests := map[string]struct {
			sourceSelection, targetSelection image.Selection
			expectedPixels                   [][]image.Color
		}{
			"1x1 images": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				},
			},
			"target bigger than source 1": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(9, 10, 11, 12),
					},
				},
			},
			"target bigger than source 2": {
				sourceSelection: newImage([][]image.Color{
					{image.RGBA(1, 2, 3, 4)},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{image.RGBA(5, 6, 7, 8)},
					{image.RGBA(9, 10, 11, 12)},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{image.RGBA(1, 2, 3, 4)},
					{image.RGBA(9, 10, 11, 12)},
				},
			},
			"2x1 images": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				},
			},
			"source clamped x": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0).WithSize(1, 1),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 6, 7, 8), image.RGBA(6, 7, 8, 9),
					},
				},
			},
			"source clamped y": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
					{
						image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0).WithSize(1, 1),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
					{
						image.RGBA(6, 7, 8, 9),
					},
				},
			},
			"1x2 images": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
					{
						image.RGBA(9, 10, 11, 12),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(6, 7, 8, 9),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
					{
						image.RGBA(9, 10, 11, 12),
					},
				},
			},
			"target out boundaries x=-1": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4), image.RGBA(5, 6, 7, 8),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).Selection(-1, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
				},
			},
			"target out boundaries y=-1": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(5, 6, 7, 8),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				}).Selection(0, -1),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(5, 6, 7, 8),
					},
				},
			},
			"source higher than target": {
				sourceSelection: newImage([][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
					{
						image.RGBA(5, 6, 7, 8),
					},
				}).WholeImageSelection(),
				targetSelection: newImage([][]image.Color{
					{
						image.RGBA(9, 10, 11, 12),
					},
				}).Selection(0, 0),
				expectedPixels: [][]image.Color{
					{
						image.RGBA(1, 2, 3, 4),
					},
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				tool := blend.NewSource()
				// when
				tool.BlendSourceToTarget(test.sourceSelection, test.targetSelection)
				// then
				assertColors(t, test.targetSelection.Image(), test.expectedPixels)
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
	img := image.New(width, height, fake.NewAcceleratedImage(width, height))
	selection := img.WholeImageSelection()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			selection.SetColor(x, y, pixels[y][x])
		}
	}
	return img
}
