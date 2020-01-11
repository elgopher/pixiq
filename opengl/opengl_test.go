package opengl_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/opengl"
)

var mainThreadLoop *opengl.MainThreadLoop

func TestMain(m *testing.M) {
	var exit int
	opengl.StartMainThreadLoop(func(main *opengl.MainThreadLoop) {
		mainThreadLoop = main
		exit = m.Run()
	})
	os.Exit(exit)
}

func TestNew(t *testing.T) {
	t.Run("should panic when MainThreadLoop is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			opengl.New(nil)
		})
	})
	t.Run("should create OpenGL using supplied MainThreadLoop", func(t *testing.T) {
		// when
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		images := openGL.AcceleratedImages()
		windows := openGL.Windows()
		// then
		assert.NotNil(t, images)
		assert.NotNil(t, windows)
	})
}

func TestTextures_New(t *testing.T) {
	t.Run("should create AcceleratedImage", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		images := openGL.AcceleratedImages()
		// when
		image := images.New(0, 0)
		// then
		assert.NotNil(t, image)
	})
}

func TestTexture_Upload(t *testing.T) {
	t.Run("should upload pixels", func(t *testing.T) {
		color1 := pixiq.RGBA(10, 20, 30, 40)
		color2 := pixiq.RGBA(50, 60, 70, 80)
		color3 := pixiq.RGBA(90, 100, 110, 120)
		color4 := pixiq.RGBA(130, 140, 150, 160)

		tests := map[string]struct {
			width, height int
			inputColors   []pixiq.Color
		}{
			"1x1": {
				width:       1,
				height:      1,
				inputColors: []pixiq.Color{color1},
			},
			"2x1": {
				width:       2,
				height:      1,
				inputColors: []pixiq.Color{color1, color2},
			},
			"1x2": {
				width:       1,
				height:      2,
				inputColors: []pixiq.Color{color1, color2},
			},
			"2x2": {
				width:       2,
				height:      2,
				inputColors: []pixiq.Color{color1, color2, color3, color4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				images := openGL.AcceleratedImages()
				image := images.New(test.width, test.height)
				// when
				image.Upload(test.inputColors)
				// then
				output := make([]pixiq.Color, len(test.inputColors))
				image.Download(output)
				assert.Equal(t, test.inputColors, output)
			})
		}
	})
}

func TestWindows_Open(t *testing.T) {
	t.Run("should clamp width to platform-specific minimum if negative", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		windows := openGL.Windows()
		// when
		win := windows.Open(-1, 0)
		defer win.Close()
		// then
		require.NotNil(t, win)
		assert.GreaterOrEqual(t, win.Width(), 0)
	})
	t.Run("should clamp height to platform-specific minimum if negative", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		windows := openGL.Windows()
		// when
		win := windows.Open(0, -1)
		defer win.Close()
		// then
		require.NotNil(t, win)
		assert.GreaterOrEqual(t, win.Height(), 0)
	})
	t.Run("should open Window", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		windows := openGL.Windows()
		// when
		win := windows.Open(640, 360)
		defer win.Close()
		// then
		require.NotNil(t, win)
		assert.Equal(t, 640, win.Width())
		assert.Equal(t, 360, win.Height())
	})
	t.Run("should open two windows at the same time", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		windows := openGL.Windows()
		// when
		win1 := windows.Open(640, 360)
		defer win1.Close()
		win2 := windows.Open(320, 180)
		defer win2.Close()
		// then
		require.NotNil(t, win1)
		assert.Equal(t, 640, win1.Width())
		assert.Equal(t, 360, win1.Height())
		require.NotNil(t, win2)
		assert.Equal(t, 320, win2.Width())
		assert.Equal(t, 180, win2.Height())
	})
	t.Run("should open another Window after first one was closed", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		windows := openGL.Windows()
		win1 := windows.Open(640, 360)
		win1.Close()
		// when
		win2 := windows.Open(320, 180)
		defer win2.Close()
		// then
		require.NotNil(t, win2)
		assert.Equal(t, 320, win2.Width())
		assert.Equal(t, 180, win2.Height())
	})
	t.Run("should skip nil option", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		windows := openGL.Windows()
		// when
		win := windows.Open(0, 0, nil)
		defer win.Close()
	})
	t.Run("zoom should not affect the width and height", func(t *testing.T) {
		for zoom := -1; zoom <= 2; zoom++ {
			name := fmt.Sprintf("zoom=%d", zoom)
			t.Run(name, func(t *testing.T) {
				openGL := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				windows := openGL.Windows()
				// when
				win := windows.Open(640, 360, opengl.Zoom(zoom))
				defer win.Close()
				// then
				require.NotNil(t, win)
				assert.Equal(t, 640, win.Width())
				assert.Equal(t, 360, win.Height())
			})
		}
	})
}

func TestWindow_Draw(t *testing.T) {
	t.Run("should draw screen image", func(t *testing.T) {
		color1 := pixiq.RGBA(10, 20, 30, 40)
		color2 := pixiq.RGBA(50, 60, 70, 80)
		color3 := pixiq.RGBA(90, 100, 110, 120)
		color4 := pixiq.RGBA(130, 140, 150, 160)

		t.Run("1x1", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			defer openGL.Destroy()
			windows := openGL.Windows()
			window := windows.Open(1, 1, opengl.NoDecorationHint())
			defer window.Close()
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(1, 1)
			image.WholeImageSelection().SetColor(0, 0, color1)
			// when
			window.Draw(image)
			// then
			expected := []pixiq.Color{color1}
			assert.Equal(t, expected, framebufferPixels(0, 0, 1, 1))
		})
		t.Run("1x2", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			defer openGL.Destroy()
			windows := openGL.Windows()
			window := windows.Open(1, 2, opengl.NoDecorationHint())
			defer window.Close()
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(1, 2)
			image.WholeImageSelection().SetColor(0, 0, color1)
			image.WholeImageSelection().SetColor(0, 1, color2)
			// when
			window.Draw(image)
			// then
			expected := []pixiq.Color{color2, color1}
			assert.Equal(t, expected, framebufferPixels(0, 0, 1, 2))
		})
		t.Run("2x1", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			defer openGL.Destroy()
			windows := openGL.Windows()
			window := windows.Open(2, 1, opengl.NoDecorationHint())
			defer window.Close()
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(2, 1)
			image.WholeImageSelection().SetColor(0, 0, color1)
			image.WholeImageSelection().SetColor(1, 0, color2)
			// when
			window.Draw(image)
			// then
			expected := []pixiq.Color{color1, color2}
			assert.Equal(t, expected, framebufferPixels(0, 0, 2, 1))
		})
		t.Run("2x2", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			defer openGL.Destroy()
			windows := openGL.Windows()
			window := windows.Open(2, 2, opengl.NoDecorationHint())
			defer window.Close()
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(2, 2)
			selection := image.WholeImageSelection()
			selection.SetColor(0, 0, color1)
			selection.SetColor(1, 0, color2)
			selection.SetColor(0, 1, color3)
			selection.SetColor(1, 1, color4)
			// when
			window.Draw(image)
			// then
			expected := []pixiq.Color{color3, color4, color1, color2}
			assert.Equal(t, expected, framebufferPixels(0, 0, 2, 2))
		})

		t.Run("zoom < 1 should not change the framebuffer size", func(t *testing.T) {
			for zoom := -1; zoom < 1; zoom++ {
				name := fmt.Sprintf("zoom=%d", zoom)
				t.Run(name, func(t *testing.T) {
					openGL := opengl.New(mainThreadLoop)
					defer openGL.Destroy()
					windows := openGL.Windows()
					window := windows.Open(1, 1, opengl.NoDecorationHint(), opengl.Zoom(zoom))
					defer window.Close()
					images := pixiq.NewImages(openGL.AcceleratedImages())
					image := images.New(1, 1)
					image.WholeImageSelection().SetColor(0, 0, color1)
					// when
					window.Draw(image)
					// then
					expected := []pixiq.Color{color1}
					assert.Equal(t, expected, framebufferPixels(0, 0, 1, 1))
				})
			}
		})

		t.Run("zoom > 1 should make framebuffer zoom times bigger", func(t *testing.T) {
			for zoom := 2; zoom < 4; zoom++ {
				name := fmt.Sprintf("zoom=%d", zoom)
				t.Run(name, func(t *testing.T) {
					openGL := opengl.New(mainThreadLoop)
					defer openGL.Destroy()
					windows := openGL.Windows()
					window := windows.Open(1, 1, opengl.NoDecorationHint(), opengl.Zoom(zoom))
					defer window.Close()
					images := pixiq.NewImages(openGL.AcceleratedImages())
					image := images.New(1, 1)
					image.WholeImageSelection().SetColor(0, 0, color1)
					// when
					window.Draw(image)
					// then
					expected := make([]pixiq.Color, zoom*zoom)
					for i := 0; i < len(expected); i++ {
						expected[i] = color1
					}
					assert.Equal(t, expected, framebufferPixels(0, 0, int32(zoom), int32(zoom)))
				})
			}
		})
	})
}

func framebufferPixels(x, y, width, height int32) []pixiq.Color {
	size := (height - y) * (width - x)
	frameBuffer := make([]pixiq.Color, size)
	mainThreadLoop.Execute(func() {
		gl.ReadPixels(x, y, width, height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(frameBuffer))
	})
	return frameBuffer
}

func TestRun(t *testing.T) {
	t.Run("should run provided callback", func(t *testing.T) {
		var callbackExecuted bool
		mainThreadLoop.Execute(func() {
			opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
				callbackExecuted = true
			})
		})
		assert.True(t, callbackExecuted)
	})
	t.Run("should create pixiq objects using OpenGL acceleration and windows", func(t *testing.T) {
		var (
			actualGL     *opengl.OpenGL
			actualImages *pixiq.Images
			actualLoops  *pixiq.ScreenLoops
		)
		mainThreadLoop.Execute(func() {
			opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, loops *pixiq.ScreenLoops) {
				actualGL = gl
				actualImages = images
				actualLoops = loops
			})
		})
		assert.NotNil(t, actualGL)
		assert.NotNil(t, actualGL.Windows())
		assert.NotNil(t, actualGL.AcceleratedImages())
		assert.NotNil(t, actualImages)
		assert.NotNil(t, actualLoops)
	})

}

func TestWindow_Poll(t *testing.T) {
	t.Run("should return EmptyEvent and false when there is no keyboard events", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		win := openGL.Windows().Open(1, 1)
		defer win.Close()
		// when
		event, ok := win.Poll()
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
				openGL := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				win := openGL.Windows().Open(0, 0, opengl.Zoom(test.zoom))
				defer win.Close()
				// when
				zoom := win.Zoom()
				// expect
				assert.Equal(t, test.expectedZoom, zoom)
			})
		}
	})
}
