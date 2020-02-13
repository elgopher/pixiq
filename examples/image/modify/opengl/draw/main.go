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
		window, err := gl.OpenWindow(200, 200, opengl.Zoom(2))
		if err != nil {
			panic(err)
		}
		// xy -> rgb
		vertices := []float32{
			-1, 1, 1, 0, 0, // top-left -> red
			1, 1, 0, 1, 0, // top-right -> green
			1, -1, 0, 0, 1, // bottom-right -> blue
			-1, -1, 1, 1, 1, // bottom-left -> white
		}
		buffer, err := gl.NewFloatVertexBuffer(len(vertices))
		if err != nil {
			panic(err)
		}
		if err := buffer.Upload(0, vertices); err != nil {
			panic(err)
		}
		array, err := gl.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec3})
		if err != nil {
			panic(err)
		}
		xy := opengl.VertexBufferPointer{Offset: 0, Stride: 5, Buffer: buffer}
		if err := array.Set(0, xy); err != nil {
			panic(err)
		}
		color := opengl.VertexBufferPointer{Offset: 2, Stride: 5, Buffer: buffer}
		if err := array.Set(1, color); err != nil {
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
	layout(location = 1) in vec3 vertexColor;
	out vec4 interpolatedColor;
	
	void main() {
		gl_Position = vec4(vertexPosition, 0.0, 1.0);
		interpolatedColor = vec4(vertexColor, 1.0);
	}
`

const fragmentShaderSrc = `
	#version 330 core
	
	in vec4 interpolatedColor;
	out vec4 color;
	
	void main() {
		color = interpolatedColor;
	}
`

type command struct {
	vertexArray *opengl.VertexArray
}

func (c command) RunGL(renderer *opengl.Renderer, _ []image.AcceleratedImageSelection) error {
	return renderer.DrawArrays(c.vertexArray, opengl.TriangleFan, 0, 4)
}
