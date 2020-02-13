package main

import (
	"errors"
	"fmt"

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
		// xy -> st
		vertices := []float32{
			-1, 1, 0, 1, // top-left
			1, 1, 1, 1, // top-right
			1, -1, 1, 0, // bottom-right
			-1, -1, 0, 0, // bottom-left
		}
		buffer, err := gl.NewFloatVertexBuffer(len(vertices))
		if err != nil {
			panic(err)
		}
		if err := buffer.Upload(0, vertices); err != nil {
			panic(err)
		}
		array, err := gl.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
		if err != nil {
			panic(err)
		}
		xy := opengl.VertexBufferPointer{Offset: 0, Stride: 4, Buffer: buffer}
		if err := array.Set(0, xy); err != nil {
			panic(err)
		}
		st := opengl.VertexBufferPointer{Offset: 2, Stride: 4, Buffer: buffer}
		if err := array.Set(1, st); err != nil {
			panic(err)
		}
		cmd, err := program.AcceleratedCommand(&command{
			vertexArray: array,
		})
		if err != nil {
			panic(err)
		}
		sampledImage, err := gl.NewImage(2, 2)
		if err != nil {
			panic(err)
		}
		selection := sampledImage.WholeImageSelection()
		selection.SetColor(0, 0, image.RGB(255, 0, 0))
		selection.SetColor(1, 0, image.RGB(0, 255, 0))
		selection.SetColor(0, 1, image.RGB(255, 255, 255))
		selection.SetColor(1, 1, image.RGB(0, 0, 255))
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			if err := screen.Modify(cmd, selection); err != nil {
				panic(err)
			}
		})
	})
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

type command struct {
	vertexArray *opengl.VertexArray
}

func (c command) RunGL(renderer *opengl.Renderer, selections []image.AcceleratedImageSelection) error {
	if len(selections) != 1 {
		return errors.New("invalid number of selections")
	}
	if err := renderer.BindTexture(0, "tex", selections[0].Image); err != nil {
		return fmt.Errorf("error binding texture: %v", err)
	}
	if err := renderer.DrawArrays(c.vertexArray, opengl.TriangleFan, 0, 4); err != nil {
		return fmt.Errorf("error drawing arrays: %v", err)
	}
	return nil
}
