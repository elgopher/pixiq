package opengl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/opengl"
)

func TestAcceleratedCommand_Run(t *testing.T) {
	t.Run("should execute command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		texture := openGL.NewAcceleratedImage(1, 1)
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
				acceleratedCommand := program.AcceleratedCommand(command)
				// when
				acceleratedCommand.Run(output, test.selections)
				// then
				assert.Equal(t, 1, command.executionCount)
				assert.Equal(t, test.selections, command.selections)
				assert.NotNil(t, command.renderer)
			})
		}
	})
	t.Run("should panic when output image is nil", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		command := program.AcceleratedCommand(&emptyCommand{})
		assert.Panics(t, func() {
			// when
			command.Run(image.AcceleratedImageSelection{}, []image.AcceleratedImageSelection{})
		})
	})
	t.Run("should panic when output image and program were created in different OpenGL contexts", func(t *testing.T) {
		imageContext, _ := opengl.New(mainThreadLoop)
		defer imageContext.Destroy()
		programContext, _ := opengl.New(mainThreadLoop)
		defer programContext.Destroy()
		img := imageContext.NewAcceleratedImage(1, 1)
		program := workingProgram(t, programContext)
		command := program.AcceleratedCommand(&emptyCommand{})
		assert.Panics(t, func() {
			// when
			command.Run(image.AcceleratedImageSelection{
				Image: img,
			}, []image.AcceleratedImageSelection{})
		})
	})
	t.Run("vertex buffer can be used inside command", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		program := workingProgram(t, openGL)
		output := openGL.NewAcceleratedImage(1, 1)
		command := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			buffer := openGL.NewFloatVertexBuffer(1)
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
		program := workingProgram(t, openGL)
		output := openGL.NewAcceleratedImage(1, 1)
		command := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			array := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float})
			defer array.Delete()
			buffer := openGL.NewFloatVertexBuffer(1)
			defer buffer.Delete()
			array.Set(0, opengl.VertexBufferPointer{
				Buffer: buffer,
				Offset: 0,
				Stride: 1,
			})
		}})
		// when
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
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
				img := openGL.NewAcceleratedImage(test.width, test.height)
				img.Upload(make([]image.Color, test.width*test.height))
				program := workingProgram(t, openGL)
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
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
		commands := map[string]opengl.Command{
			"nil":   nil,
			"empty": &emptyCommand{},
		}
		for name, command := range commands {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img := openGL.NewAcceleratedImage(2, 1)
				pixels := []image.Color{image.RGB(1, 2, 3), image.RGB(4, 5, 6)}
				img.Upload(pixels)
				program := workingProgram(t, openGL)
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
			img := openGL.NewAcceleratedImage(1, 1)
			img.Upload(make([]image.Color, 1))
			program := workingProgram(t, openGL)
			glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
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
			"vec2": {
				vertexShaderSrc: `
								#version 330 core
								layout(location = 0) in vec2 vertexPosition;
								void main() {
									gl_Position = vec4(vertexPosition.x-1, vertexPosition.y-2, 0, 1);
								}
								`,
				typ:  opengl.Vec2,
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
				typ:  opengl.Vec3,
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
				typ:  opengl.Vec4,
				data: []float32{1, 2, 3, 4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img := openGL.NewAcceleratedImage(1, 1)
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
				array := openGL.NewVertexArray(opengl.VertexLayout{test.typ})
				buffer := openGL.NewFloatVertexBuffer(len(test.data))
				buffer.Upload(0, test.data)
				vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: len(test.data)}
				array.Set(0, vertexPosition)
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					renderer.DrawArrays(array, opengl.Points, 0, 1)
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
		img := openGL.NewAcceleratedImage(1, 1)
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
		array := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float, opengl.Vec3})
		require.NoError(t, err)
		buffer := openGL.NewFloatVertexBuffer(4)
		buffer.Upload(0, []float32{0, 0.2, 0.4, 0.6})
		vertexPositionX := opengl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 4}
		array.Set(0, vertexPositionX)
		vertexColor := opengl.VertexBufferPointer{Buffer: buffer, Offset: 1, Stride: 4}
		array.Set(1, vertexColor)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.DrawArrays(array, opengl.Points, 0, 1)
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
				expectedColors: []image.Color{image.Transparent, color},
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
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img := openGL.NewAcceleratedImage(test.width, test.height)
				img.Upload(make([]image.Color, test.width*test.height))
				program := compileProgram(t, openGL,
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
				array := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2})
				buffer := openGL.NewFloatVertexBuffer(8)
				buffer.Upload(0, []float32{
					-1, 1, // top-left
					1, 1, // top-right
					1, -1, // bottom-right
					-1, -1}, // bottom-left
				)
				vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: 2}
				array.Set(0, vertexPosition)
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					renderer.DrawArrays(array, opengl.TriangleFan, 0, 4)
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
		img := openGL.NewAcceleratedImage(2, 1)
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
		array := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2})
		buffer := openGL.NewFloatVertexBuffer(4)
		buffer.Upload(0, []float32{-0.5, 0, 0.5, 0})
		vertexPositionX := opengl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 2}
		array.Set(0, vertexPositionX)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.DrawArrays(array, opengl.Points, 0, 2)
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
			layout          opengl.VertexLayout
		}{
			"float instead of vec2": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in vec2 vertexPosition;
					void main() {
						gl_Position = vec4(vertexPosition, 0, 1);
					}
					`,
				layout: opengl.VertexLayout{opengl.Float},
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
				layout: opengl.VertexLayout{opengl.Vec2, opengl.Vec4},
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
				layout: opengl.VertexLayout{opengl.Vec3},
			},
			"vec4, vec4 instead of float": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in float vertexPosition;
					void main() {
						gl_Position = vec4(vertexPosition, 0, 0, 1); 
					}
					`,
				layout: opengl.VertexLayout{opengl.Vec4, opengl.Vec4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img := openGL.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 2))
				vertexShader, err := openGL.CompileVertexShader(test.vertexShaderSrc)
				require.NoError(t, err)
				fragmentShader, err := openGL.CompileFragmentShader(`
								#version 330 core
								void main() {}
								`)
				require.NoError(t, err)
				program, err := openGL.LinkProgram(vertexShader, fragmentShader)
				require.NoError(t, err)
				array := openGL.NewVertexArray(test.layout)
				buffer := openGL.NewFloatVertexBuffer(10)
				buffer.Upload(0, []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Offset: 0, Stride: 10}
				for i := range test.layout {
					array.Set(i, vertexPosition)
				}
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					assert.Panics(t, func() {
						renderer.DrawArrays(array, opengl.Points, 0, 1)
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
			layout          opengl.VertexLayout
		}{
			"len(vertex array) > len(shader)": {
				vertexShaderSrc: `
					#version 330 core
					layout(location = 0) in vec4 vertexPosition;
					void main() {
						gl_Position = vertexPosition;
					}
					`,
				layout: opengl.VertexLayout{opengl.Vec4, opengl.Vec4},
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
				layout: opengl.VertexLayout{opengl.Vec4},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				img := openGL.NewAcceleratedImage(1, 1)
				img.Upload(make([]image.Color, 2))
				vertexShader, err := openGL.CompileVertexShader(test.vertexShaderSrc)
				require.NoError(t, err)
				fragmentShader, err := openGL.CompileFragmentShader(`
								#version 330 core
								void main() {}
								`)
				require.NoError(t, err)
				program, err := openGL.LinkProgram(vertexShader, fragmentShader)
				require.NoError(t, err)
				array := openGL.NewVertexArray(test.layout)
				buffer := openGL.NewFloatVertexBuffer(8)
				buffer.Upload(0, []float32{0, 0, 0, 0, 0, 0, 0, 0})
				for location := range test.layout {
					array.Set(location, opengl.VertexBufferPointer{Buffer: buffer, Offset: location * 4, Stride: 8})
				}
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
					// when
					renderer.DrawArrays(array, opengl.Points, 0, 1)
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
	t.Run("can't bind texture without uniformName", func(t *testing.T) {
		names := []string{"", " ", "  ", "\n", "\t"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				output := openGL.NewAcceleratedImage(1, 1)
				tex := openGL.NewAcceleratedImage(1, 1)
				program := workingProgram(t, openGL)
				command := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
					assert.Panics(t, func() {
						// when
						renderer.BindTexture(0, name, tex)
					})
				}})
				command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
			})
		}
	})
	t.Run("can't bind texture with uniformName not specified in program", func(t *testing.T) {
		names := []string{"foo", "bar"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				openGL, _ := opengl.New(mainThreadLoop)
				defer openGL.Destroy()
				output := openGL.NewAcceleratedImage(1, 1)
				tex := openGL.NewAcceleratedImage(1, 1)
				program := compileProgram(t, openGL,
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
				command := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
					assert.Panics(t, func() {
						// when
						renderer.BindTexture(0, name, tex)
					})
				}})
				command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
			})
		}
	})
	t.Run("can't bind texture created in a different context", func(t *testing.T) {
		openGL1, _ := opengl.New(mainThreadLoop)
		defer openGL1.Destroy()
		openGL2, _ := opengl.New(mainThreadLoop)
		defer openGL2.Destroy()
		output := openGL1.NewAcceleratedImage(1, 1)
		tex := openGL2.NewAcceleratedImage(1, 1)
		program := workingProgram(t, openGL1)
		command := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			assert.Panics(t, func() {
				// when
				renderer.BindTexture(0, "tex", tex)
			})
		}})
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("can't bind texture with negative texture unit", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		output := openGL.NewAcceleratedImage(1, 1)
		tex := openGL.NewAcceleratedImage(1, 1)
		program := workingProgram(t, openGL)
		command := program.AcceleratedCommand(&command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			assert.Panics(t, func() {
				// when
				renderer.BindTexture(-1, "tex", tex)
			})
		}})
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("can bind texture", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		output := openGL.NewAcceleratedImage(1, 1)
		tex := openGL.NewAcceleratedImage(1, 1)
		program := compileProgram(t, openGL,
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
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.BindTexture(0, "tex", tex)
		}}
		command := program.AcceleratedCommand(glCommand)
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
	t.Run("should draw point by sampling texture", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		img := openGL.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		tex := openGL.NewAcceleratedImage(1, 1)
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
		array := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
		buffer := openGL.NewFloatVertexBuffer(2)
		buffer.Upload(0, []float32{0.0, 0.0})
		vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: 2}
		array.Set(0, vertexPosition)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
			// when
			renderer.BindTexture(0, "tex", tex)
			renderer.DrawArrays(array, opengl.Points, 0, 1)
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
		img := openGL.NewAcceleratedImage(1, 1)
		img.Upload(make([]image.Color, 1))
		tex1 := openGL.NewAcceleratedImage(1, 1)
		tex1.Upload([]image.Color{image.RGBA(5, 6, 7, 8)})
		tex2 := openGL.NewAcceleratedImage(1, 1)
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
		array := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
		buffer := openGL.NewFloatVertexBuffer(2)
		buffer.Upload(0, []float32{0.0, 0.0})
		vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: 2}
		array.Set(0, vertexPosition)
		glCommand := &command{runGL: func(renderer *opengl.Renderer, _ []image.AcceleratedImageSelection) {
			// when
			renderer.BindTexture(0, "tex1", tex1)
			renderer.BindTexture(1, "tex2", tex2)
			renderer.DrawArrays(array, opengl.Points, 0, 1)
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
