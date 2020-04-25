package main

import (
	"log"
	"time"

	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/glfw"
	"github.com/jacekolszak/pixiq/image"
	"github.com/jacekolszak/pixiq/loop"
)

func main() {
	glfw.RunOrDie(func(openGL *glfw.OpenGL) {
		var (
			context = openGL.Context()
			buffer  = makeVertexBuffer(context)
			array   = makeVertexArray(context, buffer)
			program = compileProgram(context)
			cmd     = program.AcceleratedCommand(&drawColorfulRectangle{vertexArray: array})
			window  = openWindow(openGL)
		)
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			screen.Modify(cmd)

			if window.ShouldClose() {
				frame.StopLoopEventually()
			}
		})
	})
}

type drawColorfulRectangle struct {
	vertexArray *gl.VertexArray
}

func (c drawColorfulRectangle) RunGL(renderer *gl.Renderer, _ []image.AcceleratedImageSelection) {
	var opacity = float32(time.Now().Nanosecond()) / float32(time.Second)
	renderer.SetFloat("opacity", opacity)
	renderer.DrawArrays(c.vertexArray, gl.TriangleFan, 0, 4)
}

func makeVertexBuffer(context *gl.Context) *gl.FloatVertexBuffer {
	// xy -> rgb
	vertices := []float32{
		-1, 1, 1, 0, 0, // top-left -> red
		1, 1, 0, 1, 0, // top-right -> green
		1, -1, 0, 0, 1, // bottom-right -> blue
		-1, -1, 1, 1, 1, // bottom-left -> white
	}
	buffer := context.NewFloatVertexBuffer(len(vertices), gl.StaticDraw)
	buffer.Upload(0, vertices)
	return buffer
}

func makeVertexArray(context *gl.Context, buffer *gl.FloatVertexBuffer) *gl.VertexArray {
	array := context.NewVertexArray(gl.VertexLayout{gl.Vec2, gl.Vec3})
	xy := gl.VertexBufferPointer{Offset: 0, Stride: 5, Buffer: buffer}
	array.Set(0, xy)
	color := gl.VertexBufferPointer{Offset: 2, Stride: 5, Buffer: buffer}
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

func openWindow(openGL *glfw.OpenGL) *glfw.Window {
	window, err := openGL.OpenWindow(200, 200, glfw.Zoom(2))
	if err != nil {
		log.Panicf("OpenWindow failed: %v", err)
	}
	return window
}
