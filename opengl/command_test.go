package opengl_test

import (
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/opengl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float, opengl.Vec3})
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
				img, _ := openGL.NewAcceleratedImage(test.width, test.height)
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
				array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2})
				require.NoError(t, err)
				buffer, err := openGL.NewFloatVertexBuffer(8)
				require.NoError(t, err)
				err = buffer.Upload(0, []float32{
					-1, 1, // top-left
					1, 1, // top-right
					1, -1, // bottom-right
					-1, -1}, // bottom-left
				)
				require.NoError(t, err)
				vertexPosition := opengl.VertexBufferPointer{Buffer: buffer, Stride: 2}
				err = array.Set(0, vertexPosition)
				require.NoError(t, err)
				glCommand := &command{runGL: func(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
					// when
					renderer.DrawArrays(array, opengl.TriangleFan, 0, 4)
					return nil
				}}
				command, _ := program.AcceleratedCommand(glCommand)
				err = command.Run(image.AcceleratedImageSelection{
					Location: test.outputLocation,
					Image:    img,
				}, []image.AcceleratedImageSelection{})
				// then
				require.NoError(t, err)
				assertColors(t, test.expectedColors, img)
			})
		}
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
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Float, opengl.Vec3})
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
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
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
		array, err := openGL.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
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
