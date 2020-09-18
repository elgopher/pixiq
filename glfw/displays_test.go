package glfw_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// All tests here are rather plumbing tests than real tests
// verifying if glfw package integrates properly with GLFW library.

func TestDisplays(t *testing.T) {
	t.Run("should create API object", func(t *testing.T) {
		displays, err := glfw.Displays(mainThreadLoop)
		require.NoError(t, err)
		assert.NotNil(t, displays)
	})
}

func TestDisplaysAPI_All(t *testing.T) {
	t.Run("should return all displays", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		// when
		all := displays.All()
		// then
		assert.NotEmpty(t, all)
	})
}

func TestDisplaysAPI_Primary(t *testing.T) {
	t.Run("should return primary display", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		primary, ok := displays.Primary()
		require.True(t, ok)
		assert.NotNil(t, primary)
	})
}

func TestDisplay_Name(t *testing.T) {
	t.Run("should return display's name", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		// when
		name := display.Name()
		// then
		assert.NotEmpty(t, name)
	})
}

func TestDisplay_Workarea(t *testing.T) {
	t.Run("should return display's workarea", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		// when
		workarea := display.Workarea()
		assert.True(t, workarea.X() >= 0)
		assert.True(t, workarea.Y() >= 0)
		assert.True(t, workarea.Width() > 0)
		assert.True(t, workarea.Height() > 0)
	})
}

func TestDisplay_VideoMode(t *testing.T) {
	t.Run("should return current video mode for display", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		// when
		mode := display.VideoMode()
		// then
		assert.True(t, mode.Width() > 0)
		assert.True(t, mode.Height() > 0)
		assert.True(t, mode.RefreshRate() >= 0)
	})
}

func TestDisplay_PhysicalSize(t *testing.T) {
	t.Run("should return physical size for display", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		// when
		size := display.PhysicalSize()
		// then
		assert.True(t, size.Width() > 0)
		assert.True(t, size.Height() > 0)
	})
}

func TestDisplay_VideoModes(t *testing.T) {
	t.Run("should return all vide modes for display", func(t *testing.T) {
		displays, _ := glfw.Displays(mainThreadLoop)
		display, _ := displays.Primary()
		// when
		modes := display.VideoModes()
		// then
		assert.True(t, len(modes) > 0)
		// and
		for _, mode := range modes {
			assert.True(t, mode.Width() > 0)
			assert.True(t, mode.Height() > 0)
			assert.True(t, mode.RefreshRate() >= 0)
		}
	})
}
