package gl_test

import (
	"github.com/jacekolszak/pixiq/gl"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestContextOf(t *testing.T) {
	t.Run("should panic when api is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			gl.ContextOf(nil)
		})
	})
	t.Run("should create context", func(t *testing.T) {
		context := gl.ContextOf(apiStub{})
		assert.NotNil(t, context)
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
				context := gl.ContextOf(apiStub{})
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
				context := gl.ContextOf(apiStub{})
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
				context := gl.ContextOf(apiStub{})
				buffer := context.NewFloatVertexBuffer(test.size)
				assert.Panics(t, func() {
					// when
					buffer.Upload(test.offset, test.data)
				})
			})
		}
	})
	t.Run("should panic when offset is negative", func(t *testing.T) {
		context := gl.ContextOf(apiStub{})
		buffer := context.NewFloatVertexBuffer(1)
		assert.Panics(t, func() {
			// when
			buffer.Upload(-1, []float32{1})
		})
	})
}

type apiStub struct {
}

func (a apiStub) GenBuffers(n int32, buffers *uint32) {
}

func (a apiStub) BindBuffer(target uint32, buffer uint32) {
}

func (a apiStub) BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
}

func (a apiStub) BufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
}

func (a apiStub) GetBufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
}

func (a apiStub) DeleteBuffers(n int32, buffers *uint32) {
}
