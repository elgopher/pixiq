package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func newScreenPolygon(vertexPositionLocation int32, texturePositionLocation int32) *screenPolygon {
	var vertexArrayID uint32
	var vertexBufferID uint32
	gl.GenVertexArrays(1, &vertexArrayID)
	gl.BindVertexArray(vertexArrayID)
	gl.GenBuffers(1, &vertexBufferID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferID)
	data := []float32{
		-1, 1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
		1, 1, 1, 1,
		1, -1, 1, 0,
		-1, -1, 0, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
	const stride int32 = 4 * 4
	const vec2size int32 = 2
	gl.VertexAttribPointer(
		uint32(vertexPositionLocation),
		vec2size,
		gl.FLOAT,
		false,
		stride,
		gl.PtrOffset(0),
	)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(
		uint32(texturePositionLocation),
		vec2size,
		gl.FLOAT,
		false,
		stride,
		gl.PtrOffset(8),
	)
	gl.EnableVertexAttribArray(1)
	return &screenPolygon{vertexArrayID: vertexArrayID, vertexBufferID: vertexBufferID}
}

type screenPolygon struct {
	vertexArrayID  uint32
	vertexBufferID uint32
}

func (p *screenPolygon) draw() {
	gl.BindVertexArray(p.vertexArrayID)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}
