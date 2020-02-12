package opengl_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			_, _ = opengl.New(nil)
		})
	})
	t.Run("should create OpenGL using supplied MainThreadLoop", func(t *testing.T) {
		// when
		openGL, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// then
		assert.NotNil(t, openGL)
	})
	t.Run("should create 2 objects working at the same time", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			openGL, err := opengl.New(mainThreadLoop)
			require.NoError(t, err)
			defer openGL.Destroy()
		}
	})
}

func TestOpenGL_NewImage(t *testing.T) {
	t.Run("should return error for negative width", func(t *testing.T) {
		openGL, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		img, err := openGL.NewImage(-1, 0)
		// then
		require.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should return error for negative height", func(t *testing.T) {
		openGL, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		img, err := openGL.NewImage(0, -1)
		// then
		require.Error(t, err)
		assert.Nil(t, img)
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
				openGL, err := opengl.New(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				// when
				img, err := openGL.NewImage(test.width, test.height)
				// then
				require.NoError(t, err)
				assert.NotNil(t, img)
				assert.Equal(t, test.width, img.Width())
				assert.Equal(t, test.height, img.Height())
			})
		}
	})
}

func TestOpenGL_NewTexture(t *testing.T) {
	t.Run("should return error for negative width", func(t *testing.T) {
		openGL, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		img, err := openGL.NewAcceleratedImage(-1, 0)
		// then
		require.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should return error for negative height", func(t *testing.T) {
		openGL, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		img, err := openGL.NewAcceleratedImage(0, -1)
		// then
		require.Error(t, err)
		assert.Nil(t, img)
	})
	t.Run("should create AcceleratedImage", func(t *testing.T) {
		openGL, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer openGL.Destroy()
		// when
		img, err := openGL.NewAcceleratedImage(0, 0)
		// then
		require.NoError(t, err)
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
			"0x0": {
				width:       0,
				height:      0,
				inputColors: []image.Color{},
			},
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
				openGL, err := opengl.New(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				img, err := openGL.NewAcceleratedImage(test.width, test.height)
				require.NoError(t, err)
				// when
				img.Upload(test.inputColors)
				// then
				assertColors(t, test.inputColors, img)
			})
		}
	})
	t.Run("2 OpenGL contexts", func(t *testing.T) {
		gl1, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer gl1.Destroy()
		gl2, err := opengl.New(mainThreadLoop)
		require.NoError(t, err)
		defer gl2.Destroy()
		img1, err := gl1.NewAcceleratedImage(1, 1)
		require.NoError(t, err)
		img2, err := gl2.NewAcceleratedImage(1, 1)
		require.NoError(t, err)
		// when
		img1.Upload([]image.Color{color1})
		img2.Upload([]image.Color{color2})
		// then
		assertColors(t, []image.Color{color1}, img1)
		assertColors(t, []image.Color{color2}, img2)
	})
}

func TestRunOrDie(t *testing.T) {
	t.Run("should run provided callback", func(t *testing.T) {
		var callbackExecuted bool
		mainThreadLoop.Execute(func() {
			opengl.RunOrDie(func(gl *opengl.OpenGL) {
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
			opengl.RunOrDie(func(gl *opengl.OpenGL) {
				actualGL = gl
			})
		})
		assert.NotNil(t, actualGL)
	})
}

func TestOpenGL_OpenWindow(t *testing.T) {
	t.Run("should constrain width to platform-specific minimum if negative", func(t *testing.T) {
		openGL, err := opengl.New(mainThreadLoop)
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
		openGL, err := opengl.New(mainThreadLoop)
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
		openGL, err := opengl.New(mainThreadLoop)
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
		openGL, err := opengl.New(mainThreadLoop)
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
		openGL, err := opengl.New(mainThreadLoop)
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
		openGL, err := opengl.New(mainThreadLoop)
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
				openGL, err := opengl.New(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				// when
				win, err := openGL.OpenWindow(640, 360, opengl.Zoom(test.zoom))
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
				openGL, err := opengl.New(mainThreadLoop)
				require.NoError(t, err)
				defer openGL.Destroy()
				// when
				win, err := openGL.OpenWindow(640, 360, opengl.Zoom(test.zoom))
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

func TestOpenGL_CompileVertexShader(t *testing.T) {
	t.Run("should return error for incorrect shader", func(t *testing.T) {
		tests := map[string]string{
			"golang code": "package main\nfunc main() {}",
		}
		for name, source := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				shader, err := openGL.CompileVertexShader(source)
				// then
				assert.Error(t, err)
				assert.Nil(t, shader)
			})
		}
	})
	t.Run("should compile shader", func(t *testing.T) {
		tests := map[string]string{
			"GLSL 1.10": "void main() {}",
			"minimal": `
				#version 330 core
				void main() {
					gl_Position = vec4(0, 0, 0, 0);
				}
				`,
		}
		for name, source := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				shader, err := openGL.CompileVertexShader(source)
				// then
				require.NoError(t, err)
				assert.NotNil(t, shader)
			})
		}
	})
	t.Run("should not panic for empty shader", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		// when
		_, _ = openGL.CompileVertexShader("")
	})
}

func TestOpenGL_CompileFragmentShader(t *testing.T) {
	t.Run("should return error for incorrect shader", func(t *testing.T) {
		tests := map[string]string{
			"golang code": "package main\nfunc main() {}",
		}
		for name, source := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				shader, err := openGL.CompileFragmentShader(source)
				assert.Error(t, err)
				assert.Nil(t, shader)
			})
		}
	})
	t.Run("should compile shader", func(t *testing.T) {
		tests := map[string]string{
			"GLSL 1.10": "void main() {}",
			"minimal": `
				#version 330 core
				void main() {}
				`,
		}
		for name, source := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				shader, err := openGL.CompileFragmentShader(source)
				require.NoError(t, err)
				assert.NotNil(t, shader)
			})
		}
	})
	t.Run("should not panic for empty shader", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		// when
		_, _ = openGL.CompileFragmentShader("")
	})
}

func TestOpenGL_LinkProgram(t *testing.T) {
	t.Run("should return error when vertex shader is nil", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		fragmentShader, _ := openGL.CompileFragmentShader("")
		// when
		program, err := openGL.LinkProgram(nil, fragmentShader)
		// then
		assert.Error(t, err)
		assert.Nil(t, program)
	})
	t.Run("should return error when fragment shader is nil", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vertexShader, _ := openGL.CompileVertexShader("")
		// when
		program, err := openGL.LinkProgram(vertexShader, nil)
		// then
		assert.Error(t, err)
		assert.Nil(t, program)
	})
	t.Run("should return error", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vertexShader, err := openGL.CompileVertexShader(`
								#version 330 core
								void noMain() {}
								`)
		require.NoError(t, err)
		fragmentShader, err := openGL.CompileFragmentShader(`
								#version 330 core
								void noMainEither() {}
								`)
		require.NoError(t, err)
		// when
		program, err := openGL.LinkProgram(vertexShader, fragmentShader)
		// then
		assert.Error(t, err)
		assert.Nil(t, program)
	})
	t.Run("should return program", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vertexShader, _ := openGL.CompileVertexShader(`
								#version 330 core
								void main() {
									gl_Position = vec4(0, 0, 0, 0);
								}
								`)
		fragmentShader, _ := openGL.CompileFragmentShader(`
								#version 330 core
								void main() {}
								`)
		// when
		program, err := openGL.LinkProgram(vertexShader, fragmentShader)
		// then
		require.NoError(t, err)
		assert.NotNil(t, program)
	})
}

func TestProgram_AcceleratedCommand(t *testing.T) {
	t.Run("should return error when passed command is nil", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		// when
		cmd, err := program.AcceleratedCommand(nil)
		assert.Error(t, err)
		assert.Nil(t, cmd)
	})
	t.Run("should return command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		// when
		cmd, err := program.AcceleratedCommand(&commandMock{})
		require.NoError(t, err)
		assert.NotNil(t, cmd)
	})
}

func TestAcceleratedCommand_Run(t *testing.T) {
	t.Run("should execute command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		texture, _ := openGL.NewAcceleratedImage(1, 1)
		output := image.AcceleratedImageSelection{
			Image: texture,
		}
		tests := map[string]struct {
			selections []image.AcceleratedImageSelection
		}{
			"empty selections": {
				selections: []image.AcceleratedImageSelection{},
			},
			"one selection": {
				selections: []image.AcceleratedImageSelection{{}},
			},
			"two selections": {
				selections: []image.AcceleratedImageSelection{{}, {}},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				command := &commandMock{}
				acceleratedCommand, _ := program.AcceleratedCommand(command)
				// when
				err := acceleratedCommand.Run(output, test.selections)
				// then
				require.NoError(t, err)
				assert.Equal(t, 1, command.executionCount)
				assert.Equal(t, test.selections, command.selections)
				assert.NotNil(t, command.renderer)
			})
		}
	})
	t.Run("should return error when command returned error", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img, _ := openGL.NewAcceleratedImage(1, 1)
		program := workingProgram(t, openGL)
		command, _ := program.AcceleratedCommand(&failingCommand{})
		// when
		err := command.Run(image.AcceleratedImageSelection{Image: img}, []image.AcceleratedImageSelection{})
		assert.Error(t, err)
	})
	t.Run("should return error when output image is nil", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		command, _ := program.AcceleratedCommand(&emptyCommand{})
		// when
		err := command.Run(image.AcceleratedImageSelection{}, []image.AcceleratedImageSelection{})
		assert.Error(t, err)
	})
	t.Run("should return error when output image and program were created in different OpenGL contexts", func(t *testing.T) {
		imageContext, _ := opengl.New(mainThreadLoop)
		defer imageContext.Destroy()
		programContext, _ := opengl.New(mainThreadLoop)
		defer programContext.Destroy()
		img, _ := imageContext.NewAcceleratedImage(1, 1)
		program := workingProgram(t, programContext)
		command, _ := program.AcceleratedCommand(&emptyCommand{})
		// when
		err := command.Run(image.AcceleratedImageSelection{
			Image: img,
		}, []image.AcceleratedImageSelection{})
		assert.Error(t, err)
	})
	t.Run("vertex buffer can be used inside command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		output, _ := openGL.NewAcceleratedImage(1, 1)
		command, _ := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			buffer, err := openGL.NewFloatVertexBuffer(1)
			require.NoError(t, err)
			values := []float32{1}
			require.NoError(t, buffer.Upload(0, values))
			require.NoError(t, buffer.Download(0, values))
			buffer.Delete()
			return nil
		}})
		// when
		err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
		assert.NoError(t, err)
	})
	t.Run("vertex array can be used inside command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		output, _ := openGL.NewAcceleratedImage(1, 1)
		command, _ := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
			require.NoError(t, err)
			defer array.Delete()
			buffer, _ := openGL.NewFloatVertexBuffer(1)
			defer buffer.Delete()
			err = array.Set(0, opengl.VertexBufferPointer{
				Buffer: buffer,
				Offset: 0,
				Stride: 1,
			})
			require.NoError(t, err)
			return nil
		}})
		// when
		err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
		assert.NoError(t, err)
	})
	t.Run("clear image fragment with color", func(t *testing.T) {
		color := image.RGBA(1, 2, 3, 4)
		tests := map[string]struct {
			width, height  int
			location       image.AcceleratedImageLocation
			expectedColors []image.Color
		}{
			"empty location": {
				width: 1, height: 1,
				location:       image.AcceleratedImageLocation{},
				expectedColors: []image.Color{image.Transparent},
			},
			"x out of bounds": {
				width: 1, height: 1,
				location:       image.AcceleratedImageLocation{X: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent},
			},
			"y out of bounds": {
				width: 1, height: 1,
				location:       image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent},
			},
			"whole image": {
				width: 1, height: 1,
				location:       image.AcceleratedImageLocation{Width: 1, Height: 1},
				expectedColors: []image.Color{color},
			},
			"top left corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Width: 1, Height: 1},
				expectedColors: []image.Color{color, image.Transparent, image.Transparent, image.Transparent},
			},
			"top row": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Width: 2, Height: 1},
				expectedColors: []image.Color{color, color, image.Transparent, image.Transparent},
			},
			"left column": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Width: 1, Height: 2},
				expectedColors: []image.Color{color, image.Transparent, color, image.Transparent},
			},
			"top right corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{X: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, color, image.Transparent, image.Transparent},
			},
			"bottom left corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, image.Transparent, color, image.Transparent},
			},
			"bottom right corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{X: 1, Y: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, image.Transparent, image.Transparent, color},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img, _ := openGL.NewAcceleratedImage(test.width, test.height)
				img.Upload(make([]image.Color, test.width*test.height))
				program := workingProgram(t, openGL)
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
					renderer.Clear(color)
					return nil
				}}
				command, _ := program.AcceleratedCommand(glCommand)
				// when
				err := command.Run(image.AcceleratedImageSelection{
					Location: test.location,
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				require.NoError(t, err)
				assertColors(t, test.expectedColors, img)
			})
		}
	})
	t.Run("should not change the image pixels when command does not do anything", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img, _ := openGL.NewAcceleratedImage(2, 1)
		pixels := []image.Color{image.RGB(1, 2, 3), image.RGB(4, 5, 6)}
		img.Upload(pixels)
		program := workingProgram(t, openGL)
		command, _ := program.AcceleratedCommand(&emptyCommand{})
		// when
		err := command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{
				X:      0,
				Y:      0,
				Width:  1,
				Height: 1,
			},
			Image: img,
		}, []image.AcceleratedImageSelection{})
		// then
		require.NoError(t, err)
		assertColors(t, pixels, img)
	})
}

func TestRenderer_Clear(t *testing.T) {
	tests := map[string]struct {
		color image.Color
	}{
		"1": {
			color: image.RGBA(1, 2, 3, 4),
		},
		"2": {
			color: image.RGBA(5, 6, 7, 8),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			openGL, _ := opengl.New(mainThreadLoop)
			defer openGL.Destroy()
			img, _ := openGL.NewAcceleratedImage(1, 1)
			img.Upload(make([]image.Color, 1))
			program := workingProgram(t, openGL)
			glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
				// when
				renderer.Clear(test.color)
				return nil
			}}
			command, _ := program.AcceleratedCommand(glCommand)
			err := command.Run(image.AcceleratedImageSelection{
				Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
				Image:    img,
			}, []image.AcceleratedImageSelection{})
			// then
			require.NoError(t, err)
			assertColors(t, []image.Color{test.color}, img)
		})
	}
}

func TestRenderer_DrawArrays(t *testing.T) {
	t.Run("should draw point using one vertex attribute", func(t *testing.T) {
		tests := map[string]struct {
			vertexShaderSrc string
			typ             opengl.Type
			data            []float32
		}{
			"float": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in float vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition - 1, 0, 0, 1);
								}
								`,
				typ:  opengl.Float,
				data: []float32{1},
			},
			"float2": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec2 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, 0, 1);
								}
								`,
				typ:  opengl.Float2,
				data: []float32{1, 2},
			},
			"float3": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec3 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, vertexPosition.z-3, 1);
								}
								`,
				typ:  opengl.Float3,
				data: []float32{1, 2, 3},
			},
			"float4": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec4 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, vertexPosition.z-3, vertexPosition.w-3);
								}
								`,
				typ:  opengl.Float4,
				data: []float32{1, 2, 3, 4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img, _ := openGL.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 1))
				vertexShader, err := openGL.CompileVertexShader(test.vertexShaderSrc)
				require.NoError(t, err)
				fragmentShader, err := openGL.CompileFragmentShader(`
								#version 330 core
								out vec4 color;
								void main() {
									color = vec4(0.2, 0.4, 0.6, 0.8);
								}
								`)
				require.NoError(t, err)
				program, err := openGL.LinkProgram(vertexShader, fragmentShader)
				require.NoError(t, err)
				array, err := openGL.NewVertexArray(opengl.VertexLayout{test.typ})
				require.NoError(t, err)
				buffer, err := openGL.NewFloatVertexBuffer(len(test.data))
				require.NoError(t, err)
				err = buffer.Upload(0, test.data)
				require.NoError(t, err)
				vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: len(test.data)}
				err = array.Set(0, vertexPosition)
				require.NoError(t, err)
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
					// when
					renderer.DrawArrays(array, opengl.Points, 0, 1)
					return nil
				}}
				command, _ := program.AcceleratedCommand(glCommand)
				err = command.Run(image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				require.NoError(t, err)
				assertColors(t, []image.Color{image.RGBA(51, 102, 153, 204)}, img)
			})
		}
	})
	t.Run("should draw point using 2 vertex attributes", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img, _ := openGL.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		vertexShader, err := openGL.CompileVertexShader(`
								#version 330 core
								layout(location = 0) in float vertexPositionX;
								layout(location = 1) in vec3 vertexColor;
								out vec4 interpolatedColor;
								void main() {
									gl_Position = vec4(vertexPositionX, 0, 0, 1);
									interpolatedColor = vec4(vertexColor, 1.);
								}
								`)
		require.NoError(t, err)
		fragmentShader, err := openGL.CompileFragmentShader(`
								#version 330 core
								in vec4 interpolatedColor;
								out vec4 color;
								void main() {
									color = interpolatedColor;
								}
								`)
		require.NoError(t, err)
		program, err := openGL.LinkProgram(vertexShader, fragmentShader)
		require.NoError(t, err)
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float, opengl.Float3})
		require.NoError(t, err)
		buffer, err := openGL.NewFloatVertexBuffer(4)
		require.NoError(t, err)
		err = buffer.Upload(0, []float32{0, 0.2, 0.4, 0.6})
		require.NoError(t, err)
		vertexPositionX := opengl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 4}
		err = array.Set(0, vertexPositionX)
		vertexColor := opengl.VertexBufferPointer{Buffer: buffer, Offset: 1, Stride: 4}
		err = array.Set(1, vertexColor)
		require.NoError(t, err)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			// when
			renderer.DrawArrays(array, opengl.Points, 0, 1)
			return nil
		}}
		command, _ := program.AcceleratedCommand(glCommand)
		err = command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		require.NoError(t, err)
		assertColors(t, []image.Color{image.RGB(51, 102, 153)}, img)
	})
	t.Run("should draw two points", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img, _ := openGL.NewAcceleratedImage(2, 1)
		img.Upload(make([]image.Color, 2))
		vertexShader, err := openGL.CompileVertexShader(`
								#version 330 core
								layout(location = 0) in vec2 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition, 0, 1);
								}
								`)
		require.NoError(t, err)
		fragmentShader, err := openGL.CompileFragmentShader(`
								#version 330 core
								out vec4 color;
								void main() {
									color = vec4(1.0, 0.89, 0.8, 0.7);
								}
								`)
		require.NoError(t, err)
		program, err := openGL.LinkProgram(vertexShader, fragmentShader)
		require.NoError(t, err)
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float, opengl.Float3})
		require.NoError(t, err)
		buffer, err := openGL.NewFloatVertexBuffer(4)
		require.NoError(t, err)
		err = buffer.Upload(0, []float32{0, 0, 1, 0})
		require.NoError(t, err)
		vertexPositionX := opengl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 2}
		err = array.Set(0, vertexPositionX)
		require.NoError(t, err)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			// when
			renderer.DrawArrays(array, opengl.Points, 0, 2)
			return nil
		}}
		command, _ := program.AcceleratedCommand(glCommand)
		err = command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 2, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		require.NoError(t, err)
		assertColors(t, []image.Color{image.RGBA(255, 227, 204, 178), image.RGBA(255, 227, 204, 178)}, img)
	})
}

func TestRenderer_BindTexture(t *testing.T) {
	t.Run("can't bind texture without uniformName", func(t *testing.T) {
		names := []string{"", " ", "  ", "\n", "\t"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				output, _ := openGL.NewAcceleratedImage(1, 1)
				tex, _ := openGL.NewAcceleratedImage(1, 1)
				program := workingProgram(t, openGL)
				command, _ := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
					// when
					return renderer.BindTexture(0, name, tex)
				}})
				err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
				// then
				assert.Error(t, err)
			})
		}
	})
	t.Run("can't bind texture with uniformName not specified in program", func(t *testing.T) {
		names := []string{"foo", "bar"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				output, _ := openGL.NewAcceleratedImage(1, 1)
				tex, _ := openGL.NewAcceleratedImage(1, 1)
				program := workingProgram(t, openGL) // TODO It's better to have a shader here
				command, _ := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
					// when
					return renderer.BindTexture(0, name, tex)
				}})
				err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
				// then
				assert.Error(t, err)
			})
		}
	})
	t.Run("can't bind texture created in a different context", func(t *testing.T) {
		openGL1, _ := opengl.New(mainThreadLoop)
		defer openGL1.Destroy()
		openGL2, _ := opengl.New(mainThreadLoop)
		defer openGL2.Destroy()
		output, _ := openGL1.NewAcceleratedImage(1, 1)
		tex, _ := openGL2.NewAcceleratedImage(1, 1)
		program := workingProgram(t, openGL1)
		command, _ := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			// when
			return renderer.BindTexture(0, "tex", tex)
		}})
		err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
		// then
		assert.Error(t, err)
	})
	t.Run("can't bind texture with negative texture unit", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		output, _ := openGL.NewAcceleratedImage(1, 1)
		tex, _ := openGL.NewAcceleratedImage(1, 1)
		program := workingProgram(t, openGL)
		command, _ := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			// when
			return renderer.BindTexture(-1, "tex", tex)
		}})
		err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
		// then
		assert.Error(t, err)
	})
	t.Run("can bind texture", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		output, _ := openGL.NewAcceleratedImage(1, 1)
		tex, _ := openGL.NewAcceleratedImage(1, 1)
		program := workingProgram(t, openGL) // TODO It's better to have a shader here
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			// when
			return renderer.BindTexture(0, "tex", tex)
		}}
		command, _ := program.AcceleratedCommand(glCommand)
		err := command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
		// then
		assert.NoError(t, err)
	})
	t.Run("should draw point by sampling texture", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img, _ := openGL.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		tex, _ := openGL.NewAcceleratedImage(1, 1)
		tex.Upload([]image.Color{image.RGBA(1, 2, 3, 4)})
		program := compileProgram(t,
			openGL,
			`
				#version 330 core
				layout(location = 0) in vec2 xy;	
				void main() {
					gl_Position = vec4(xy, 0.0, 1.0);
				}
				`,
			`
				#version 330 core
				uniform sampler2D tex;
				out vec4 color;
				void main() {
					color = texture(tex, vec2(0.0, 0.0));
				}
				`)
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float2, opengl.Float2})
		require.NoError(t, err)
		buffer, err := openGL.NewFloatVertexBuffer(2)
		require.NoError(t, err)
		err = buffer.Upload(0, []float32{0.0, 0.0})
		require.NoError(t, err)
		vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: 2}
		err = array.Set(0, vertexPosition)
		require.NoError(t, err)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
			// when
			err := renderer.BindTexture(0, "tex", tex)
			if err != nil {
				return err
			}
			renderer.DrawArrays(array, opengl.Points, 0, 1)
			return nil
		}}
		command, _ := program.AcceleratedCommand(glCommand)
		err = command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		require.NoError(t, err)
		assertColors(t, []image.Color{image.RGBA(1, 2, 3, 4)}, img)
	})
	t.Run("should draw point by sampling 2 textures", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img, _ := openGL.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		tex1, _ := openGL.NewAcceleratedImage(1, 1)
		tex1.Upload([]image.Color{image.RGBA(5, 6, 7, 8)})
		tex2, _ := openGL.NewAcceleratedImage(1, 1)
		tex2.Upload([]image.Color{image.RGBA(9, 10, 11, 12)})
		program := compileProgram(t,
			openGL,
			`
				#version 330 core
				layout(location = 0) in vec2 xy;	
				void main() {
					gl_Position = vec4(xy, 0.0, 1.0);
				}
				`,
			`
				#version 330 core
				uniform sampler2D tex1;
				uniform sampler2D tex2;
				out vec4 color;
				void main() {
					color = texture(tex1, vec2(0.0, 0.0)) + texture(tex2, vec2(0.0, 0.0));
				}
				`)
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float2, opengl.Float2})
		require.NoError(t, err)
		buffer, err := openGL.NewFloatVertexBuffer(2)
		require.NoError(t, err)
		err = buffer.Upload(0, []float32{0.0, 0.0})
		require.NoError(t, err)
		vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: 2}
		err = array.Set(0, vertexPosition)
		require.NoError(t, err)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, _ []image.AcceleratedImageSelection) error {
			// when
			err := renderer.BindTexture(0, "tex1", tex1)
			if err != nil {
				return err
			}
			err = renderer.BindTexture(1, "tex2", tex2)
			if err != nil {
				return err
			}
			renderer.DrawArrays(array, opengl.Points, 0, 1)
			return nil
		}}
		command, _ := program.AcceleratedCommand(glCommand)
		err = command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		require.NoError(t, err)
		assertColors(t, []image.Color{image.RGBA(5+9, 6+10, 7+11, 8+12)}, img)
	})
}

func assertColors(t *testing.T, expected []image.Color, img *opengl.AcceleratedImage) {
	output := make([]image.Color, len(expected))
	img.Download(output)
	assert.Equal(t, expected, output)
}

func TestOpenGL_NewFloatVertexBuffer(t *testing.T) {
	t.Run("should return error when size is negative", func(t *testing.T) {
		tests := map[string]int{
			"size -1": -1,
			"size -2": -2,
		}
		for name, size := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				buffer, err := openGL.NewFloatVertexBuffer(size)
				// then
				assert.Error(t, err)
				assert.Nil(t, buffer)
			})
		}
	})
	t.Run("should create FloatVertexBuffer", func(t *testing.T) {
		tests := map[string]int{
			"size 0": 0,
			"size 1": 1,
		}
		for name, size := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				buffer, err := openGL.NewFloatVertexBuffer(size)
				// then
				require.NoError(t, err)
				assert.NotNil(t, buffer)
				// and
				assert.Equal(t, size, buffer.Size())
			})
		}
	})
}

func TestFloatVertexBuffer_Upload(t *testing.T) {
	t.Run("should return error when trying to upload slice bigger than size", func(t *testing.T) {
		tests := map[string]struct {
			offset int
			size   int
			data   []float32
		}{
			"size 0, offset 0, data len 1": {
				data: []float32{1},
			},
			"size 1, offset 0, data len 2": {
				size: 1,
				data: []float32{1, 2},
			},
			"size 0, offset 1, data len 1": {
				offset: 1,
				data:   []float32{1},
			},
			"size 1, offset 1, data len 1": {
				size:   1,
				offset: 1,
				data:   []float32{1},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				buffer, _ := openGL.NewFloatVertexBuffer(test.size)
				defer buffer.Delete()
				// when
				err := buffer.Upload(test.offset, test.data)
				assert.Error(t, err)
			})
		}
	})
	t.Run("should return error when offset is negative", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		// when
		err := buffer.Upload(-1, []float32{1})
		assert.Error(t, err)
	})
	t.Run("should upload data", func(t *testing.T) {
		tests := map[string]struct {
			size     int
			offset   int
			input    []float32
			expected []float32
		}{
			"offset 0": {
				size:     1,
				offset:   0,
				input:    []float32{1},
				expected: []float32{1},
			},
			"offset 0, size 2": {
				size:     2,
				offset:   0,
				input:    []float32{1, 2},
				expected: []float32{1, 2},
			},
			"offset 1": {
				size:     2,
				offset:   1,
				input:    []float32{1},
				expected: []float32{1},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				buffer, _ := openGL.NewFloatVertexBuffer(test.size)
				defer buffer.Delete()
				// when
				err := buffer.Upload(test.offset, test.input)
				// then
				require.NoError(t, err)
				// and
				output := make([]float32, len(test.expected))
				err = buffer.Download(test.offset, output)
				require.NoError(t, err)
				assert.InDeltaSlice(t, test.expected, output, 1e-35)
			})
		}
	})

}

func TestFloatVertexBuffer_Download(t *testing.T) {
	t.Run("should return error when offset is negative", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		output := make([]float32, 1)
		// when
		err := buffer.Download(-1, output)
		assert.Error(t, err)
	})
	t.Run("should download data", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		tests := map[string]struct {
			input          []float32
			offset         int
			output         []float32
			expectedOutput []float32
		}{
			"empty output slice": {
				input:          []float32{1},
				output:         make([]float32, 0),
				expectedOutput: []float32{},
			},
			"nil output slice": {
				input:          []float32{1},
				output:         nil,
				expectedOutput: nil,
			},
			"one element slice": {
				input:          []float32{1},
				output:         make([]float32, 1),
				expectedOutput: []float32{1},
			},
			"two elements slice": {
				input:          []float32{1, 2},
				output:         make([]float32, 2),
				expectedOutput: []float32{1, 2},
			},
			"output slice bigger than buffer": {
				input:          []float32{1},
				output:         make([]float32, 2),
				expectedOutput: []float32{1, 0},
			},
			"offset: 1": {
				offset:         1,
				input:          []float32{1, 2},
				output:         make([]float32, 1),
				expectedOutput: []float32{2},
			},
			"output slice bigger than remaining buffer": {
				offset:         1,
				input:          []float32{1, 2},
				output:         make([]float32, 2),
				expectedOutput: []float32{2, 0},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				buffer, _ := openGL.NewFloatVertexBuffer(len(test.input))
				defer buffer.Delete()
				_ = buffer.Upload(0, test.input)
				// when
				err := buffer.Download(test.offset, test.output)
				// then
				require.NoError(t, err)
				assert.InDeltaSlice(t, test.expectedOutput, test.output, 1e-35)
			})
		}
	})
}

func TestOpenGL_NewVertexArray(t *testing.T) {
	t.Run("should return error", func(t *testing.T) {
		tests := map[string]struct {
			layout opengl.VertexLayout
		}{
			"nil layout": {
				layout: nil,
			},
			"empty layout": {
				layout: opengl.VertexLayout{},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				// when
				vao, err := openGL.NewVertexArray(test.layout)
				// then
				assert.Error(t, err)
				assert.Nil(t, vao)
			})
		}
	})
	t.Run("should create vertex array", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		// when
		vao, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		// then
		require.NoError(t, err)
		assert.NotNil(t, vao)
		// cleanup
		vao.Delete()
	})
}

func TestVertexArray_Set(t *testing.T) {
	t.Run("should return error when offset is negative", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		pointer := opengl.VertexBufferPointer{
			Buffer: buffer,
			Offset: -1,
			Stride: 1,
		}
		// when
		err := vao.Set(0, pointer)
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when stride is negative", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		pointer := opengl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: -1,
		}
		// when
		err := vao.Set(0, pointer)
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when location is negative", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		pointer := opengl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: 1,
		}
		// when
		err := vao.Set(-1, pointer)
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when location is higher than number of arguments", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		pointer := opengl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: 1,
		}
		// when
		err := vao.Set(1, pointer)
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when buffer is nil", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		pointer := opengl.VertexBufferPointer{
			Buffer: nil,
			Offset: 0,
			Stride: 1,
		}
		// when
		err := vao.Set(0, pointer)
		// then
		assert.Error(t, err)
	})
	t.Run("should return error when buffer was not created by context", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		vertexBufferNotCreatedInContext := &opengl.FloatVertexBuffer{}
		pointer := opengl.VertexBufferPointer{
			Buffer: vertexBufferNotCreatedInContext,
			Offset: 0,
			Stride: 1,
		}
		// when
		err := vao.Set(0, pointer)
		// then
		assert.Error(t, err)
	})
	t.Run("should set", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		vao, _ := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
		defer vao.Delete()
		buffer, _ := openGL.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		pointer := opengl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: 1,
		}
		// when
		err := vao.Set(0, pointer)
		// then
		assert.NoError(t, err)
	})
}

func workingProgram(t *testing.T, openGL *opengl.OpenGL) *opengl.Program {
	return compileProgram(t, openGL,
		`#version 330 core
						void main() {
							gl_Position = vec4(0, 0, 0, 0);
						}`,
		`#version 330 core
						 uniform sampler2D tex;
						 out vec4 color;
						 void main() {
						 	color = texture(tex, vec2(0,0));
						 }`)
}

func compileProgram(t *testing.T, openGL *opengl.OpenGL,
	vertexShaderSrc, fragmentShaderSrc string) *opengl.Program {
	vertexShader, err := openGL.CompileVertexShader(vertexShaderSrc)
	require.NoError(t, err)
	fragmentShader, err := openGL.CompileFragmentShader(fragmentShaderSrc)
	require.NoError(t, err)
	program, err := openGL.LinkProgram(vertexShader, fragmentShader)
	require.NoError(t, err)
	return program
}

type commandMock struct {
	executionCount int
	selections     []image.AcceleratedImageSelection
	renderer       *opengl.Renderer
}

func (f *commandMock) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
	f.executionCount++
	f.selections = selections
	f.renderer = renderer
	return nil
}

type failingCommand struct{}

func (f *failingCommand) RunGL(*opengl.Renderer, []image.AcceleratedImageSelection) error {
	return errors.New("command failed")
}

type command struct {
	runGL func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error
}

func (c *command) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
	return c.runGL(renderer, selections)
}

type emptyCommand struct {
}

func (e emptyCommand) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
	return nil
}
