package opengl

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func compileProgram(vertexShaderSrc, fragmentShaderSrc string) (*program, error) {
	vertexShader, err := compileVertexShader(vertexShaderSrc)
	if err != nil {
		return nil, err
	}
	defer vertexShader.delete()

	fragmentShader, err := compileFragmentShader(fragmentShaderSrc)
	if err != nil {
		return nil, err
	}
	defer fragmentShader.delete()

	program, err := linkProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	return program, nil
}

type program struct {
	id                      uint32
	vertexPositionLocation  int32
	texturePositionLocation int32
}

func (p *program) use() {
	gl.UseProgram(p.id)
}

func (p *program) uniformNames() map[string]int32 {
	names := map[string]int32{}
	var count, length, size int32
	var xtype uint32
	nameMaxSize := int32(64)
	name := make([]byte, nameMaxSize)
	gl.GetProgramiv(p.id, gl.ACTIVE_UNIFORMS, &count)
	for i := int32(0); i < count; i++ {
		gl.GetActiveUniform(p.id, uint32(i), nameMaxSize, &length, &size, &xtype, &name[0])
		names[string(name)] = i
	}
	return names
}

func linkProgram(shaders ...*shader) (*program, error) {
	programID := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(programID, shader.id)
	}
	gl.LinkProgram(programID)
	var success int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var infoLogLen int32
		gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &infoLogLen)
		infoLog := make([]byte, infoLogLen)
		if infoLogLen > 0 {
			gl.GetProgramInfoLog(programID, infoLogLen, nil, &infoLog[0])
		}
		return nil, fmt.Errorf("error linking program: %s", string(infoLog))
	}
	return &program{
		id:                      programID,
		vertexPositionLocation:  0,
		texturePositionLocation: 1,
	}, nil
}

type shader struct {
	id uint32
}

func compileVertexShader(src string) (*shader, error) {
	return compileShader(gl.VERTEX_SHADER, src)
}

func compileFragmentShader(src string) (*shader, error) {
	return compileShader(gl.FRAGMENT_SHADER, src)
}

func compileShader(xtype uint32, src string) (*shader, error) {
	if src == "" {
		src = " "
	}
	shaderID := gl.CreateShader(xtype)
	srcXString, free := gl.Strs(src)
	defer free()
	length := int32(len(src))
	gl.ShaderSource(shaderID, 1, srcXString, &length)
	gl.CompileShader(shaderID)
	var success int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLen)
		infoLog := make([]byte, logLen)
		if logLen > 0 {
			gl.GetShaderInfoLog(shaderID, logLen, nil, &infoLog[0])
		}
		return nil, fmt.Errorf("error compiling shader: %s", string(infoLog))
	}
	return &shader{id: shaderID}, nil
}

func (s *shader) delete() {
	gl.DeleteShader(s.id)
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
