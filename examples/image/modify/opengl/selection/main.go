package main

import (
	"log"

	"github.com/elgopher/pixiq/gl"
	"github.com/elgopher/pixiq/glfw"
	"github.com/elgopher/pixiq/image"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		var (
			context = openGL.Context()
			buffer  = makeVertexBuffer(context)
			array   = makeVertexArray(context, buffer)
			program = compileProgram(context)
			cmd     = program.AcceleratedCommand(&drawSelection{vertexArray: array})
			window  = openWindow(openGL)
		)
		sampledImage := openGL.NewImage(2, 2)
		selection := sampledImage.WholeImageSelection()
		selection.SetColor(0, 0, image.RGB(255, 0, 0))
		selection.SetColor(1, 0, image.RGB(0, 255, 0))
		selection.SetColor(0, 1, image.RGB(255, 255, 255))
		selection.SetColor(1, 1, image.RGB(0, 0, 255))

		for {
			window.Screen().Modify(cmd, selection)
			window.Draw()
			if window.ShouldClose() {
				break
			}
		}
	})
}

func makeVertexArray(context *gl.Context, buffer *gl.FloatVertexBuffer) *gl.VertexArray {
	array := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec2})
	xy := gl.VertexBufferPointer{Offset: 0, Stride: 4, Buffer: buffer}
	array.Set(0, xy)
	st := gl.VertexBufferPointer{Offset: 2, Stride: 4, Buffer: buffer}
	array.Set(1, st)
	return array
}

func makeVertexBuffer(ctx *gl.Context) *gl.FloatVertexBuffer {
	// xy -> st
	vertices := []float32{
		-1, 1, 0, 1, // top-left
		1, 1, 1, 1, // top-right
		1, -1, 1, 0, // bottom-right
		-1, -1, 0, 0, // bottom-left
	}
	buffer := ctx.NewFloatVertexBuffer(len(vertices), gl.StaticDraw)
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

func compileProgram(context *gl.Context) *gl.Program {
	vertexShader, err := context.CompileVertexShader(vertexShaderSrc)
	if err != nil {
		log.Panicf("CompileVertexShader failed: %v", err)
	}
	fragmentShader, err := context.CompileFragmentShader(fragmentShaderSrc)
	if err != nil {
		log.Panicf("CompileFragmentShader failed: %v", err)
	}
	program, err := context.LinkProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Panicf("LinkProgram failed: %v", err)
	}
	return program
}

type drawSelection struct {
	vertexArray *gl.VertexArray
}

func (c drawSelection) RunGL(renderer *gl.Renderer, selections []image.AcceleratedImageSelection) {
	if len(selections) != 1 {
		panic("invalid number of selections")
	}
	renderer.BindTexture(0, "tex", selections[0].Image)
	renderer.DrawArrays(c.vertexArray, gl.TriangleFan, 0, 4)
}

func openWindow(openGL *glfw.OpenGL) *glfw.Window {
	window, err := openGL.OpenWindow(200, 200, glfw.Zoom(2))
	if err != nil {
		log.Panicf("OpenWindow failed: %v", err)
	}
	return window
}
