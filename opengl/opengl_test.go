package opengl_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq"
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
	t.Run("should create OpenGL using supplied MainThreadLoop", func(t *testing.T) {
		// when
		openGL := opengl.New(mainThreadLoop)
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

func TestGlfwWindows_Open(t *testing.T) {
	t.Run("should clamp width to platform-specific minimum if negative", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		windows := openGL.Windows()
		// when
		win := windows.Open(-1, 0)
		require.NotNil(t, win)
		assert.GreaterOrEqual(t, win.Width(), 0)
	})
	t.Run("should clamp height to platform-specific minimum if negative", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		windows := openGL.Windows()
		// when
		win := windows.Open(0, -1)
		require.NotNil(t, win)
		assert.GreaterOrEqual(t, win.Height(), 0)
	})
	t.Run("should open Window", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		windows := openGL.Windows()
		// when
		win := windows.Open(640, 360)
		require.NotNil(t, win)
		assert.Equal(t, 640, win.Width())
		assert.Equal(t, 360, win.Height())
	})
	t.Run("should open two windows at the same time", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		windows := openGL.Windows()
		// when
		win1 := windows.Open(640, 360)
		win2 := windows.Open(320, 180)
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
		windows := openGL.Windows()
		win1 := windows.Open(640, 360)
		win1.Close()
		// when
		win2 := windows.Open(320, 180)
		require.NotNil(t, win2)
		assert.Equal(t, 320, win2.Width())
		assert.Equal(t, 180, win2.Height())
	})
}

func TestGlfwWindow_Draw(t *testing.T) {
	t.Run("should draw screen image", func(t *testing.T) {
		color1 := pixiq.RGBA(10, 20, 30, 40)
		color2 := pixiq.RGBA(50, 60, 70, 80)
		color3 := pixiq.RGBA(90, 100, 110, 120)
		color4 := pixiq.RGBA(130, 140, 150, 160)

		t.Run("1x1", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			windows := openGL.Windows()
			window := windows.Open(1, 1, opengl.NoDecoration{})
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(1, 1)
			image.WholeImageSelection().SetColor(0, 0, color1)
			// when
			window.Draw(image)
			// then
			assert.Equal(t, []pixiq.Color{color1}, framebufferPixels(0, 0, 1, 1))
		})
		t.Run("1x2", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			windows := openGL.Windows()
			window := windows.Open(1, 2, opengl.NoDecoration{})
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(1, 2)
			image.WholeImageSelection().SetColor(0, 0, color1)
			image.WholeImageSelection().SetColor(0, 1, color2)
			// when
			window.Draw(image)
			// then
			assert.Equal(t, []pixiq.Color{color2, color1}, framebufferPixels(0, 0, 1, 2))
		})
		t.Run("2x1", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			windows := openGL.Windows()
			window := windows.Open(2, 1, opengl.NoDecoration{})
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(2, 1)
			image.WholeImageSelection().SetColor(0, 0, color1)
			image.WholeImageSelection().SetColor(1, 0, color2)
			// when
			window.Draw(image)
			// then
			assert.Equal(t, []pixiq.Color{color1, color2}, framebufferPixels(0, 0, 2, 1))
		})
		t.Run("2x2", func(t *testing.T) {
			openGL := opengl.New(mainThreadLoop)
			windows := openGL.Windows()
			window := windows.Open(2, 2, opengl.NoDecoration{})
			images := pixiq.NewImages(openGL.AcceleratedImages())
			image := images.New(2, 2)
			image.WholeImageSelection().SetColor(0, 0, color1)
			image.WholeImageSelection().SetColor(1, 0, color2)
			image.WholeImageSelection().SetColor(0, 1, color3)
			image.WholeImageSelection().SetColor(1, 1, color4)
			// when
			window.Draw(image)
			// then
			assert.Equal(t, []pixiq.Color{color3, color4, color1, color2}, framebufferPixels(0, 0, 2, 2))
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
			opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, screens *pixiq.Screens) {
				callbackExecuted = true
			})
		})
		assert.True(t, callbackExecuted)
	})
	t.Run("should create pixiq objects using OpenGL acceleration and windows", func(t *testing.T) {
		var (
			actualGL      *opengl.OpenGL
			actualImages  *pixiq.Images
			actualScreens *pixiq.Screens
		)
		mainThreadLoop.Execute(func() {
			opengl.Run(func(gl *opengl.OpenGL, images *pixiq.Images, screens *pixiq.Screens) {
				actualGL = gl
				actualImages = images
				actualScreens = screens
			})
		})
		assert.NotNil(t, actualGL)
		assert.NotNil(t, actualImages)
		assert.NotNil(t, actualScreens)
	})

}
