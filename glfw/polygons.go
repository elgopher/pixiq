package glfw

import (
	gl33 "github.com/go-gl/gl/v3.3-core/gl"

	"github.com/jacekolszak/pixiq/gl"
)

func newScreenPolygon(context *gl.Context) *screenPolygon {
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
	buffer := context.NewFloatVertexBuffer(len(data), gl.DynamicDraw)
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
	return &screenPolygon{vao: vao, vbo: buffer, rect: data, api: context.API()}
}

type screenPolygon struct {
	vao  *gl.VertexArray
	vbo  *gl.FloatVertexBuffer
	rect rect
	api  gl.API
}

func (p *screenPolygon) draw(xRight, yBottom float32) {
	p.rect.SetTopRight(xRight, 1)
	p.rect.SetBottomRight(xRight, yBottom)
	p.rect.SetBottomLeft(-1, yBottom)
	p.vbo.Upload(0, p.rect)
	p.api.BindVertexArray(p.vao.ID())
	p.api.DrawArrays(gl33.TRIANGLE_FAN, 0, 4)
}

func (p *screenPolygon) delete() {
	p.vao.Delete()
	p.vbo.Delete()
}

type rect []float32

func (r rect) SetTopRight(x, y float32) {
	r[4] = x
	r[5] = y
}

func (r rect) SetBottomRight(x, y float32) {
	r[8] = x
	r[9] = y
}

func (r rect) SetBottomLeft(x, y float32) {
	r[12] = x
	r[13] = y
}
