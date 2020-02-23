package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func newScreenPolygon(context *context) *screenPolygon {
	const (
		vertexLocation  = 0
		textureLocation = 1
	)
	var vertexArrayID, vertexBufferID uint32
	context.GenVertexArrays(1, &vertexArrayID)
	context.BindVertexArray(vertexArrayID)
	context.GenBuffers(1, &vertexBufferID)
	context.BindBuffer(gl.ARRAY_BUFFER, vertexBufferID)
	data := []float32{
		-1, 1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
		1, 1, 1, 1,
		1, -1, 1, 0,
		-1, -1, 0, 0,
	}
	context.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
	const stride int32 = 4 * 4
	const vec2size int32 = 2
	context.VertexAttribPointer(
		vertexLocation,
		vec2size,
		gl.FLOAT,
		false,
		stride,
		gl.PtrOffset(0),
	)
	context.EnableVertexAttribArray(0)
	context.VertexAttribPointer(
		textureLocation,
		vec2size,
		gl.FLOAT,
		false,
		stride,
		gl.PtrOffset(8),
	)
	context.EnableVertexAttribArray(1)
	return &screenPolygon{vertexArrayID: vertexArrayID, context: context}
}

type screenPolygon struct {
	vertexArrayID uint32
	context       *context
}

func (p *screenPolygon) draw() {
	p.context.BindVertexArray(p.vertexArrayID)
	p.context.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}
