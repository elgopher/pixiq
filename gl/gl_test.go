package gl_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/image"
)

func TestNewContext(t *testing.T) {
	t.Run("should panic when api is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			gl.NewContext(nil)
		})
	})
	t.Run("should create context", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		assert.NotNil(t, context)
	})
}

func TestContextAPI(t *testing.T) {
	t.Run("should return API", func(t *testing.T) {
		api := &apiStub{}
		context := gl.NewContext(api)
		// when
		actualAPI := context.API()
		assert.Same(t, api, actualAPI)
	})
}

func TestContext_NewFloatVertexBuffer(t *testing.T) {
	t.Run("should panic when size is negative", func(t *testing.T) {
		tests := map[string]int{
			"size -1": -1,
			"size -2": -2,
		}
		for name, size := range tests {
			t.Run(name, func(t *testing.T) {
				context := gl.NewContext(apiStub{})
				// when
				assert.Panics(t, func() {
					context.NewFloatVertexBuffer(size)
				})
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
				context := gl.NewContext(apiStub{})
				// when
				buffer := context.NewFloatVertexBuffer(size)
				// then
				assert.NotNil(t, buffer)
				// and
				assert.Equal(t, size, buffer.Size())
			})
		}
	})
}

func TestFloatVertexBuffer_Upload(t *testing.T) {
	t.Run("should panic when trying to upload slice bigger than size", func(t *testing.T) {
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
				context := gl.NewContext(apiStub{})
				buffer := context.NewFloatVertexBuffer(test.size)
				assert.Panics(t, func() {
					// when
					buffer.Upload(test.offset, test.data)
				})
			})
		}
	})
	t.Run("should panic when offset is negative", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		buffer := context.NewFloatVertexBuffer(1)
		assert.Panics(t, func() {
			// when
			buffer.Upload(-1, []float32{1})
		})
	})
}

func TestFloatVertexBuffer_Download(t *testing.T) {
	t.Run("should panic when offset is negative", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		buffer := context.NewFloatVertexBuffer(1)
		defer buffer.Delete()
		output := make([]float32, 1)
		assert.Panics(t, func() {
			// when
			buffer.Download(-1, output)
		})
	})
	t.Run("should panic when buffer has been deleted", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		buffer := context.NewFloatVertexBuffer(1)
		buffer.Delete()
		output := make([]float32, 1)
		// when
		assert.Panics(t, func() {
			// when
			buffer.Download(0, output)
		})
	})
}

func TestOpenGL_NewVertexArray(t *testing.T) {
	t.Run("should panic", func(t *testing.T) {
		tests := map[string]struct {
			layout gl.VertexLayout
		}{
			"nil layout": {
				layout: nil,
			},
			"empty layout": {
				layout: gl.VertexLayout{},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				context := gl.NewContext(apiStub{})
				assert.Panics(t, func() {
					// when
					vao := context.NewVertexArray(test.layout)
					// then
					assert.Nil(t, vao)
				})
			})
		}
	})
}

func TestVertexArray_Set(t *testing.T) {
	t.Run("should panic when offset is negative", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		defer vao.Delete()
		buffer := context.NewFloatVertexBuffer(1)
		pointer := gl.VertexBufferPointer{
			Buffer: buffer,
			Offset: -1,
			Stride: 1,
		}
		assert.Panics(t, func() {
			// when
			vao.Set(0, pointer)
		})
	})
	t.Run("should panic when stride is negative", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		buffer := context.NewFloatVertexBuffer(1)
		pointer := gl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: -1,
		}
		assert.Panics(t, func() {
			// when
			vao.Set(0, pointer)
		})
	})
	t.Run("should panic when location is negative", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		buffer := context.NewFloatVertexBuffer(1)
		pointer := gl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: 1,
		}
		assert.Panics(t, func() {
			// when
			vao.Set(-1, pointer)
		})
	})
	t.Run("should panic when location is higher than number of arguments", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		buffer := context.NewFloatVertexBuffer(1)
		pointer := gl.VertexBufferPointer{
			Buffer: buffer,
			Offset: 0,
			Stride: 1,
		}
		assert.Panics(t, func() {
			// when
			vao.Set(1, pointer)
		})
	})
	t.Run("should panic when buffer is nil", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		pointer := gl.VertexBufferPointer{
			Buffer: nil,
			Offset: 0,
			Stride: 1,
		}
		assert.Panics(t, func() {
			// when
			vao.Set(0, pointer)
		})
	})
	t.Run("should panic when buffer was not created by context", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vao := context.NewVertexArray(gl.VertexLayout{gl.Float})
		vertexBufferNotCreatedInContext := &gl.FloatVertexBuffer{}
		pointer := gl.VertexBufferPointer{
			Buffer: vertexBufferNotCreatedInContext,
			Offset: 0,
			Stride: 1,
		}
		assert.Panics(t, func() {
			// when
			vao.Set(0, pointer)
		})
	})
}
func TestOpenGL_LinkProgram(t *testing.T) {
	t.Run("should panic when vertex shader is nil", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		fragmentShader, _ := context.CompileFragmentShader("")
		assert.Panics(t, func() {
			// when
			_, _ = context.LinkProgram(nil, fragmentShader)
		})

	})
	t.Run("should panic when fragment shader is nil", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		vertexShader, _ := context.CompileVertexShader("")
		assert.Panics(t, func() {
			// when
			_, _ = context.LinkProgram(vertexShader, nil)
		})
	})
}
func TestContext_NewAcceleratedImage(t *testing.T) {
	t.Run("should panic for negative width", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		// when
		assert.Panics(t, func() {
			context.NewAcceleratedImage(-1, 0)
		})
	})
	t.Run("should panic for negative height", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		assert.Panics(t, func() {
			// when
			context.NewAcceleratedImage(0, -1)
		})
	})
	t.Run("should panic for too big width", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		capabilities := context.Capabilities()
		assert.Panics(t, func() {
			// when
			context.NewAcceleratedImage(capabilities.MaxTextureSize()+1, 1)
		})
	})
	t.Run("should panic for too big height", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		capabilities := context.Capabilities()
		assert.Panics(t, func() {
			// when
			context.NewAcceleratedImage(1, capabilities.MaxTextureSize()+1)
		})
	})
}
func TestProgram_AcceleratedCommand(t *testing.T) {
	t.Run("should return command", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		program := workingProgram(context)
		// when
		cmd := program.AcceleratedCommand(&emptyCommand{})
		assert.NotNil(t, cmd)
	})
}

func TestAcceleratedCommand_Run(t *testing.T) {
	t.Run("should execute command", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		program := workingProgram(context)
		texture := context.NewAcceleratedImage(1, 1)
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
		context := gl.NewContext(apiStub{})
		program := workingProgram(context)
		command := program.AcceleratedCommand(&emptyCommand{})
		assert.Panics(t, func() {
			// when
			command.Run(image.AcceleratedImageSelection{}, []image.AcceleratedImageSelection{})
		})
	})
	t.Run("should panic when output image and program were created in different OpenGL contexts", func(t *testing.T) {
		imageContext := gl.NewContext(apiStub{})
		programContext := gl.NewContext(apiStub{})
		img := imageContext.NewAcceleratedImage(1, 1)
		program := workingProgram(programContext)
		command := program.AcceleratedCommand(&emptyCommand{})
		assert.Panics(t, func() {
			// when
			command.Run(image.AcceleratedImageSelection{
				Image: img,
			}, []image.AcceleratedImageSelection{})
		})
	})
}

func TestRenderer_BindTexture(t *testing.T) {
	t.Run("can't bind texture without uniformName", func(t *testing.T) {
		names := []string{"", " ", "  ", "\n", "\t"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				context := gl.NewContext(apiStub{})
				var (
					output  = context.NewAcceleratedImage(1, 1)
					tex     = context.NewAcceleratedImage(1, 1)
					program = workingProgram(context)
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
	t.Run("can't bind texture with negative texture unit", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		var (
			output  = context.NewAcceleratedImage(1, 1)
			tex     = context.NewAcceleratedImage(1, 1)
			program = workingProgram(context)
			command = program.AcceleratedCommand(&command{runGL: func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
				assert.Panics(t, func() {
					// when
					renderer.BindTexture(-1, "tex", tex)
				})
			}})
		)
		command.Run(image.AcceleratedImageSelection{Image: output}, []image.AcceleratedImageSelection{})
	})
}

func TestOpenGL_Error(t *testing.T) {
	t.Run("should no return error", func(t *testing.T) {
		context := gl.NewContext(apiStub{})
		// when
		err := context.Error()
		// then
		assert.NoError(t, err)
	})
}

func workingProgram(context *gl.Context) *gl.Program {
	var (
		vertexShader, _   = context.CompileVertexShader("")
		fragmentShader, _ = context.CompileFragmentShader("")
		program, _        = context.LinkProgram(vertexShader, fragmentShader)
	)
	return program
}

type emptyCommand struct{}

func (e emptyCommand) RunGL(_ *gl.Renderer, _ []image.AcceleratedImageSelection) {}

type commandMock struct {
	executionCount int
	selections     []image.AcceleratedImageSelection
	renderer       *gl.Renderer
}

func (f *commandMock) RunGL(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
	f.executionCount++
	f.selections = selections
	f.renderer = renderer
}

type command struct {
	runGL func(renderer *gl.Renderer, selections []image.AcceleratedImageSelection)
}

func (c *command) RunGL(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
	c.runGL(renderer, selections)
}

type apiStub struct{}

func (a apiStub) GenBuffers(n int32, buffers *uint32) {}

func (a apiStub) BindBuffer(target uint32, buffer uint32)                                   {}
func (a apiStub) BufferData(target uint32, size int, data unsafe.Pointer, usage uint32)     {}
func (a apiStub) BufferSubData(target uint32, offset int, size int, data unsafe.Pointer)    {}
func (a apiStub) GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {}
func (a apiStub) DeleteBuffers(n int32, buffers *uint32)                                    {}
func (a apiStub) GenVertexArrays(n int32, arrays *uint32)                                   {}
func (a apiStub) DeleteVertexArrays(n int32, arrays *uint32)                                {}
func (a apiStub) BindVertexArray(array uint32)                                              {}
func (a apiStub) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
}
func (a apiStub) EnableVertexAttribArray(index uint32)                                    {}
func (a apiStub) CreateShader(xtype uint32) uint32                                        { return 0 }
func (a apiStub) ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {}
func (a apiStub) CompileShader(shader uint32)                                             {}
func (a apiStub) GetShaderiv(shader uint32, pname uint32, params *int32) {
	const compileStatus = 0x8B81
	const glTrue = 1
	if pname == compileStatus {
		*params = glTrue
	}
}
func (a apiStub) GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {}
func (a apiStub) DeleteShader(shader uint32)                                                   {}
func (a apiStub) GoStr(cstr *uint8) string                                                     { return "" }
func (a apiStub) Strs(strs ...string) (cstrs **uint8, free func()) {
	return nil, func() {}
}
func (a apiStub) AttachShader(program uint32, shader uint32) {}
func (a apiStub) LinkProgram(program uint32)                 {}
func (a apiStub) GetProgramiv(program uint32, pname uint32, params *int32) {
	const linkStatus = 0x8B82
	const glTrue = 1
	if pname == linkStatus {
		*params = glTrue
	}
}
func (a apiStub) GetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {}
func (a apiStub) UseProgram(program uint32)                                                      {}
func (a apiStub) CreateProgram() uint32                                                          { return 0 }
func (a apiStub) GetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
}
func (a apiStub) GetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
}
func (a apiStub) GetAttribLocation(program uint32, name *uint8) int32                          { return 0 }
func (a apiStub) Enable(cap uint32)                                                            {}
func (a apiStub) BindFramebuffer(target uint32, framebuffer uint32)                            {}
func (a apiStub) Scissor(x int32, y int32, width int32, height int32)                          {}
func (a apiStub) Viewport(x int32, y int32, width int32, height int32)                         {}
func (a apiStub) ClearColor(red float32, green float32, blue float32, alpha float32)           {}
func (a apiStub) Clear(mask uint32)                                                            {}
func (a apiStub) DrawArrays(mode uint32, first int32, count int32)                             {}
func (a apiStub) Uniform1f(location int32, v0 float32)                                         {}
func (a apiStub) Uniform2f(location int32, v0 float32, v1 float32)                             {}
func (a apiStub) Uniform3f(location int32, v0 float32, v1 float32, v2 float32)                 {}
func (a apiStub) Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32)     {}
func (a apiStub) Uniform1i(location int32, v0 int32)                                           {}
func (a apiStub) Uniform2i(location int32, v0 int32, v1 int32)                                 {}
func (a apiStub) Uniform3i(location int32, v0 int32, v1 int32, v2 int32)                       {}
func (a apiStub) Uniform4i(location int32, v0 int32, v1 int32, v2 int32, v3 int32)             {}
func (a apiStub) UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {}
func (a apiStub) UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {}
func (a apiStub) ActiveTexture(texture uint32)                                                 {}
func (a apiStub) BindTexture(target uint32, texture uint32)                                    {}
func (a apiStub) GetIntegerv(pname uint32, data *int32) {
	const maxTextureSize = 0x0D33
	if pname == maxTextureSize {
		*data = 1024 * 1024
	}
}
func (a apiStub) GenTextures(n int32, textures *uint32) {}
func (a apiStub) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
}
func (a apiStub) TexParameteri(target uint32, pname uint32, param int32) {}
func (a apiStub) GenFramebuffers(n int32, framebuffers *uint32)          {}
func (a apiStub) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
}
func (a apiStub) TexSubImage2D(target uint32, level int32, xoffset int32, yoffset int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
}
func (a apiStub) GetTexImage(target uint32, level int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
}
func (a apiStub) GetError() uint32 { return 0 }
func (a apiStub) ReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
}
func (a apiStub) Finish()                             {}
func (a apiStub) Ptr(data interface{}) unsafe.Pointer { return nil }
func (a apiStub) PtrOffset(offset int) unsafe.Pointer { return nil }
