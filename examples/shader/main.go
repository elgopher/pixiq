package main

import (
	"fmt"

	"github.com/jacekolszak/pixiq/colornames"
	program2 "github.com/jacekolszak/pixiq/glsl/program"
	"github.com/jacekolszak/pixiq/loop"
	"github.com/jacekolszak/pixiq/opengl"
)

func main() {
	// Use OpenGL on PCs with Linux, Windows and MacOS.
	// This package can open windows and draw images on them.
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		window, err := gl.OpenWindow(80, 16, opengl.Zoom(5))
		if err != nil {
			panic(err)
		}

		program := gl.DrawProgram()
		program.SetVertexShader(vertexShaderSrc)
		program.SetFragmentShader(fragmentShaderSrc)
		compiledProgram, err := program.Compile()
		if err != nil {
			panic(err)
		}
		vertexPosition := compiledProgram.GetVertexAttributeLocation("vertexPosition")
		texturePosition := compiledProgram.GetVertexAttributeLocation("texturePosition")
		fmt.Println(vertexPosition)
		fmt.Println(texturePosition)

		vbo := gl.NewFloatVertexBuffer(program2.StaticDraw)
		data := []float32{
			-1, -1, 0, 1, // (x,y) -> (u,v), that is: vertexPosition -> texturePosition
			1, -1, 1, 1,
			1, 1, 1, 0,
			-1, -1, 0, 1,
			1, 1, 1, 0,
			-1, 1, 0, 0,
		}
		vbo.Update(0, data)

		// TODO Shared VAO is not possible
		vao := compiledProgram.NewVertexArrayObject()
		vao.SetVertexAttribute(vertexPosition, vbo.Pointer(0, 2, 4))
		vao.SetVertexAttribute(texturePosition, vbo.Pointer(2, 2, 4))
		fmt.Println(vao)

		// TODO FINISH THIS ONE
		compiledProgram.NewCall(func(call program2.DrawCall) {
			call.BindVertexArrayObject(vao)
			call.BindTexture0(nil)
			call.Draw(program2.Triangles, 0, 6)
		})

		// Create a loop for a screen. OpenGL's Window is a Screen (some day
		// in the future Pixiq may support different platforms such as mobile
		// or browser, therefore we need a Screen abstraction).
		// Each iteration of the loop is a Frame.
		loop.Run(window, func(frame *loop.Frame) {
			screen := frame.Screen()
			screen.SetColor(40, 8, colornames.White)
		})
	})
}

const vertexShaderSrc = `
#version 330 core

in vec2 vertexPosition;
in vec2 texturePosition;

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
