package opengl_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/image"
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
		// then
		assert.NotNil(t, openGL)
	})
	t.Run("should create 2 objects working at the same time", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			openGL := opengl.New(mainThreadLoop)
			defer openGL.Destroy()
		}
	})
}

func TestOpenGL_NewImage(t *testing.T) {
	t.Run("should create Image", func(t *testing.T) {
		tests := map[string]struct {
			width  int
			height int
		}{
			"0x0": {
				width:  0,
				height: 0,
			},
			"1x2": {
				width:  1,
				height: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				img := openGL.NewImage(test.width, test.height)
				// then
				assert.NotNil(t, img)
				assert.Equal(t, test.width, img.Width())
				assert.Equal(t, test.height, img.Height())
			})
		}
	})
}

func TestOpenGL_NewAcceleratedImage(t *testing.T) {
	t.Run("should create AcceleratedImage", func(t *testing.T) {
		openGL := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		// when
		img := openGL.NewAcceleratedImage(0, 0)
		// then
		assert.NotNil(t, img)
	})
}

func TestTexture_Upload(t *testing.T) {
	color1 := image.RGBA(10, 20, 30, 40)
	color2 := image.RGBA(50, 60, 70, 80)
	color3 := image.RGBA(90, 100, 110, 120)
	color4 := image.RGBA(130, 140, 150, 160)

	t.Run("should upload pixels", func(t *testing.T) {
		tests := map[string]struct {
			width, height int
			inputColors   []image.Color
		}{
			"1x1": {
				width:       1,
				height:      1,
				inputColors: []image.Color{color1},
			},
			"2x1": {
				width:       2,
				height:      1,
				inputColors: []image.Color{color1, color2},
			},
			"1x2": {
				width:       1,
				height:      2,
				inputColors: []image.Color{color1, color2},
			},
			"2x2": {
				width:       2,
				height:      2,
				inputColors: []image.Color{color1, color2, color3, color4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img := openGL.NewAcceleratedImage(test.width, test.height)
				// when
				img.Upload(test.inputColors)
				// then
				output := make([]image.Color, len(test.inputColors))
				img.Download(output)
				assert.Equal(t, test.inputColors, output)
			})
		}
	})
	t.Run("2 OpenGL contexts", func(t *testing.T) {
		gl1 := opengl.New(mainThreadLoop)
		defer gl1.Destroy()
		gl2 := opengl.New(mainThreadLoop)
		defer gl2.Destroy()
		img1 := gl1.NewAcceleratedImage(1, 1)
		img2 := gl2.NewAcceleratedImage(1, 1)
		// when
		img1.Upload([]image.Color{color1})
		img2.Upload([]image.Color{color2})
		// then
		output := make([]image.Color, 1)
		img1.Download(output)
		assert.Equal(t, []image.Color{color1}, output)
		// and
		img2.Download(output)
		assert.Equal(t, []image.Color{color2}, output)
	})
}

func TestRun(t *testing.T) {
	t.Run("should run provided callback", func(t *testing.T) {
		var callbackExecuted bool
		mainThreadLoop.Execute(func() {
			opengl.Run(func(gl *opengl.OpenGL) {
				callbackExecuted = true
			})
		})
		assert.True(t, callbackExecuted)
	})
	t.Run("should start a MainThreadLoop and create OpenGL object", func(t *testing.T) {
		var (
			actualGL *opengl.OpenGL
		)
		mainThreadLoop.Execute(func() {
			opengl.Run(func(gl *opengl.OpenGL) {
				actualGL = gl
			})
		})
		assert.NotNil(t, actualGL)
	})
}
