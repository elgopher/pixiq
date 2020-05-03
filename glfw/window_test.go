package glfw_test

import (
	"fmt"
	"testing"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gl2 "github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
)

func TestWindow_DrawIntoBackBuffer(t *testing.T) {
	t.Run("should draw screen image", func(t *testing.T) {
		color1 := image.RGBA(10, 20, 30, 40)
		color2 := image.RGBA(50, 60, 70, 80)
		color3 := image.RGBA(90, 100, 110, 120)
		color4 := image.RGBA(130, 140, 150, 160)

		t.Run("1x1", func(t *testing.T) {
			openGL, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
			window, err := openGL.OpenWindow(1, 1, glfw.NoDecorationHint())
			require.NoError(t, err)
			defer window.Close()
			window.Image().WholeImageSelection().SetColor(0, 0, color1)
			// when
			window.DrawIntoBackBuffer()
			// then
			expected := []image.Color{color1}
			assert.Equal(t, expected, framebufferPixels(window.ContextAPI(), 0, 0, 1, 1))
		})
		t.Run("1x2", func(t *testing.T) {
			openGL, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
			window, err := openGL.OpenWindow(1, 2, glfw.NoDecorationHint())
			require.NoError(t, err)
			defer window.Close()
			img := window.Image()
			img.WholeImageSelection().SetColor(0, 0, color1)
			img.WholeImageSelection().SetColor(0, 1, color2)
			// when
			window.DrawIntoBackBuffer()
			// then
			expected := []image.Color{color2, color1}
			assert.Equal(t, expected, framebufferPixels(window.ContextAPI(), 0, 0, 1, 2))
		})
		t.Run("2x1", func(t *testing.T) {
			openGL, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
			window, err := openGL.OpenWindow(2, 1, glfw.NoDecorationHint())
			require.NoError(t, err)
			defer window.Close()
			img := window.Image()
			img.WholeImageSelection().SetColor(0, 0, color1)
			img.WholeImageSelection().SetColor(1, 0, color2)
			// when
			window.DrawIntoBackBuffer()
			// then
			expected := []image.Color{color1, color2}
			assert.Equal(t, expected, framebufferPixels(window.ContextAPI(), 0, 0, 2, 1))
		})
		t.Run("2x2", func(t *testing.T) {
			openGL, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
			window, err := openGL.OpenWindow(2, 2, glfw.NoDecorationHint())
			require.NoError(t, err)
			defer window.Close()
			img := window.Image()
			selection := img.WholeImageSelection()
			selection.SetColor(0, 0, color1)
			selection.SetColor(1, 0, color2)
			selection.SetColor(0, 1, color3)
			selection.SetColor(1, 1, color4)
			// when
			window.DrawIntoBackBuffer()
			// then
			expected := []image.Color{color3, color4, color1, color2}
			assert.Equal(t, expected, framebufferPixels(window.ContextAPI(), 0, 0, 2, 2))
		})

		t.Run("zoom < 1 should not change the framebuffer size", func(t *testing.T) {
			for zoom := -1; zoom < 1; zoom++ {
				name := fmt.Sprintf("zoom=%d", zoom)
				t.Run(name, func(t *testing.T) {
					openGL, err := glfw.NewOpenGL(mainThreadLoop)
					require.NoError(t, err)
					defer openGL.Destroy()
					window, err := openGL.OpenWindow(1, 1, glfw.NoDecorationHint(), glfw.Zoom(zoom))
					require.NoError(t, err)
					defer window.Close()
					img := window.Image()
					img.WholeImageSelection().SetColor(0, 0, color1)
					// when
					window.DrawIntoBackBuffer()
					// then
					expected := []image.Color{color1}
					assert.Equal(t, expected, framebufferPixels(window.ContextAPI(), 0, 0, 1, 1))
				})
			}
		})

		t.Run("zoom > 1 should make framebuffer zoom times bigger", func(t *testing.T) {
			for zoom := 2; zoom < 4; zoom++ {
				name := fmt.Sprintf("zoom=%d", zoom)
				t.Run(name, func(t *testing.T) {
					openGL, err := glfw.NewOpenGL(mainThreadLoop)
					require.NoError(t, err)
					defer openGL.Destroy()
					window, err := openGL.OpenWindow(1, 1, glfw.NoDecorationHint(), glfw.Zoom(zoom))
					require.NoError(t, err)
					defer window.Close()
					img := window.Image()
					img.WholeImageSelection().SetColor(0, 0, color1)
					// when
					window.DrawIntoBackBuffer()
					// then
					expected := make([]image.Color, zoom*zoom)
					for i := 0; i < len(expected); i++ {
						expected[i] = color1
					}
					assert.Equal(t, expected, framebufferPixels(window.ContextAPI(), 0, 0, int32(zoom), int32(zoom)))
				})
			}
		})

		t.Run("two windows", func(t *testing.T) {
			openGL, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
			window1, err := windowOfColor(openGL, color1)
			require.NoError(t, err)
			defer window1.Close()
			window2, err := windowOfColor(openGL, color2)
			require.NoError(t, err)
			defer window2.Close()
			// when
			window1.DrawIntoBackBuffer()
			// then
			expected := []image.Color{color1}
			assert.Equal(t, expected, framebufferPixels(window1.ContextAPI(), 0, 0, 1, 1))
			// when
			window2.DrawIntoBackBuffer()
			// then
			expected = []image.Color{color2}
			assert.Equal(t, expected, framebufferPixels(window2.ContextAPI(), 0, 0, 1, 1))
		})

		t.Run("two OpenGL instances", func(t *testing.T) {
			openGL1, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL1.Destroy()
			openGL2, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL2.Destroy()
			window1, err := windowOfColor(openGL1, color1)
			require.NoError(t, err)
			defer window1.Close()
			window2, err := windowOfColor(openGL2, color2)
			require.NoError(t, err)
			defer window2.Close()
			// when
			window1.DrawIntoBackBuffer()
			// then
			expected := []image.Color{color1}
			assert.Equal(t, expected, framebufferPixels(window1.ContextAPI(), 0, 0, 1, 1))
			// when
			window2.DrawIntoBackBuffer()
			// then
			expected = []image.Color{color2}
			assert.Equal(t, expected, framebufferPixels(window2.ContextAPI(), 0, 0, 1, 1))
		})
	})
}

func windowOfColor(openGL *glfw.OpenGL, color image.Color) (*glfw.Window, error) {
	window, err := openGL.OpenWindow(1, 1, glfw.NoDecorationHint())
	if err != nil {
		return nil, err
	}
	selection := window.Image().WholeImageSelection()
	selection.SetColor(0, 0, color)
	return window, err
}

func framebufferPixels(context gl2.API, x, y, width, height int32) []image.Color {
	size := (height - y) * (width - x)
	frameBuffer := make([]image.Color, size)
	context.ReadPixels(x, y, width, height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(frameBuffer))
	return frameBuffer
}

func TestWindow_PollKeyboardEvent(t *testing.T) {
	t.Run("should return EmptyEvent and false when there is no keyboard events", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		win, err := openGL.OpenWindow(1, 1)
		require.NoError(t, err)
		defer win.Close()
		// when
		event, ok := win.PollKeyboardEvent()
		// then
		assert.Equal(t, keyboard.EmptyEvent, event)
		assert.False(t, ok)
	})
}

func TestWindow_Zoom(t *testing.T) {
	t.Run("should return specified zoom for window", func(t *testing.T) {
		tests := map[string]struct {
			zoom         int
			expectedZoom int
		}{
			"zoom -1": {
				zoom:         -1,
				expectedZoom: 1,
			},
			"zoom 0": {
				zoom:         0,
				expectedZoom: 1,
			},
			"zoom 1": {
				zoom:         1,
				expectedZoom: 1,
			},
			"zoom 2": {
				zoom:         2,
				expectedZoom: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, err := glfw.NewOpenGL(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				win, err := openGL.OpenWindow(0, 0, glfw.Zoom(test.zoom))
				require.NoError(t, err)
				defer win.Close()
				// when
				zoom := win.Zoom()
				// expect
				assert.Equal(t, test.expectedZoom, zoom)
			})
		}
	})
}

func TestWindow_Image(t *testing.T) {
	t.Run("should provide screen image", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
		}{
			"1x2": {
				width:  1,
				height: 2,
			},
			"3x4": {
				width:  3,
				height: 4,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, err := glfw.NewOpenGL(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				win, err := openGL.OpenWindow(test.width, test.height, glfw.NoDecorationHint())
				require.NoError(t, err)
				defer win.Close()
				// when
				img := win.Image()
				// then
				require.NotNil(t, img)
				assert.Equal(t, test.width, img.Width())
				assert.Equal(t, test.height, img.Height())
			})
		}
	})
	t.Run("zoom should not affect the screen size", func(t *testing.T) {
		tests := map[string]struct {
			zoom int
		}{
			"zoom = -1": {
				zoom: 1,
			},
			"zoom = 0": {
				zoom: 0,
			},
			"zoom = 1": {
				zoom: 1,
			},
			"zoom = 2": {
				zoom: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, err := glfw.NewOpenGL(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				win, err := openGL.OpenWindow(640, 360, glfw.Zoom(test.zoom))
				require.NoError(t, err)
				// when
				screen := win.Image()
				// then
				assert.Equal(t, 640, screen.Width())
				assert.Equal(t, 360, screen.Height())
			})
		}
	})
	t.Run("initial screen is transparent", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		win, err := openGL.OpenWindow(1, 1, glfw.NoDecorationHint())
		require.NoError(t, err)
		transparent := image.RGBA(0, 0, 0, 0)
		// when
		img := win.Image()
		// then
		selection := img.WholeImageSelection()
		assert.Equal(t, transparent, selection.Color(0, 0))
	})
}

func TestWindow_ContextAPI(t *testing.T) {
	t.Run("should return context API", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win, _ := openGL.OpenWindow(1, 1)
		// when
		api := win.ContextAPI()
		// then
		assert.NotNil(t, api)
	})
}
