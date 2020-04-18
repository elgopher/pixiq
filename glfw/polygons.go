package glfw

import (
	gl33 "github.com/go-gl/gl/v3.3-core/gl"

	"github.com/jacekolszak/pixiq/gl"
)

func newScreenPolygon(context *gl.Context, api gl.API) *screenPolygon {
	const (
		vertexLocation  = 0
		textureLocation = 1
	)
	data := []float32{
		-1, 1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
		1, 1, 1, 1,
		1, -1, 1, 0,
		-1, -1, 0, 0,
	}
	buffer := context.NewFloatVertexBuffer(len(data), gl.Static)
	buffer.Upload(0, data)

	vao := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec2})
	vao.Set(vertexLocation, gl.VertexBufferPointer{
		Buffer: buffer,
		Offset: 0,
		Stride: 4,
	})
	vao.Set(textureLocation, gl.VertexBufferPointer{
		Buffer: buffer,
		Offset: 2,
		Stride: 4,
	})
	return &screenPolygon{vertexArrayID: vao.ID(), api: api}
}

type screenPolygon struct {
	vertexArrayID uint32
	api           gl.API
}

func (p *screenPolygon) draw() {
	p.api.BindVertexArray(p.vertexArrayID)
	p.api.DrawArrays(gl33.TRIANGLE_FAN, 0, 4)
}
