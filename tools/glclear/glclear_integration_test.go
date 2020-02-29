package glclear_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/tools/glclear"
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
func TestNew_Integration(t *testing.T) {
	t.Run("should create tool", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		command := openGL.Context().NewClearCommand()
		// when
		tool := glclear.New(command)
		assert.NotNil(t, tool)
	})
}

func TestTool_Clear(t *testing.T) {
	openGL, err := glfw.NewOpenGL(mainThreadLoop)
	require.NoError(t, err)
	defer openGL.Destroy()
	command := openGL.Context().NewClearCommand()
	t.Run("should clear selection", func(t *testing.T) {
		colors := []image.Color{
			image.RGBA(10, 20, 30, 40),
			image.RGBA(50, 60, 70, 80),
		}
		for _, color := range colors {
			tests := map[string]struct {
				width, height  int
				expectedColors [2][2]image.Color
			}{
				"top left corner": {
					width: 1, height: 1,
					expectedColors: [2][2]image.Color{
						{color, image.Transparent},
						{image.Transparent, image.Transparent},
					},
				},
				"top line": {
					width: 2, height: 1,
					expectedColors: [2][2]image.Color{
						{color, color},
						{image.Transparent, image.Transparent},
					},
				},
				"left line": {
					width: 1, height: 2,
					expectedColors: [2][2]image.Color{
						{color, image.Transparent},
						{color, image.Transparent},
					},
				},
			}
			for name, test := range tests {
				testName := fmt.Sprintf("%s %v", name, color)
				t.Run(testName, func(t *testing.T) {
					img := openGL.NewImage(2, 2)
					selection := img.Selection(0, 0).WithSize(test.width, test.height)
					tool := glclear.New(command)
					tool.SetColor(color)
					// when
					tool.Clear(selection)
					// then
					assertColors(t, img, test.expectedColors)
				})
			}
		}
	})
}

func assertColors(t *testing.T, img *image.Image, expectedColorLines [2][2]image.Color) {
	selection := img.WholeImageSelection()
	assert.Equal(t, expectedColorLines[0][0], selection.Color(0, 0))
	assert.Equal(t, expectedColorLines[0][1], selection.Color(1, 0))
	assert.Equal(t, expectedColorLines[1][0], selection.Color(0, 1))
	assert.Equal(t, expectedColorLines[1][1], selection.Color(1, 1))
}
