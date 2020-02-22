package opengl_test

import (
	"github.com/jacekolszak/pixiq/gl"
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
	t.Run("should panic for too big width", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		capabilities := context.Capabilities()
		assert.Panics(t, func() {
			// when
			context.NewAcceleratedImage(capabilities.MaxTextureSize()+1, 1)
		})
	})
	t.Run("should panic for too big height", func(t *testing.T) {
		openGL, _ := opengl.New(mainThreadLoop)
		defer openGL.Destroy()
		context := gl.ContextOf(openGL)
		capabilities := context.Capabilities()
		assert.Panics(t, func() {
			// when
			context.NewAcceleratedImage(1, capabilities.MaxTextureSize()+1)
		})
	})
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
