package opengl_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq"
	"github.com/jacekolszak/pixiq/opengl"
)

var mainThreadLoop *opengl.MainThreadLoop

func TestMain(m *testing.M) {
	opengl.StartMainThreadLoop(func(main *opengl.MainThreadLoop) {
		mainThreadLoop = main
		exit := m.Run()
		mainThreadLoop.StopEventually()
		os.Exit(exit)
	})
}

func TestNew(t *testing.T) {
	t.Run("should create OpenGL using supplied MainThreadLoop", func(t *testing.T) {
		// when
		gl := opengl.New(mainThreadLoop)
		images := gl.AcceleratedImages()
		windows := gl.SystemWindows()
		// then
		assert.NotNil(t, images)
		assert.NotNil(t, windows)
	})
}

func TestTextures_New(t *testing.T) {
	t.Skip()
	t.Run("should create AcceleratedImage", func(t *testing.T) {
		gl := opengl.New(mainThreadLoop)
		images := gl.AcceleratedImages()
		// when
		image := images.New(0, 0)
		// then
		assert.NotNil(t, image)
	})
}

func TestTexture_Upload(t *testing.T) {
	t.Skip()
	t.Run("should upload pixels", func(t *testing.T) {
		gl := opengl.New(mainThreadLoop)
		images := gl.AcceleratedImages()
		image := images.New(1, 1)
		color := pixiq.RGBA(10, 20, 30, 40)
		input := []pixiq.Color{color}
		// when
		image.Upload(input)
		// then
		output := make([]pixiq.Color, 1)
		image.Download(output)
		assert.Equal(t, input, output)
	})
}

func TestGlfwWindows_Open(t *testing.T) {
	t.Skip()
	t.Run("should open window", func(t *testing.T) {
		gl := opengl.New(mainThreadLoop)
		windows := gl.SystemWindows()
		// when
		window := windows.Open(640, 360)
		// then
		assert.NotNil(t, window)
	})
}

func TestGlfwWindow_Draw(t *testing.T) {
	t.Skip()
	t.Run("should draw image inside window", func(t *testing.T) {
		gl := opengl.New(mainThreadLoop)
		windows := gl.SystemWindows()
		window := windows.Open(1, 1)
		images := pixiq.NewImages(gl.AcceleratedImages())
		image := images.New(1, 1)
		color := pixiq.RGBA(10, 20, 30, 40)
		image.WholeImageSelection().SetColor(0, 0, color)
		// when
		window.Draw(image)
		// then
		// TODO verify framebuffer 0 - use readPixels or something similar
	})
}
