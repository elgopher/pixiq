package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func newFrameImagePolygon(loop *MainThreadLoop) *frameImagePolygon {
	var vertexArrayId uint32
	var vertexBufferId uint32
	loop.Execute(func() {
		gl.GenVertexArrays(1, &vertexArrayId)
		gl.BindVertexArray(vertexArrayId)
		gl.GenBuffers(1, &vertexBufferId)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferId)
		data := []float32{
			-1, -1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
			1, -1, 1, 1,
			1, 1, 1, 0,
			//
			-1, -1, 0, 1,
			1, 1, 1, 0,
			-1, 1, 0, 0,
		}
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
		// vertexPosition (attribute 0)
		const stride int32 = 4 * 4
		const vec2size int32 = 2
		gl.VertexAttribPointer(
			0,
			vec2size,
			gl.FLOAT,
			false,
			stride,
			gl.PtrOffset(0),
		)
		gl.EnableVertexAttribArray(0)
		// texturePosition (attribute 1)
		gl.VertexAttribPointer(
			1,
			vec2size,
			gl.FLOAT,
			false,
			stride,
			gl.PtrOffset(8),
		)
		gl.EnableVertexAttribArray(1)
	})
	return &frameImagePolygon{vertexArrayId: vertexArrayId, vertexBufferId: vertexBufferId}
}

type frameImagePolygon struct {
	vertexArrayId  uint32
	vertexBufferId uint32
}

func (p *frameImagePolygon) draw() {
	gl.BindBuffer(gl.ARRAY_BUFFER, p.vertexBufferId)
	gl.BindVertexArray(p.vertexArrayId)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
