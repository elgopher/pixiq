package glfw_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gl2 "github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/keyboard"
	"github.com/jacekolszak/pixiq/mouse"
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
			window.Screen().SetColor(0, 0, color1)
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
			window.Screen().SetColor(0, 0, color1)
			window.Screen().SetColor(0, 1, color2)
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
			window.Screen().SetColor(0, 0, color1)
			window.Screen().SetColor(1, 0, color2)
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
			screen := window.Screen()
			screen.SetColor(0, 0, color1)
			screen.SetColor(1, 0, color2)
			screen.SetColor(0, 1, color3)
			screen.SetColor(1, 1, color4)
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
					window.Screen().SetColor(0, 0, color1)
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
					window.Screen().SetColor(0, 0, color1)
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

		t.Run("should draw screen despite the state of OpenGL context", func(t *testing.T) {
			tests := map[string]func(ctx *gl2.Context, win *glfw.Window){
				"BlendFunc": func(ctx *gl2.Context, win *glfw.Window) {
					ctx.API().Enable(gl.BLEND)
					ctx.API().BlendFunc(0, 0)
				},
				"BindFramebuffer": func(ctx *gl2.Context, win *glfw.Window) {
					var fb uint32
					ctx.API().GenFramebuffers(1, &fb)
					ctx.API().BindFramebuffer(gl.FRAMEBUFFER, fb)
				},
				"Viewport": func(ctx *gl2.Context, win *glfw.Window) {
					ctx.API().Viewport(0, 0, 0, 0)
				},
				"Scissor": func(ctx *gl2.Context, win *glfw.Window) {
					ctx.API().Enable(gl.SCISSOR_TEST)
					ctx.API().Scissor(0, 0, 0, 0)
				},
				"BindTexture": func(ctx *gl2.Context, win *glfw.Window) {
					img := ctx.NewAcceleratedImage(1, 1)
					img.Upload([]image.Color{image.RGBA(0, 0, 0, 0)})
					texture := img.TextureID()
					ctx.API().BindTexture(gl.TEXTURE_2D, texture)
				},
				"BindTexture when screen is uploaded": func(ctx *gl2.Context, win *glfw.Window) {
					win.Screen().Image().Upload()
					img := ctx.NewAcceleratedImage(1, 1)
					img.Upload([]image.Color{image.RGBA(0, 0, 0, 0)})
					texture := img.TextureID()
					ctx.API().BindTexture(gl.TEXTURE_2D, texture)
				},
				"UseProgram": func(ctx *gl2.Context, win *glfw.Window) {
					program := ctx.API().CreateProgram()
					ctx.API().UseProgram(program)
				},
			}
			for name, testFunction := range tests {
				t.Run(name, func(t *testing.T) {
					openGL, _ := glfw.NewOpenGL(mainThreadLoop)
					defer openGL.Destroy()
					win, _ := windowOfColor(openGL, color1)
					defer win.Close()
					// when
					testFunction(openGL.Context(), win)
					win.DrawIntoBackBuffer()
					// then
					expected := []image.Color{color1}
					assert.Equal(t, expected, framebufferPixels(win.ContextAPI(), 0, 0, 1, 1))
				})
			}

		})
	})

	t.Run("should panic for closed window", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win, _ := openGL.OpenWindow(1, 1)
		win.Close()
		assert.Panics(t, func() {
			// when
			win.DrawIntoBackBuffer()
		})
	})

	t.Run("should draw perfect pixels when window size is not a zoom multiplication", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		videoMode := display.VideoModes()[0] // TODO avoid using 1:1, 2:1, 1:2 ratios
		zoom := videoMode.Height() / 2
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		window, err := openGL.OpenFullScreenWindow(videoMode, glfw.Zoom(zoom))
		require.NoError(t, err)
		defer window.Close()
		c00 := image.RGBA(10, 20, 30, 40)
		c10 := image.RGBA(50, 60, 70, 80)
		c01 := image.RGBA(90, 100, 110, 120)
		c11 := image.RGBA(130, 140, 150, 160)
		window.Screen().SetColor(0, 0, c00)
		window.Screen().SetColor(1, 0, c10)
		window.Screen().SetColor(0, 1, c01)
		window.Screen().SetColor(1, 1, c11)
		// when
		window.DrawIntoBackBuffer() // TODO This does not work because size was not yet updated
		// then
		fb := framebufferPixels(window.ContextAPI(), 0, 0, 1, 1)
		assert.Equal(t, c01, fb[0])
		fb = framebufferPixels(window.ContextAPI(), int32(videoMode.Width()/2), 0, 1, 1)
		assert.Equal(t, c11, fb[0])
		fb = framebufferPixels(window.ContextAPI(), 0, int32(videoMode.Height()/2), 1, 1)
		assert.Equal(t, c00, fb[0])
		fb = framebufferPixels(window.ContextAPI(), int32(videoMode.Width()/2), int32(videoMode.Height()/2), 1, 1)
		assert.Equal(t, c10, fb[0])
	})

}

func TestWindow_Draw(t *testing.T) {
	t.Run("should panic for closed window", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win, _ := openGL.OpenWindow(1, 1)
		win.Close()
		assert.Panics(t, func() {
			// when
			win.Draw()
		})
	})
}

func TestWindow_SwapBuffers(t *testing.T) {
	t.Run("should panic for closed window", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win, _ := openGL.OpenWindow(1, 1)
		win.Close()
		assert.Panics(t, func() {
			// when
			win.SwapBuffers()
		})
	})
}

func TestWindow_Close(t *testing.T) {
	t.Run("second Close does nothing", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win, _ := openGL.OpenWindow(1, 1)
		win.Close()
		// when
		win.Close()
	})
	t.Run("second Close on a second open window does nothing", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win1, _ := openGL.OpenWindow(1, 1)
		defer win1.Close()
		win2, _ := openGL.OpenWindow(1, 1)
		win2.Close()
		// when
		win2.Close()
	})
}

func TestWindow_ShouldClose(t *testing.T) {
	t.Run("ShouldClose on closed window returns false", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win, _ := openGL.OpenWindow(1, 1)
		win.Close()
		// when
		shouldClose := win.ShouldClose()
		assert.False(t, shouldClose)
	})
	t.Run("ShouldClose on a second closed window returns false", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		win1, _ := openGL.OpenWindow(1, 1)
		defer win1.Close()
		win2, _ := openGL.OpenWindow(1, 1)
		win2.Close()
		// when
		shouldClose := win2.ShouldClose()
		assert.False(t, shouldClose)
	})
}

func windowOfColor(openGL *glfw.OpenGL, color image.Color) (*glfw.Window, error) {
	window, err := openGL.OpenWindow(1, 1, glfw.NoDecorationHint())
	if err != nil {
		return nil, err
	}
	window.Screen().SetColor(0, 0, color)
	return window, err
}

func framebufferPixels(context gl2.API, x, y, width, height int32) []image.Color {
	size := height * width
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

func TestWindow_PollMouseEvent(t *testing.T) {
	t.Run("should return EmptyEvent and false when there is no mouse events", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		win, err := openGL.OpenWindow(1, 1)
		require.NoError(t, err)
		defer win.Close()
		// when
		_, _ = win.PollMouseEvent() // mouse.MoveEvent is always returned first
		event, ok := win.PollMouseEvent()
		// then
		assert.Equal(t, mouse.EmptyEvent, event)
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

func TestWindow_Screen(t *testing.T) {
	t.Run("should provide screen selection", func(t *testing.T) {
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
				screen := win.Screen()
				// then
				assert.Equal(t, 0, screen.ImageX())
				assert.Equal(t, 0, screen.ImageY())
				assert.Equal(t, test.width, screen.Width())
				assert.Equal(t, test.height, screen.Height())
				// and
				require.NotNil(t, screen.Image())
				assert.Equal(t, test.width, screen.Image().Width())
				assert.Equal(t, test.height, screen.Image().Height())
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
				screen := win.Screen()
				// then
				assert.Equal(t, 640, screen.Width())
				assert.Equal(t, 360, screen.Height())
				// and
				assert.Equal(t, 640, screen.Image().Width())
				assert.Equal(t, 360, screen.Image().Height())
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
		screen := win.Screen()
		// then
		assert.Equal(t, transparent, screen.Color(0, 0))
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

func TestWindow_SetCursor(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()

	t.Run("should panic for nil cursor", func(t *testing.T) {
		win, _ := openGL.OpenWindow(1, 1)
		defer win.Close()
		assert.Panics(t, func() {
			win.SetCursor(nil)
		})
	})

	t.Run("should set custom cursor for window", func(t *testing.T) {
		win, _ := openGL.OpenWindow(1, 1)
		defer win.Close()
		cursorSelection := openGL.NewImage(1, 1).WholeImageSelection()
		cursor := openGL.NewCursor(cursorSelection)
		// when
		win.SetCursor(cursor)
	})

	t.Run("should set standard cursor for window", func(t *testing.T) {
		win, _ := openGL.OpenWindow(1, 1)
		defer win.Close()
		cursor := openGL.NewStandardCursor(glfw.Arrow)
		// when
		win.SetCursor(cursor)
	})
}

func TestWindow_Resize(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()

	t.Run("should eventually resize window", func(t *testing.T) {
		tests := map[string]struct {
			newWidth, newHeight, newZoom                int
			expectedWidth, expectedHeight, expectedZoom int
		}{
			"320x180x1": {
				newWidth:       320,
				newHeight:      180,
				newZoom:        1,
				expectedWidth:  320,
				expectedHeight: 180,
				expectedZoom:   1,
			},
			"640x360x2": {
				newWidth:       640,
				newHeight:      360,
				newZoom:        2,
				expectedWidth:  1280,
				expectedHeight: 720,
				expectedZoom:   2,
			},
			"320x180x0": {
				newWidth:       320,
				newHeight:      180,
				newZoom:        0,
				expectedWidth:  320,
				expectedHeight: 180,
				expectedZoom:   1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				window, _ := openGL.OpenWindow(640, 360)
				defer window.Close()
				// when
				window.Resize(test.newWidth, test.newHeight, test.newZoom)
				// then
				assert.Eventually(t, func() bool {
					return window.Width() == test.expectedWidth &&
						window.Height() == test.expectedHeight &&
						window.Zoom() == test.expectedZoom
				}, 1*time.Second, 10*time.Millisecond)
			})
		}
	})

	t.Run("should resize screen", func(t *testing.T) {
		originalWidth, originalHeight := 640, 320
		tests := map[string]struct {
			newWidth, newHeight, newZoom int
		}{
			"nothing has changed": {
				newWidth:  originalWidth,
				newHeight: originalHeight,
				newZoom:   1,
			},
			"zoom changed to 2": {
				newWidth:  originalWidth,
				newHeight: originalHeight,
				newZoom:   2,
			},
			"zoom changed to 0": {
				newWidth:  originalWidth,
				newHeight: originalHeight,
				newZoom:   0,
			},
			"half the size": {
				newWidth:  originalWidth / 2,
				newHeight: originalHeight / 2,
				newZoom:   1,
			},
			"double the size": {
				newWidth:  originalWidth * 2,
				newHeight: originalHeight * 2,
				newZoom:   1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				window, _ := openGL.OpenWindow(originalWidth, originalHeight)
				defer window.Close()
				fillWithColors(window.Screen())
				originalImage := openGL.NewImage(originalWidth, originalHeight)
				defer originalImage.Delete()
				originalSelection := originalImage.WholeImageSelection()
				copySourceTo(window.Screen(), originalSelection)
				// when
				window.Resize(test.newWidth, test.newHeight, test.newZoom)
				// then
				resizedScreen := window.Screen()
				assert.Equal(t, test.newWidth, resizedScreen.Width())
				assert.Equal(t, test.newHeight, resizedScreen.Height())
				// and
				assertSelectionEqual(t, originalSelection, resizedScreen)
			})
		}

	})
}

func assertSelectionEqual(t *testing.T, expected image.Selection, actual image.Selection) {
	for y := 0; y < actual.Height(); y++ {
		for x := 0; x < actual.Width(); x++ {
			assert.Equal(t, expected.Color(x, y), actual.Color(x, y))
		}
	}
}

func copySourceTo(source image.Selection, target image.Selection) {
	for y := 0; y < source.Height(); y++ {
		for x := 0; x < source.Width(); x++ {
			target.SetColor(x, y, source.Color(x, y))
		}
	}
}

func fillWithColors(selection image.Selection) {
	r, g, b, a := 10, 20, 30, 40
	for y := 0; y < selection.Height(); y++ {
		for x := 0; x < selection.Width(); x++ {
			selection.SetColor(x, y, image.RGBAi(r, g, b, a))
			if r++; r > 255 {
				r = 0
			}
			if g++; g > 255 {
				g = 0
			}
			if b++; b > 255 {
				b = 0
			}
			if a++; a > 255 {
				a = 0
			}
		}
	}
}

func TestWindow_SetPosition(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()

	t.Run("should eventually set position of window", func(t *testing.T) {
		window, _ := openGL.OpenWindow(640, 360)
		defer window.Close()
		// when
		newX := 100
		newY := 200
		window.SetPosition(newX, newY)
		// then
		assert.Eventually(t, func() bool {
			return newX == window.X() &&
				newY == window.Y()
		}, 1*time.Second, 10*time.Millisecond)
	})
}

func TestWindow_SetDecorationHint(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()

	t.Run("should hide decorations", func(t *testing.T) {
		window, _ := openGL.OpenWindow(640, 360)
		defer window.Close()
		// when
		window.SetDecorationHint(true)
		// then
		assert.True(t, window.Decorated())
	})

	t.Run("should show decorations", func(t *testing.T) {
		window, _ := openGL.OpenWindow(640, 360)
		defer window.Close()
		// when
		window.SetDecorationHint(false)
		// then
		assert.False(t, window.Decorated())
	})
}

func TestWindow_EnterFullScreen(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()

	t.Run("should enter full screen using first video mode", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		// current video mode on MacOS is returning not supported full screen video mode
		videoMode := display.VideoModes()[0]
		window, _ := openGL.OpenWindow(320, 200)
		defer window.Close()
		// when
		window.EnterFullScreen(videoMode, 2)
		// then
		var (
			expectedWindowWidth  = videoMode.Width()
			expectedWindowHeight = videoMode.Height()
			expectedScreenWidth  = expectedWindowWidth / 2
			expectedScreenHeight = expectedWindowHeight / 2
		)
		assert.Eventually(t, func() bool {
			screen := window.Screen()
			return expectedWindowWidth == window.Width() &&
				expectedWindowHeight == window.Height() &&
				expectedScreenWidth == screen.Width() &&
				expectedScreenHeight == screen.Height()
		}, 1*time.Second, 10*time.Millisecond)
	})
}

func TestWindow_ExitFullScreen(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()

	t.Run("should exit full screen", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		videoMode := display.VideoMode()
		window, _ := openGL.OpenFullScreenWindow(videoMode, glfw.Zoom(2))
		defer window.Close()
		// when
		window.ExitFullScreen()
		// then
		assert.Eventually(t, func() bool {
			return videoMode.Width() == window.Width() &&
				videoMode.Height() == window.Height() &&
				2 == window.Zoom() &&
				videoMode.Width()/2 == window.Screen().Width() &&
				videoMode.Height()/2 == window.Screen().Height()
		}, 1*time.Second, 10*time.Millisecond)
	})

	t.Run("should exit full screen after executing EnterFullScreen", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		videoMode := display.VideoMode()
		window, _ := openGL.OpenWindow(320, 200, glfw.Zoom(2))
		defer window.Close()
		x, y := window.X(), window.Y()
		window.EnterFullScreen(videoMode, 1)
		// when
		window.ExitFullScreen()
		// then
		assert.Eventually(t, func() bool {
			return 640 == window.Width() &&
				400 == window.Height() &&
				2 == window.Zoom() &&
				x == window.X() &&
				y == window.Y() &&
				320 == window.Screen().Width() &&
				200 == window.Screen().Height()
		}, 1*time.Second, 10*time.Millisecond)
	})
}

func TestWindow_SetAutoIconifyHint(t *testing.T) {
	openGL, _ := glfw.NewOpenGL(mainThreadLoop)
	defer openGL.Destroy()
	displays, _ := glfw.Displays(mainThreadLoop)
	display, _ := displays.Primary()
	videoMode := display.VideoMode()

	t.Run("should no iconify full screen window on focus lost", func(t *testing.T) {
		window, _ := openGL.OpenFullScreenWindow(videoMode)
		defer window.Close()
		// when
		window.SetAutoIconifyHint(false)
		// then
		assert.False(t, window.AutoIconify())
	})

	t.Run("should iconify full screen window on focus lost", func(t *testing.T) {
		window, _ := openGL.OpenFullScreenWindow(videoMode)
		defer window.Close()
		// when
		window.SetAutoIconifyHint(true)
		// then
		assert.True(t, window.AutoIconify())
	})

}
