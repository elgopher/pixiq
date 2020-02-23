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

func TestContext_LinkProgram(t *testing.T) {
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

func TestContext_Capabilities(t *testing.T) {
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

func TestContext_NewAcceleratedImage(t *testing.T) {
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

func TestAcceleratedImage_Upload(t *testing.T) {
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
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				img := context.NewAcceleratedImage(test.width, test.height)
				// when
				img.Upload(test.inputColors)
				// then
				assertColors(t, test.inputColors, img)
			})
		}
	})
	t.Run("2 OpenGL contexts", func(t *testing.T) {
		gl1, _ := opengl.New(mainThreadLoop)
		defer gl1.Destroy()
		context1 := gl.ContextOf(gl1)
		gl2, _ := opengl.New(mainThreadLoop)
		defer gl2.Destroy()
		context2 := gl.ContextOf(gl2)

		img1 := context1.NewAcceleratedImage(1, 1)
		img2 := context2.NewAcceleratedImage(1, 1)
		// when
		img1.Upload([]image.Color{color1})
		img2.Upload([]image.Color{color2})
		// then
		assertColors(t, []image.Color{color1}, img1)
		assertColors(t, []image.Color{color2}, img2)
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
			context := gl.ContextOf(openGL)
			img := context.NewAcceleratedImage(1, 1)
			img.Upload(make([]image.Color, 1))
			program := workingProgram(t, context)
			glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
				// when
				renderer.Clear(test.color)
			}}
			command := program.AcceleratedCommand(glCommand)
			command.Run(image.AcceleratedImageSelection{
				Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
				Image:    img,
			}, []image.AcceleratedImageSelection{})
			// then
			assertColors(t, []image.Color{test.color}, img)
		})
	}
}

func TestRenderer_DrawArrays(t *testing.T) {
	t.Run("should draw point using one vertex attribute", func(t *testing.T) {
		tests := map[string]struct {
			vertexShaderSrc string
			typ             gl.Type
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
				typ:  gl.Float,
				data: []float32{1},
			},
			"vec2": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec2 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, 0, 1);
								}
								`,
				typ:  gl.Vec2,
				data: []float32{1, 2},
			},
			"vec3": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec3 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, vertexPosition.z-3, 1);
								}
								`,
				typ:  gl.Vec3,
				data: []float32{1, 2, 3},
			},
			"vec4": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec4 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, vertexPosition.z-3, vertexPosition.w-3);
								}
								`,
				typ:  gl.Vec4,
				data: []float32{1, 2, 3, 4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				img := context.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 1))
				vertexShader, err := context.CompileVertexShader(test.vertexShaderSrc)
				require.NoError(t, err)
				fragmentShader, err := context.CompileFragmentShader(`
								#version 330 core
								out vec4 color;
								void main() {
									color = vec4(0.2, 0.4, 0.6, 0.8);
								}
								`)
				require.NoError(t, err)
				program, err := context.LinkProgram(vertexShader, fragmentShader)
				require.NoError(t, err)
				array := context.NewVertexArray(gl.VertexLayout{test.typ})
				buffer := context.NewFloatVertexBuffer(len(test.data))
				buffer.Upload(0, test.data)
				vertexPosition := gl.VertexBufferPointer{Buffer: buffer, Stride: len(test.data)}
				array.Set(0, vertexPosition)
				glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					renderer.DrawArrays(array, gl.Points, 0, 1)
				}}
				command := program.AcceleratedCommand(glCommand)
				command.Run(image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				assertColors(t, []image.Color{image.RGBA(51, 102, 153, 204)}, img)
			})
		}
	})
	t.Run("should draw point using 2 vertex attributes", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		img := context.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		vertexShader, err := context.CompileVertexShader(`
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
		fragmentShader, err := context.CompileFragmentShader(`
								#version 330 core
								in vec4 interpolatedColor;
								out vec4 color;
								void main() {
									color = interpolatedColor;
								}
								`)
		require.NoError(t, err)
		program, err := context.LinkProgram(vertexShader, fragmentShader)
		require.NoError(t, err)
		array := context.NewVertexArray(gl.VertexLayout{gl.Float, gl.Vec3})
		require.NoError(t, err)
		buffer := context.NewFloatVertexBuffer(4)
		buffer.Upload(0, []float32{0, 0.2, 0.4, 0.6})
		vertexPositionX := gl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 4}
		array.Set(0, vertexPositionX)
		vertexColor := gl.VertexBufferPointer{Buffer: buffer, Offset: 1, Stride: 4}
		array.Set(1, vertexColor)
		glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.DrawArrays(array, gl.Points, 0, 1)
		}}
		command := program.AcceleratedCommand(glCommand)
		command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		assertColors(t, []image.Color{image.RGB(51, 102, 153)}, img)
	})
	t.Run("should draw triangle fan with location specified", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		color := image.RGBA(51, 102, 153, 204)
		tests := map[string]struct {
			width, height  int
			outputLocation image.AcceleratedImageLocation
			expectedColors []image.Color
		}{
			"X:1": {
				outputLocation: image.AcceleratedImageLocation{X: 1, Width: 1, Height: 1},
				width:          2, height: 1,
				expectedColors: []image.Color{image.Transparent, color},
			},
			"Y:1": {
				outputLocation: image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 1},
				width:          1, height: 2,
				expectedColors: []image.Color{color, image.Transparent},
			},
			"Width:2": {
				outputLocation: image.AcceleratedImageLocation{Width: 2, Height: 1},
				width:          2, height: 1,
				expectedColors: []image.Color{color, color},
			},
			"Width:3": {
				outputLocation: image.AcceleratedImageLocation{Width: 3, Height: 1},
				width:          3, height: 1,
				expectedColors: []image.Color{color, color, color},
			},
			"Height:2": {
				outputLocation: image.AcceleratedImageLocation{Width: 1, Height: 2},
				width:          1, height: 2,
				expectedColors: []image.Color{color, color},
			},
			"Height:3": {
				outputLocation: image.AcceleratedImageLocation{Width: 1, Height: 3},
				width:          1, height: 3,
				expectedColors: []image.Color{color, color, color},
			},
			"Height:2 and Y:1": {
				outputLocation: image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 2},
				width:          1, height: 3,
				expectedColors: []image.Color{color, color, image.Transparent},
			},
			"Height out of bounds": {
				outputLocation: image.AcceleratedImageLocation{Width: 1, Height: 3},
				width:          1, height: 2,
				expectedColors: []image.Color{color, color},
			},
			"middle row": {
				outputLocation: image.AcceleratedImageLocation{Y: 1, Width: 1, Height: 1},
				width:          1, height: 3,
				expectedColors: []image.Color{image.Transparent, color, image.Transparent},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				img := context.NewAcceleratedImage(test.width, test.height)
				img.Upload(make([]image.Color, test.width*test.height))
				program := compileProgram(t, context,
					`
								#version 330 core
								layout(location = 0) in vec2 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition, 0, 1);
								}
								`,
					`
								#version 330 core
								out vec4 color;
								void main() {
									color = vec4(0.2, 0.4, 0.6, 0.8);
								}
								`,
				)
				array := context.NewVertexArray(gl.VertexLayout{gl.Vec2})
				buffer := context.NewFloatVertexBuffer(8)
				buffer.Upload(0, []float32{
					-1, 1, // top-left
					1, 1, // top-right
					1, -1, // bottom-right
					-1, -1}, // bottom-left
				)
				vertexPosition := gl.VertexBufferPointer{Buffer: buffer, Stride: 2}
				array.Set(0, vertexPosition)
				glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					renderer.DrawArrays(array, gl.TriangleFan, 0, 4)
				}}
				command := program.AcceleratedCommand(glCommand)
				command.Run(image.AcceleratedImageSelection{
					Location: test.outputLocation,
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				assertColors(t, test.expectedColors, img)
			})
		}
	})
	t.Run("should draw two points", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		img := context.NewAcceleratedImage(2, 1)
		img.Upload(make([]image.Color, 2))
		vertexShader, err := context.CompileVertexShader(`
								#version 330 core
								layout(location = 0) in vec2 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition, 0, 1);
								}
								`)
		require.NoError(t, err)
		fragmentShader, err := context.CompileFragmentShader(`
								#version 330 core
								out vec4 color;
								void main() {
									color = vec4(1.0, 0.89, 0.8, 0.7);
								}
								`)
		require.NoError(t, err)
		program, err := context.LinkProgram(vertexShader, fragmentShader)
		require.NoError(t, err)
		array := context.NewVertexArray(gl.VertexLayout{gl.Vec2})
		buffer := context.NewFloatVertexBuffer(4)
		buffer.Upload(0, []float32{-0.5, 0, 0.5, 0})
		vertexPositionX := gl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 2}
		array.Set(0, vertexPositionX)
		glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.DrawArrays(array, gl.Points, 0, 2)
		}}
		command := program.AcceleratedCommand(glCommand)
		command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 2, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		assertColors(t, []image.Color{image.RGBA(255, 227, 204, 178), image.RGBA(255, 227, 204, 178)}, img)
	})
	t.Run("should panic on shader attributes and vertex array mismatch", func(t *testing.T) {
		tests := map[string]struct {
			vertexShaderSrc string
			layout          gl.VertexLayout
		}{
			"float instead of vec2": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in vec2 vertexPosition;
					void main() {
						gl_Position = vec4(vertexPosition, 0, 1);
					}
					`,
				layout: gl.VertexLayout{gl.Float},
			},
			"vec2, vec4 instead of vec2, vec3": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in vec2 vertexPosition1;
					layout(location = 1) in vec3 vertexPosition2;
					void main() {
						vec3 vertexPosition = vec3(vertexPosition1, vertexPosition2.z);
						gl_Position = vec4(vertexPosition, 1);
					}
					`,
				layout: gl.VertexLayout{gl.Vec2, gl.Vec4},
			},
			"vec3 instead of float, vec4": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in float vertexPosition1;
					layout(location = 1) in vec4  vertexPosition2;
					void main() {
						gl_Position = vec4(vertexPosition1, vertexPosition2.yzw); 
					}
					`,
				layout: gl.VertexLayout{gl.Vec3},
			},
			"vec4, vec4 instead of float": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in float vertexPosition;
					void main() {
						gl_Position = vec4(vertexPosition, 0, 0, 1); 
					}
					`,
				layout: gl.VertexLayout{gl.Vec4, gl.Vec4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				img := context.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 2))
				vertexShader, err := context.CompileVertexShader(test.vertexShaderSrc)
				require.NoError(t, err)
				fragmentShader, err := context.CompileFragmentShader(`
								#version 330 core
								void main() {}
								`)
				require.NoError(t, err)
				program, err := context.LinkProgram(vertexShader, fragmentShader)
				require.NoError(t, err)
				array := context.NewVertexArray(test.layout)
				buffer := context.NewFloatVertexBuffer(10)
				buffer.Upload(0, []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				vertexPosition := gl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 10}
				for i := range test.layout {
					array.Set(i, vertexPosition)
				}
				glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					assert.Panics(t, func() {
						renderer.DrawArrays(array, gl.Points, 0, 1)
					})
				}}
				command := program.AcceleratedCommand(glCommand)
				command.Run(image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
					Image:    img,
				}, []image.AcceleratedImageSelection{})
			})
		}
	})
	t.Run("should not return error when vertex array and program have different len of attributes", func(t *testing.T) {
		tests := map[string]struct {
			vertexShaderSrc string
			layout          gl.VertexLayout
		}{
			"len(vertex array) > len(shader)": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in vec4 vertexPosition;
					void main() {
						gl_Position = vertexPosition;
					}
					`,
				layout: gl.VertexLayout{gl.Vec4, gl.Vec4},
			},
			"len(vertex array) < len(shader)": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in vec4 vertexPosition1;
					layout(location = 1) in vec4 vertexPosition2;
					void main() {
						gl_Position = vertexPosition1 + vertexPosition2; 
					}
					`,
				layout: gl.VertexLayout{gl.Vec4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				img := context.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 2))
				vertexShader, err := context.CompileVertexShader(test.vertexShaderSrc)
				require.NoError(t, err)
				fragmentShader, err := context.CompileFragmentShader(`
								#version 330 core
								void main() {}
								`)
				require.NoError(t, err)
				program, err := context.LinkProgram(vertexShader, fragmentShader)
				require.NoError(t, err)
				array := context.NewVertexArray(test.layout)
				buffer := context.NewFloatVertexBuffer(8)
				buffer.Upload(0, []float32{0, 0, 0, 0, 0, 0, 0, 0})
				for location := range test.layout {
					array.Set(location, gl.VertexBufferPointer{Buffer: buffer, Offset: location * 4, Stride: 8})
				}
				glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					renderer.DrawArrays(array, gl.Points, 0, 1)
				}}
				command := program.AcceleratedCommand(glCommand)
				command.Run(image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
					Image:    img,
				}, []image.AcceleratedImageSelection{})
			})
		}

	})
}

func TestRenderer_BindTexture(t *testing.T) {
	t.Run("can't bind texture with uniformName not specified in program", func(t *testing.T) {
		names := []string{"foo", "bar"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				context := gl.ContextOf(openGL)
				var (
					output  = context.NewAcceleratedImage(1, 1)
					tex     = context.NewAcceleratedImage(1, 1)
					program = compileProgram(t, context,
						`#version 330 core
						void main() {
							gl_Position = vec4(0, 0, 0, 0);
						}`,
						`#version 330 core
						 uniform sampler2D tex;
						 out vec4 color;
						 void main() {
						 	color = texture(tex, vec2(0,0));
						 }`,
					)
					command = program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
						assert.Panics(t, func() {
							// when
							renderer.BindTexture(0, name, tex)
						})
					}})
				)
				command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
			})
		}
	})
	t.Run("can't bind texture created in a different context", func(t *testing.T) {
		openGL1, _ := opengl.New(mainThreadLoop)
		defer openGL1.Destroy()
		context1 := gl.ContextOf(openGL1)
		openGL2, _ := opengl.New(mainThreadLoop)
		defer openGL2.Destroy()
		context2 := gl.ContextOf(openGL2)
		var (
			output  = context1.NewAcceleratedImage(1, 1)
			tex     = context2.NewAcceleratedImage(1, 1)
			program = workingProgram(t, context1)
			command = program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
				assert.Panics(t, func() {
					// when
					renderer.BindTexture(0, "tex", tex)
				})
			}})
		)
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("can bind texture", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		var (
			output  = context.NewAcceleratedImage(1, 1)
			tex     = context.NewAcceleratedImage(1, 1)
			program = compileProgram(t, context,
				`#version 330 core
						void main() {
							gl_Position = vec4(0, 0, 0, 0);
						}`,
				`#version 330 core
						 uniform sampler2D tex;
						 out vec4 color;
						 void main() {
						 	color = texture(tex, vec2(0,0));
						 }`,
			)
			glCommand = &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
				// when
				renderer.BindTexture(0, "tex", tex)
			}}
		)
		command := program.AcceleratedCommand(glCommand)
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("should draw point by sampling texture", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		img := context.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		tex := context.NewAcceleratedImage(1, 1)
		tex.Upload([]image.Color{image.RGBA(1, 2, 3, 4)})
		program := compileProgram(t,
			context,
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
		array := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec2})
		buffer := context.NewFloatVertexBuffer(2)
		buffer.Upload(0, []float32{0.0, 0.0})
		vertexPosition := gl.VertexBufferPointer{Buffer: buffer, Stride: 2}
		array.Set(0, vertexPosition)
		glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.BindTexture(0, "tex", tex)
			renderer.DrawArrays(array, gl.Points, 0, 1)
		}}
		command := program.AcceleratedCommand(glCommand)
		command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		assertColors(t, []image.Color{image.RGBA(1, 2, 3, 4)}, img)
	})
	t.Run("should draw point by sampling 2 textures", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		img := context.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		tex1 := context.NewAcceleratedImage(1, 1)
		tex1.Upload([]image.Color{image.RGBA(5, 6, 7, 8)})
		tex2 := context.NewAcceleratedImage(1, 1)
		tex2.Upload([]image.Color{image.RGBA(9, 10, 11, 12)})
		program := compileProgram(t,
			context,
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
		array := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec2})
		buffer := context.NewFloatVertexBuffer(2)
		buffer.Upload(0, []float32{0.0, 0.0})
		vertexPosition := gl.VertexBufferPointer{Buffer: buffer, Stride: 2}
		array.Set(0, vertexPosition)
		glCommand := &command{runGL: func(renderer *gl.Renderer, _ []image.AcceleratedImageSelection) {
			// when
			renderer.BindTexture(0, "tex1", tex1)
			renderer.BindTexture(1, "tex2", tex2)
			renderer.DrawArrays(array, gl.Points, 0, 1)
		}}
		command := program.AcceleratedCommand(glCommand)
		command.Run(image.AcceleratedImageSelection{
			Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
			Image:    img,
		}, []image.AcceleratedImageSelection{})
		// then
		assertColors(t, []image.Color{image.RGBA(5+9, 6+10, 7+11, 8+12)}, img)
	})
}

func TestRenderer_SetXXX(t *testing.T) {
	openGL, _ := opengl.New(mainThreadLoop)
	defer openGL.Destroy()
	context := gl.ContextOf(openGL)
	tests := map[string]struct {
		setUniform     func(name string, renderer *gl.Renderer)
		fragmentShader string
		expectedColor  image.Color
	}{
		"Float": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetFloat(name, 1)
			},
			fragmentShader: `#version 330 core
							 uniform float attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr, 0, 0, 0); 
							 }`,
			expectedColor: image.RGBA(255, 0, 0, 0),
		},
		"Vec2": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetVec2(name, 0.11, 0.2)
			},
			fragmentShader: `#version 330 core
							 uniform vec2 attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr, 0, 0); 
							 }`,
			expectedColor: image.RGBA(28, 51, 0, 0),
		},
		"Vec3": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetVec3(name, 0.11, 0.2, 0.4)
			},
			fragmentShader: `#version 330 core
							 uniform vec3 attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr, 0); 
							 }`,
			expectedColor: image.RGBA(28, 51, 102, 0),
		},
		"Vec4": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetVec4(name, 0.11, 0.2, 0.4, 0.6)
			},
			fragmentShader: `#version 330 core
							 uniform vec4 attr;
							 out vec4 color;
							 void main() {
								color = attr; 
							 }`,
			expectedColor: image.RGBA(28, 51, 102, 153),
		},
		"Int": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetInt(name, 1)
			},
			fragmentShader: `#version 330 core
							 uniform int attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr / 255.0, 0, 0, 0); 
							 }`,
			expectedColor: image.RGBA(1, 0, 0, 0),
		},
		"IVec2": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetIVec2(name, 1, 2)
			},
			fragmentShader: `#version 330 core
							 uniform ivec2 attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr.x/255.0, attr.y/255.0, 0, 0); 
							 }`,
			expectedColor: image.RGBA(1, 2, 0, 0),
		},
		"IVec3": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetIVec3(name, 1, 2, 3)
			},
			fragmentShader: `#version 330 core
							 uniform ivec3 attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr.x/255.0, attr.y/255.0, attr.z/255.0, 0); 
							 }`,
			expectedColor: image.RGBA(1, 2, 3, 0),
		},
		"IVec4": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetIVec4(name, 1, 2, 3, 4)
			},
			fragmentShader: `#version 330 core
							 uniform ivec4 attr;
							 out vec4 color;
							 void main() {
								color = vec4(attr.x/255.0, attr.y/255.0, attr.z/255.0, attr.w/255.0); 
							 }`,
			expectedColor: image.RGBA(1, 2, 3, 4),
		},
		"Mat3": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetMat3(name, [9]float32{
					0.0, 0.11, 0.6,
					0.3, 0.2, 0.15,
					0.5, 0.4, 0.25,
				})
			},
			fragmentShader: `#version 330 core
							 uniform mat3 attr;
							 out vec4 color;
							 void main() {
								float red   = attr[0][0] + attr[1][0] + attr[2][0];
								float green = attr[0][1] + attr[1][1] + attr[2][1];
								float blue  = attr[0][2] + attr[1][2] + attr[2][2];
								color = vec4(red, green, blue, 0); 
							 }`,
			expectedColor: image.RGBA(204, 181, 255, 0),
		},
		"Mat4": {
			setUniform: func(name string, renderer *gl.Renderer) {
				renderer.SetMat4(name, [16]float32{
					0.0, 0.11, 0.34, 0.1,
					0.3, 0.2, 0.15, 0.8,
					0.5, 0.4, 0.25, 0.05,
					0.01, 0.02, 0.03, 0.04,
				})
			},
			fragmentShader: `#version 330 core
							 uniform mat4 attr;
							 out vec4 color;
							 void main() {
								float red   = attr[0][0] + attr[1][0] + attr[2][0] + attr[3][0];
								float green = attr[0][1] + attr[1][1] + attr[2][1] + attr[3][1];
								float blue  = attr[0][2] + attr[1][2] + attr[2][2] + attr[3][2];
								float alpha = attr[0][3] + attr[1][3] + attr[2][3] + attr[3][3];
								color = vec4(red, green, blue, alpha); 
							 }`,
			expectedColor: image.RGBA(207, 186, 196, 252),
		},
	}
	for attributeType, test := range tests {
		t.Run(attributeType, func(t *testing.T) {

			t.Run("should panic for invalid uniform name", func(t *testing.T) {
				names := []string{"", " ", "  ", "\n", "\t"}
				for _, name := range names {
					t.Run(name, func(t *testing.T) {
						var (
							output  = context.NewAcceleratedImage(1, 1)
							program = workingProgram(t, context)
							command = program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
								assert.Panics(t, func() {
									// when
									test.setUniform(name, renderer)
								})
							}})
						)
						command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
					})
				}
			})

			t.Run("should panic for uniform name not specified in program", func(t *testing.T) {
				names := []string{"foo", "bar"}
				for _, name := range names {
					t.Run(name, func(t *testing.T) {
						var (
							output  = context.NewAcceleratedImage(1, 1)
							program = compileProgram(t, context,
								`#version 330 core
												void main() {
													gl_Position = vec4(0, 0, 0, 0);
												}`,
								test.fragmentShader,
							)
							command = program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
								assert.Panics(t, func() {
									// when
									test.setUniform(name, renderer)
								})
							}})
						)
						command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
					})
				}
			})

			t.Run("should draw point by using uniform", func(t *testing.T) {
				img := context.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 1))
				program := compileProgram(t,
					context,
					`
					#version 330 core
					layout(location = 0) in vec2 xy;	
					void main() {
						gl_Position = vec4(xy, 0.0, 1.0);
					}
					`,
					test.fragmentShader,
				)
				array := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec2})
				buffer := context.NewFloatVertexBuffer(2)
				buffer.Upload(0, []float32{0.0, 0.0})
				vertexPosition := gl.VertexBufferPointer{Buffer: buffer, Stride: 2}
				array.Set(0, vertexPosition)
				glCommand := &command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					test.setUniform("attr", renderer)
					renderer.DrawArrays(array, gl.Points, 0, 1)
				}}
				command := program.AcceleratedCommand(glCommand)
				command.Run(image.AcceleratedImageSelection{
					Location: image.AcceleratedImageLocation{Width: 1, Height: 1},
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				assertColors(t, []image.Color{test.expectedColor}, img)
			})

		})
	}

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
