package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		var (
			buffer  = makeVertexBuffer(gl)
			array   = makeVertexArray(gl, buffer)
			program = compileProgram(gl)
			cmd     = makeAcceleratedCommand(program, array)
			window  = openWindow(gl)
		)
		sampledImage, err := gl.NewImage(2, 2)
		if err != nil {
			log.Panicf("NewImage failed: %v", err)
		}
		selection := sampledImage.WholeImageSelection()
		selection.SetColor(0, 0, image.RGB(255, 0, 0))
		selection.SetColor(1, 0, image.RGB(0, 255, 0))
		selection.SetColor(0, 1, image.RGB(255, 255, 255))
		selection.SetColor(1, 1, image.RGB(0, 0, 255))

		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			if err := screen.Modify(cmd, selection); err != nil {
				log.Panicf("Modify failed: %v", err)
			}
		})
	})
}

func makeVertexArray(gl *opengl.OpenGL, buffer *opengl.FloatVertexBuffer) *opengl.VertexArray {
	array, err := gl.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec2})
	if err != nil {
		log.Panicf("NewVertexArray failed: %v", err)
	}
	xy := opengl.VertexBufferPointer{Offset: 0, Stride: 4, Buffer: buffer}
	if err := array.Set(0, xy); err != nil {
		log.Panicf("VertexBufferPointer failed: %v", err)
	}
	st := opengl.VertexBufferPointer{Offset: 2, Stride: 4, Buffer: buffer}
	if err := array.Set(1, st); err != nil {
		log.Panicf("VertexBufferPointer failed: %v", err)
	}
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
	buffer, err := gl.NewFloatVertexBuffer(len(vertices))
	if err != nil {
		log.Panicf("NewFloatVertexBuffer failed: %v", err)
	}
	if err := buffer.Upload(0, vertices); err != nil {
		log.Panicf("Upload failed: %v", err)
	}
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

func makeAcceleratedCommand(program *opengl.Program, array *opengl.VertexArray) *opengl.AcceleratedCommand {
	cmd, err := program.AcceleratedCommand(&command{
		vertexArray: array,
	})
	if err != nil {
		log.Panicf("AcceleratedCommand failed: %v", err)
	}
	return cmd
}

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

func openWindow(gl *opengl.OpenGL) *opengl.Window {
	window, err := gl.OpenWindow(200, 200, opengl.Zoom(2))
	if err != nil {
		log.Panicf("OpenWindow failed: %v", err)
	}
	return window
}
