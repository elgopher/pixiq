package main

import (
	"errors"

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
		window, err := gl.OpenWindow(1280, 720)
		if err != nil {
			panic(err)
		}
		vertices := []float32{
			-1, -1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
			1, -1, 1, 1,
			1, 1, 1, 0,
			-1, -1, 0, 1,
			1, 1, 1, 0,
			-1, 1, 0, 0,
		}
		buffer, err := gl.NewFloatVertexBuffer(len(vertices))
		if err != nil {
			panic(err)
		}
		if err := buffer.Upload(0, vertices); err != nil {
			panic(err)
		}
		array, err := gl.NewVertexArray(opengl.VertexLayout{opengl.Float2, opengl.Float2})
		if err != nil {
			panic(err)
		}
		xy := opengl.VertexBufferPointer{
			Offset: 0,
			Stride: 4,
			Buffer: buffer,
		}
		if err := array.Set(0, xy); err != nil {
			panic(err)
		}
		uv := opengl.VertexBufferPointer{
			Offset: 2,
			Stride: 4,
			Buffer: buffer,
		}
		if err := array.Set(1, uv); err != nil {
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
			if err := screen.Modify(cmd, screen); err != nil {
				panic(err)
			}
		})
	})
}

const vertexShaderSrc = `
#version 330 core

layout(location = 0) in vec2 vertexPosition;
layout(location = 1) in vec2 texturePosition;

out vec2 position;

void main() {
	gl_Position = vec4(vertexPosition, 0.0, 1.0);
	position = texturePosition;
}
`

const fragmentShaderSrc = `
#version 330 core

in vec2 position;

out vec4 color;

uniform sampler2D tex;

void main() {
	color = texture(tex, position);
}
`

type command struct {
	vertexArray *opengl.VertexArray
}

func (c command) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
	if len(selections) != 1 {
		return errors.New("invalid number of selections")
	}
	if err := renderer.BindTexture(0, "tex", selections[0].Image); err != nil {
		return err
	}
	renderer.DrawArrays(c.vertexArray, opengl.Triangles, 0, 6)
	return nil
}
