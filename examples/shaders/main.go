package main

import (
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		vertexShader, err := gl.CompileVertexShader(vertexShaderSrc)
		if err != nil {
			panic(err)
		}
		fragmentShader, err := gl.CompileFragmentShader(fragmentShaderSrc)
		if err != nil {
			panic(err)
		}
		program, err := gl.LinkProgram(vertexShader, fragmentShader)
		if err != nil {
			panic(err)
		}
		window, err := gl.OpenWindow(20, 10, opengl.Zoom(32))
		if err != nil {
			panic(err)
		}
		vertices := []float32{
			-1, -1, // x, y
			1, -1,
			1, 1,
			-1, -1,
			1, 1,
			-1, 1,
		}
		buffer, err := gl.NewFloatVertexBuffer(len(vertices))
		if err != nil {
			panic(err)
		}
		if err := buffer.Upload(0, vertices); err != nil {
			panic(err)
		}
		array, err := gl.NewVertexArray(opengl.VertexLayout{opengl.Float2})
		if err != nil {
			panic(err)
		}
		xy := opengl.VertexBufferPointer{Offset: 0, Stride: 2, Buffer: buffer}
		if err := array.Set(0, xy); err != nil {
			panic(err)
		}
		cmd, err := program.AcceleratedCommand(&command{
			vertexArray: array,
		})
		if err != nil {
			panic(err)
		}
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			if err := screen.Modify(cmd); err != nil {
				panic(err)
			}
		})
	})
}

const vertexShaderSrc = `
	#version 330 core
	
	layout(location = 0) in vec2 vertexPosition;
	
	void main() {
		gl_Position = vec4(vertexPosition, 0.0, 1.0);
	}
`

const fragmentShaderSrc = `
	#version 330 core
	
	out vec4 color;
	
	void main() {
		color = vec4(1.,0.,0.,1.);
	}
`

type command struct {
	vertexArray *opengl.VertexArray
}

func (c command) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
	renderer.DrawArrays(c.vertexArray, opengl.Triangles, 0, 6)
	return nil
}
