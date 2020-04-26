package glfw_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/tools/glblend"
)

var mainThreadLoop *glfw.MainThreadLoop

func TestMain(m *testing.M) {
	var exit int
	glfw.StartMainThreadLoop(func(main *glfw.MainThreadLoop) {
		mainThreadLoop = main
		exit = m.Run()
	})
	os.Exit(exit)
}

func TestNewSource(t *testing.T) {
	t.Run("should return source blender", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		context := openGL.Context()
		// when
		source, err := glblend.NewSource(context)
		// then
		assert.NotNil(t, source)
		assert.NoError(t, err)
	})
}

type blender interface {
	BlendSourceToTarget(source, target image.Selection)
}

func TestBlendSourceToTarget(t *testing.T) {
	var (
		color1 = image.RGBA(1, 2, 3, 4)
		color2 = image.RGBA(5, 6, 7, 8)
		color3 = image.RGBA(9, 10, 11, 12)
		color4 = image.RGBA(6, 7, 8, 9)
	)
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()
	context := openGL.Context()

	blenders := map[string]struct {
		tool                         func() blender
		color1x2, color1x3, color3x4 image.Color
		// colorTx2 is a result of blending transparent color with color2
		colorTx2 image.Color
	}{
		"Source": {
			tool: func() blender {
				b, _ := glblend.NewSource(context)
				return b
			},
			color1x2: color1,
			color1x3: color1,
			color3x4: color3,
			colorTx2: image.Transparent,
		},
		"SourceOver": {
			tool: func() blender {
				b, _ := glblend.NewSourceOver(context)
				return b
			},
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
						source := openGL.NewImage(test.width, test.height).
							WholeImageSelection()
						target := openGL.NewImage(1, 1).
							WholeImageSelection()
						originalColor := image.RGBA(1, 2, 3, 4)
						target.SetColor(0, 0, originalColor)
						// when
						blender.tool().BlendSourceToTarget(source, target)
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
						source:         newImage(openGL, colors1).Selection(-1, 0).WithSize(1, 1),
						target:         newImage(openGL, colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"x=1": {
						source:         newImage(openGL, colors1).Selection(1, 0).WithSize(1, 1),
						target:         newImage(openGL, colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"y=-1": {
						source:         newImage(openGL, colors1).Selection(0, -1).WithSize(1, 1),
						target:         newImage(openGL, colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"y=1": {
						source:         newImage(openGL, colors1).Selection(0, 1).WithSize(1, 1),
						target:         newImage(openGL, colors2).Selection(0, 0),
						expectedPixels: colorsTx2,
					},
					"y=-2 and target y=-1": {
						source:         newImage(openGL, colors1).Selection(0, -2).WithSize(1, 2),
						target:         newImage(openGL, colors2).Selection(0, -1),
						expectedPixels: colorsTx2,
					},
					"y=-1 and target x=-1": {
						source:         newImage(openGL, colors1).Selection(0, -1).WithSize(2, 1),
						target:         newImage(openGL, colors2).Selection(-1, 0),
						expectedPixels: colorsTx2,
					},
				}
				for name, test := range tests {
					t.Run(name, func(t *testing.T) {
						// when
						blender.tool().BlendSourceToTarget(test.source, test.target)
						// then
						assertColors(t, test.target.Image(), test.expectedPixels)
					})
				}
			})

			t.Run("source is not modified", func(t *testing.T) {
				sourceOriginalColors := [][]image.Color{
					{image.RGBA(1, 2, 3, 4)},
				}
				source := newImage(openGL, sourceOriginalColors).WholeImageSelection()
				target := newImage(openGL, [][]image.Color{
					{image.RGBA(5, 6, 7, 8)},
				}).WholeImageSelection()
				// when
				blender.tool().BlendSourceToTarget(source, target)
				// then
				assertColors(t, source.Image(), sourceOriginalColors)
			})

			t.Run("should blend selections", func(t *testing.T) {
				tests := map[string]struct {
					source, target image.Selection
					expectedPixels [][]image.Color
				}{
					"1x1 images": {
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1, color3,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1, color3,
							},
						}).Selection(0, 0).WithSize(1, 1),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).Selection(0, 0).WithSize(1, 1),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1, color3,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1, color2,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						source: newImage(openGL, [][]image.Color{
							{
								color1,
							},
							{
								color2,
							},
							{
								color3,
							},
						}).WholeImageSelection(),
						target: newImage(openGL, [][]image.Color{
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
						blender.tool().BlendSourceToTarget(test.source, test.target)
						// then
						assertColors(t, test.target.Image(), test.expectedPixels)
					})
				}
			})

		})
	}
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

func newImage(gl *glfw.OpenGL, pixels [][]image.Color) *image.Image {
	width := len(pixels[0])
	height := len(pixels)
	img := gl.NewImage(width, height)
	selection := img.WholeImageSelection()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			selection.SetColor(x, y, pixels[y][x])
		}
	}
	return img
}
