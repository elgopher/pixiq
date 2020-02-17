package main

import (
	"github.com/jacekolszak/pixiq/colornames"
	"github.com/jacekolszak/pixiq/fill"
	"log"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		var (
			buffer    = makeVertexBuffer(gl)
			array     = makeVertexArray(gl, buffer)
			program   = compileProgram(gl)
			cmd       = program.AcceleratedCommand(&drawSelection{vertexArray: array})
			window    = openWindow(gl)
			fillColor = fill.New(program)
		)
		sampledImage := gl.NewImage(2, 2)
		selection := sampledImage.WholeImageSelection()
		selection.SetColor(0, 0, image.RGB(255, 0, 0))
		selection.SetColor(1, 0, image.RGB(0, 255, 0))
		selection.SetColor(0, 1, image.RGB(255, 255, 255))
		selection.SetColor(1, 1, image.RGB(0, 0, 255))

		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			screen.Modify(cmd, selection)
			fillColor.Fill(screen.Selection(0, 0).WithSize(60, 60), colornames.Darkred)
			screen.SetColor(10, 10, colornames.Blue)
			screen.SetColor(11, 11, colornames.Blue)
		})
	})
}

func makeVertexArray(gl *opengl.OpenGL, buffer *opengl.FloatVertexBuffer) *opengl.VertexArray {
	array := gl.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
	xy := opengl.VertexBufferPointer{Offset: 0, Stride: 4, Buffer: buffer}
	array.Set(0, xy)
	st := opengl.VertexBufferPointer{Offset: 2, Stride: 4, Buffer: buffer}
	array.Set(1, st)
	return array
}

func makeVertexBuffer(gl *opengl.OpenGL) *opengl.FloatVertexBuffer {
	// xy -> st
	vertices := []float32{
		-1, 1, 0, 1, // top-left
		1, 1, 1, 1, // top-right
		1, -1, 1, 0, // bottom-right
		-1, -1, 0, 0, // bottom-left
	}
	buffer := gl.NewFloatVertexBuffer(len(vertices))
	buffer.Upload(0, vertices)
	return buffer
}

const vertexShaderSrc = `
	#version 330 core
	
	layout(location = 0) in vec2 xy;
	layout(location = 1) in vec2 st;
	out vec2 interpolatedST;
	
	void main() {
		gl_Position = vec4(xy, 0.0, 1.0);
		interpolatedST = st;
	}
`

const fragmentShaderSrc = `
	#version 330 core
	
	uniform sampler2D tex;
	in vec2 interpolatedST;
	out vec4 color;
	
	void main() {
		color = texture(tex, interpolatedST);
	}
`

func compileProgram(gl *opengl.OpenGL) *opengl.Program {
	vertexShader, err := gl.CompileVertexShader(vertexShaderSrc)
	if err != nil {
		log.Panicf("CompileVertexShader failed: %v", err)
	}
	fragmentShader, err := gl.CompileFragmentShader(fragmentShaderSrc)
	if err != nil {
		log.Panicf("CompileFragmentShader failed: %v", err)
	}
	program, err := gl.LinkProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Panicf("LinkProgram failed: %v", err)
	}
	return program
}

type drawSelection struct {
	vertexArray *opengl.VertexArray
}

func (c drawSelection) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) {
	if len(selections) != 1 {
		panic("invalid number of selections")
	}
	renderer.BindTexture(0, "tex", selections[0].Image)
	renderer.DrawArrays(c.vertexArray, opengl.TriangleFan, 0, 4)
}

func openWindow(gl *opengl.OpenGL) *opengl.Window {
	window, err := gl.OpenWindow(200, 200, opengl.Zoom(2))
	if err != nil {
		log.Panicf("OpenWindow failed: %v", err)
	}
	return window
}
