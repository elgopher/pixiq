package glfw

import (
	"github.com/jacekolszak/pixiq/gl"
)

func compileProgram(context *gl.Context, vertexShaderSrc, fragmentShaderSrc string) (*gl.Program, error) {
	vertexShader, err := context.CompileVertexShader(vertexShaderSrc)
	if err != nil {
		return nil, err
	}
	defer vertexShader.Delete()
	fragmentShader, err := context.CompileFragmentShader(fragmentShaderSrc)
	if err != nil {
		return nil, err
	}
	defer fragmentShader.Delete()
	program, err := context.LinkProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}
	return program, nil
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
