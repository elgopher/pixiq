package opengl_test

import (
	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/opengl"
	"github.com/stretchr/testify/assert"
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
