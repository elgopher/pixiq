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
		var images pixiq.AcceleratedImages
		var windows pixiq.SystemWindows
		// when
		gl := opengl.New(mainThreadLoop)
		images = gl.AcceleratedImages()
		windows = gl.SystemWindows()
		// then
		assert.NotNil(t, images)
		assert.NotNil(t, windows)
	})
}
