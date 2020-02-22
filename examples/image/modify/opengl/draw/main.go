package main

import (
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
	"log"
	"time"
)

func main() {
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		var (
			buffer  = makeVertexBuffer(gl)
			array   = makeVertexArray(gl, buffer)
			program = compileProgram(gl)
			cmd     = program.AcceleratedCommand(&drawColorfulRectangle{vertexArray: array})
			window  = openWindow(gl)
		)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			screen.Modify(cmd)
		})
	})
}

type drawColorfulRectangle struct {
	vertexArray *opengl.VertexArray
}

func (c drawColorfulRectangle) RunGL(renderer *opengl.Renderer, _ []image.AcceleratedImageSelection) {
	var opacity = float32(time.Now().Nanosecond()) / float32(time.Second)
	renderer.SetFloat("opacity", opacity)
	renderer.DrawArrays(c.vertexArray, opengl.TriangleFan, 0, 4)
}

func makeVertexBuffer(gl *opengl.OpenGL) *opengl.FloatVertexBuffer {
	// xy -> rgb
	vertices := []float32{
		-1, 1, 1, 0, 0, // top-left -> red
		1, 1, 0, 1, 0, // top-right -> green
		1, -1, 0, 0, 1, // bottom-right -> blue
		-1, -1, 1, 1, 1, // bottom-left -> white
	}
	buffer := gl.NewFloatVertexBuffer(len(vertices))
	buffer.Upload(0, vertices)
	return buffer
}

func makeVertexArray(gl *opengl.OpenGL, buffer *opengl.FloatVertexBuffer) *opengl.VertexArray {
	array := gl.NewVertexArray(opengl.VertexLayout{opengl.Vec2, opengl.Vec3})
	xy := opengl.VertexBufferPointer{Offset: 0, Stride: 5, Buffer: buffer}
	array.Set(0, xy)
	color := opengl.VertexBufferPointer{Offset: 2, Stride: 5, Buffer: buffer}
	array.Set(1, color)
	return array
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

	uniform float opacity;

	in vec4 interpolatedColor;
	out vec4 color;
	
	void main() {
		color = interpolatedColor * opacity;
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
	program, err := gl.LinkProgram2(vertexShader, fragmentShader)
	if err != nil {
		log.Panicf("LinkProgram failed: %v", err)
	}
	return program
}

func openWindow(gl *opengl.OpenGL) *opengl.Window {
	window, err := gl.OpenWindow(200, 200, opengl.Zoom(2))
	if err != nil {
		log.Panicf("OpenWindow failed: %v", err)
	}
	return window
}
