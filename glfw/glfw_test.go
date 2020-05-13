package glfw_test

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glfw"
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

func TestNew(t *testing.T) {
	t.Run("should panic when MainThreadLoop is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = glfw.NewOpenGL(nil)
		})
	})
	t.Run("should create OpenGL using supplied MainThreadLoop", func(t *testing.T) {
		// when
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// then
		assert.NotNil(t, openGL)
	})
	t.Run("should create 2 objects working at the same time", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			openGL, err := glfw.NewOpenGL(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
		}
	})
}

func TestOpenGL_ContextAPI(t *testing.T) {
	t.Run("should return context API", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		// when
		api := openGL.ContextAPI()
		// then
		assert.NotNil(t, api)
	})
}

func TestOpenGL_Context(t *testing.T) {
	t.Run("should return context", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		// when
		context := openGL.Context()
		// then
		assert.NotNil(t, context)
	})
	t.Run("on each invocation same context should be returned", func(t *testing.T) {
		openGL, _ := glfw.NewOpenGL(mainThreadLoop)
		defer openGL.Destroy()
		// when
		context1 := openGL.Context()
		context2 := openGL.Context()
		// then
		assert.Same(t, context1, context2)
	})
}

func TestOpenGL_NewImage(t *testing.T) {
	t.Run("should panic for negative width", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		assert.Panics(t, func() {
			// when
			openGL.NewImage(-1, 0)
		})
	})
	t.Run("should panic for negative height", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		assert.Panics(t, func() {
			// when
			openGL.NewImage(0, -1)
		})
	})
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
				openGL, err := glfw.NewOpenGL(mainThreadLoop)
				require.NoError(t, err)
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

func TestRunOrDie(t *testing.T) {
	t.Run("should run provided callback", func(t *testing.T) {
		var callbackExecuted bool
		mainThreadLoop.Execute(func() {
			glfw.RunOrDie(func(_ *glfw.OpenGL) {
				callbackExecuted = true
			})
		})
		assert.True(t, callbackExecuted)
	})
	t.Run("should start a MainThreadLoop and create OpenGL object", func(t *testing.T) {
		var (
			actualGL *glfw.OpenGL
		)
		mainThreadLoop.Execute(func() {
			glfw.RunOrDie(func(openGL *glfw.OpenGL) {
				actualGL = openGL
			})
		})
		assert.NotNil(t, actualGL)
	})
}

func TestOpenGL_OpenWindow(t *testing.T) {
	t.Run("should constrain width to platform-specific minimum if negative", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win, err := openGL.OpenWindow(-1, 0)
		require.NoError(t, err)
		defer win.Close()
		// then
		require.NotNil(t, win)
		assert.GreaterOrEqual(t, win.Width(), 0)
	})
	t.Run("should constrain height to platform-specific minimum if negative", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win, err := openGL.OpenWindow(0, -1)
		require.NoError(t, err)
		defer win.Close()
		// then
		require.NotNil(t, win)
		assert.GreaterOrEqual(t, win.Height(), 0)
	})
	t.Run("should open Window", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win, err := openGL.OpenWindow(640, 360)
		require.NoError(t, err)
		defer win.Close()
		// then
		require.NotNil(t, win)
		assert.Equal(t, 640, win.Width())
		assert.Equal(t, 360, win.Height())
	})
	t.Run("should open two windows at the same time", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win1, err := openGL.OpenWindow(640, 360)
		require.NoError(t, err)
		defer win1.Close()
		win2, err := openGL.OpenWindow(320, 180)
		require.NoError(t, err)
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
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		win1, err := openGL.OpenWindow(640, 360)
		require.NoError(t, err)
		win1.Close()
		// when
		win2, err := openGL.OpenWindow(320, 180)
		require.NoError(t, err)
		defer win2.Close()
		// then
		require.NotNil(t, win2)
		assert.Equal(t, 320, win2.Width())
		assert.Equal(t, 180, win2.Height())
	})
	t.Run("should skip nil option", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win, err := openGL.OpenWindow(0, 0, nil)
		require.NoError(t, err)
		defer win.Close()
	})
	t.Run("zoom <= 1 should not affect the width and height", func(t *testing.T) {
		tests := map[string]struct {
			zoom int
		}{
			"zoom = -1": {
				zoom: -1,
			},
			"zoom = 0": {
				zoom: 0,
			},
			"zoom = 1": {
				zoom: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, err := glfw.NewOpenGL(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				// when
				win, err := openGL.OpenWindow(640, 360, glfw.Zoom(test.zoom))
				require.NoError(t, err)
				defer win.Close()
				// then
				require.NotNil(t, win)
				assert.Equal(t, 640, win.Width())
				assert.Equal(t, 360, win.Height())
			})
		}
	})
	t.Run("zoom should affect the width and height", func(t *testing.T) {
		tests := map[string]struct {
			zoom           int
			expectedWidth  int
			expectedHeight int
		}{
			"zoom = 2": {
				zoom:           2,
				expectedWidth:  1280,
				expectedHeight: 720,
			},
			"zoom = 3": {
				zoom:           3,
				expectedWidth:  1920,
				expectedHeight: 1080,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, err := glfw.NewOpenGL(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				// when
				win, err := openGL.OpenWindow(640, 360, glfw.Zoom(test.zoom))
				require.NoError(t, err)
				defer win.Close()
				// then
				require.NotNil(t, win)
				assert.Equal(t, test.expectedWidth, win.Width())
				assert.Equal(t, test.expectedHeight, win.Height())
			})
		}
	})
}

func TestWindow_Width(t *testing.T) {
	t.Run("concurrent Width() calls should return the same value", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win, err := openGL.OpenWindow(640, 360)
		require.NoError(t, err)
		defer win.Close()
		// then
		var wg sync.WaitGroup
		goroutines := 100
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				assert.Equal(t, win.Width(), 640)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func TestWindow_Height(t *testing.T) {
	t.Run("concurrent Height() calls should return the same value", func(t *testing.T) {
		openGL, err := glfw.NewOpenGL(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		win, err := openGL.OpenWindow(640, 360)
		require.NoError(t, err)
		defer win.Close()
		// then
		var wg sync.WaitGroup
		goroutines := 100
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				assert.Equal(t, win.Height(), 360)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func TestOpenGL_NewCursor(t *testing.T) {
	openGL, err := glfw.NewOpenGL(mainThreadLoop)
	require.NoError(t, err)
	defer openGL.Destroy()
	img := openGL.NewImage(16, 16)
	selection := img.WholeImageSelection()
	t.Run("should create cursor with no options", func(t *testing.T) {
		// when
		cursor := openGL.NewCursor(selection)
		// then
		require.NotNil(t, cursor)
		cursor.Destroy()
	})
	t.Run("should create cursor with Hotspot option", func(t *testing.T) {
		tests := map[string]struct {
			x, y int
		}{
			"0,0": {},
			"1,2": {
				x: 1,
				y: 2,
			},
			"2,1": {
				x: 2,
				y: 1,
			},
			"1, selection height": {
				x: 1,
				y: selection.Height(),
			},
			"1, selection height + 1": {
				x: 1,
				y: selection.Height() + 1,
			},
			"selection width, 1": {
				x: selection.Width(),
				y: 1,
			},
			"selection width + 1, 1": {
				x: selection.Width() + 1,
				y: 1,
			},
			"-1,0": {
				x: -1,
			},
			"0,-1": {
				y: -1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// when
				cursor := openGL.NewCursor(selection, glfw.Hotspot(test.x, test.y))
				// then
				require.NotNil(t, cursor)
				cursor.Destroy()
			})
		}
	})
	t.Run("should create cursor with CursorZoom option", func(t *testing.T) {
		zooms := []int{0, 1, 2, 100}
		for _, zoom := range zooms {
			// when
			cursor := openGL.NewCursor(selection, glfw.CursorZoom(zoom))
			// then
			require.NotNil(t, cursor)
			cursor.Destroy()
		}
	})
}
