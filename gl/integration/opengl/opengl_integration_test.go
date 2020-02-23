package opengl_test

import (
	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/opengl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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

func TestContext_NewFloatVertexBuffer(t *testing.T) {
	t.Run("two buffers should have different IDs", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		// when
		buffer1 := context.NewFloatVertexBuffer(1)
		buffer2 := context.NewFloatVertexBuffer(1)
		// then
		assert.NotEqual(t, buffer1.ID(), buffer2.ID())
	})
}

func TestFloatVertexBuffer_Upload(t *testing.T) {
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
				context := gl.ContextOf(openGL)
				buffer := context.NewFloatVertexBuffer(test.size)
				defer buffer.Delete()
				// when
				buffer.Upload(test.offset, test.input)
				// then
				output := make([]float32, len(test.expected))
				buffer.Download(test.offset, output)
				assert.InDeltaSlice(t, test.expected, output, 1e-35)
			})
		}
	})
}

func TestFloatVertexBuffer_Download(t *testing.T) {
	openGL, _ := opengl.New(mainThreadLoop)
	defer openGL.Destroy()
	t.Run("should download data", func(t *testing.T) {
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
				context := gl.ContextOf(openGL)
				buffer := context.NewFloatVertexBuffer(len(test.input))
				defer buffer.Delete()
				buffer.Upload(0, test.input)
				// when
				buffer.Download(test.offset, test.output)
				// then
				assert.InDeltaSlice(t, test.expectedOutput, test.output, 1e-35)
			})
		}
	})
}

func TestContext_NewVertexArray(t *testing.T) {
	t.Run("should create vertex array", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		// when
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		// then
		assert.NotNil(t, vao)
		// cleanup
		vao.Delete()
	})
}
func TestVertexArray_Set(t *testing.T) {
	t.Run("should set", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		defer vao.Delete()
		buffer := context.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		pointer := gl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: 1,
		}
		// when
		vao.Set(0, pointer)
	})
}

func TestContext_CompileFragmentShader(t *testing.T) {
	t.Run("should return error for incorrect shader", func(t *testing.T) {
		tests := map[string]string{
			"golang code": "package main\nfunc main() {}",
		}
		for name, source := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				// when
				shader, err := context.CompileFragmentShader(source)
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
				context := gl.ContextOf(openGL)
				// when
				shader, err := context.CompileFragmentShader(source)
				require.NoError(t, err)
				assert.NotNil(t, shader)
			})
		}
	})
	t.Run("should not panic for empty shader", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		// when
		_, _ = context.CompileFragmentShader("")
	})
}

func TestContext_CompileVertexShader(t *testing.T) {
	t.Run("should return error for incorrect shader", func(t *testing.T) {
		tests := map[string]string{
			"golang code": "package main\nfunc main() {}",
		}
		for name, source := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				// when
				shader, err := context.CompileVertexShader(source)
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
				context := gl.ContextOf(openGL)
				// when
				shader, err := context.CompileVertexShader(source)
				// then
				require.NoError(t, err)
				assert.NotNil(t, shader)
			})
		}
	})
	t.Run("should not panic for empty shader", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		// when
		_, _ = context.CompileVertexShader("")
	})
}

func TestOpenGL_LinkProgram(t *testing.T) {
	t.Run("should return error", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		vertexShader, err := context.CompileVertexShader(`
								#version 330 core
								void noMain() {}
								`)
		require.NoError(t, err)
		fragmentShader, err := context.CompileFragmentShader(`
								#version 330 core
								void noMainEither() {}
								`)
		require.NoError(t, err)
		// when
		program, err := context.LinkProgram(vertexShader, fragmentShader)
		// then
		assert.Error(t, err)
		assert.Nil(t, program)
	})
	t.Run("should return program", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		vertexShader, _ := context.CompileVertexShader(`
								#version 330 core
								void main() {
									gl_Position = vec4(0, 0, 0, 0);
								}
								`)
		fragmentShader, _ := context.CompileFragmentShader(`
								#version 330 core
								void main() {}
								`)
		// when
		program, err := context.LinkProgram(vertexShader, fragmentShader)
		// then
		require.NoError(t, err)
		assert.NotNil(t, program)
	})
}

func TestOpenGL_Capabilities(t *testing.T) {
	t.Run("should return capabilities", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		// when
		capabilities := context.Capabilities()
		// then
		assert.NotNil(t, capabilities)
		assert.Greater(t, capabilities.MaxTextureSize(), 0)
	})
}

func TestOpenGL_NewAcceleratedImage(t *testing.T) {
	t.Run("should create AcceleratedImage", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		// when
		img := context.NewAcceleratedImage(0, 0)
		// then
		assert.NotNil(t, img)
	})
}

func TestAcceleratedCommand_Run(t *testing.T) {
	t.Run("vertex buffer can be used inside command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		program := workingProgram(t, context)
		output := context.NewAcceleratedImage(1, 1)
		command := program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
			buffer := context.NewFloatVertexBuffer(1)
			values := []float32{1}
			buffer.Upload(0, values)
			buffer.Download(0, values)
			buffer.Delete()
		}})
		// when
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("vertex array can be used inside command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		program := workingProgram(t, context)
		output := context.NewAcceleratedImage(1, 1)
		command := program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
			array := context.NewVertexArray(gl.VertexLayout{gl.Float})
			defer array.Delete()
			buffer := context.NewFloatVertexBuffer(1)
			defer buffer.Delete()
			array.Set(0, gl.VertexBufferPointer{
				Buffer: buffer,
				Offset: 0,
				Stride: 1,
			})
		}})
		// when
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("clear image fragment with color", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
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
			"height out of bound": {
				width: 1, height: 1,
				location:       image.AcceleratedImageLocation{Width: 1, Height: 2},
				expectedColors: []image.Color{color},
			},
			"top left corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, image.Transparent, color, image.Transparent},
			},
			"top row": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Width: 2, Height: 1},
				expectedColors: []image.Color{image.Transparent, image.Transparent, color, color},
			},
			"left column": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Width: 1, Height: 2},
				expectedColors: []image.Color{color, image.Transparent, color, image.Transparent},
			},
			"top right corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{X: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, image.Transparent, image.Transparent, color},
			},
			"bottom left corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{color, image.Transparent, image.Transparent, image.Transparent},
			},
			"bottom right corner": {
				width: 2, height: 2,
				location:       image.AcceleratedImageLocation{X: 1, Y: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, color, image.Transparent, image.Transparent},
			},
			"middle row": {
				width: 1, height: 3,
				location:       image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 1},
				expectedColors: []image.Color{image.Transparent, color, image.Transparent},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := context.NewAcceleratedImage(test.width, test.height)
				img.Upload(make([]image.Color, test.width*test.height))
				program := workingProgram(t, context)
				glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
					renderer.Clear(color)
				}}
				command := program.AcceleratedCommand(glCommand)
				// when
				command.Run(image.AcceleratedImageSelection{
					Location: test.location,
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				assertColors(t, test.expectedColors, img)
			})
		}
	})
	t.Run("should not change the image pixels when command does not do anything", func(t *testing.T) {
		commands := map[string]gl.Command{
			"nil":   nil,
			"empty": &emptyCommand{},
		}
		for name, command := range commands {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				img := context.NewAcceleratedImage(2, 1)
				pixels := []image.Color{image.RGB(1, 2, 3), image.RGB(4, 5, 6)}
				img.Upload(pixels)
				program := workingProgram(t, context)
				command := program.AcceleratedCommand(command)
				// when
				command.Run(image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{
						X:      0,
						Y:      0,
						Width:  1,
						Height: 1,
					},
					Image: img,
				}, []image.AcceleratedImageSelection{})
				// then
				assertColors(t, pixels, img)
			})
		}
	})
}

func workingProgram(t *testing.T, gl *gl.Context) *gl.Program {
	return compileProgram(t, gl,
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

func compileProgram(t *testing.T, context *gl.Context,
	vertexShaderSrc, fragmentShaderSrc string) *gl.Program {
	vertexShader, err := context.CompileVertexShader(vertexShaderSrc)
	require.NoError(t, err)
	fragmentShader, err := context.CompileFragmentShader(fragmentShaderSrc)
	require.NoError(t, err)
	program, err := context.LinkProgram(vertexShader, fragmentShader)
	require.NoError(t, err)
	return program
}

type command struct {
	runGL func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection)
}

func (c *command) RunGL(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
	c.runGL(renderer, selections)
}

type emptyCommand struct {
}

func (e emptyCommand) RunGL(_ *gl.Renderer, _ []image.AcceleratedImageSelection) {}

func assertColors(t *testing.T, expected []image.Color, img *gl.AcceleratedImage) {
	output := make([]image.Color, len(expected))
	img.Download(output)
	assert.Equal(t, expected, output)
}
