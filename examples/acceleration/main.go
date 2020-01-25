package main

import "github.com/jacekolszak/pixiq/opengl"

func main() {
	// Use OpenGL on PCs with Linux, Windows and MacOS.
	// This package can open windows and draw images on them.
	opengl.RunOrDie(func(gl *opengl.OpenGL) {
		vertexShader, err := gl.CompileVertexShader(vertexShaderSrc)
		if err != nil {
			panic(err)
		}
		fragmentShader, err := gl.CompileFragmentShader(fragmentShaderSrc)
		if err != nil {
			panic(err)
		}
		gl.LinkProgram(vertexShader, fragmentShader)
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
