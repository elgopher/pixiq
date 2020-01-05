package opengl

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func compileProgram() (*program, error) {
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
	vertePositionLocation := gl.GetAttribLocation(programID, gl.Str("vertexPosition\x00"))
	texturePositionLocation := gl.GetAttribLocation(programID, gl.Str("texturePosition\x00"))
	return &program{
		id:                      programID,
		vertexPositionLocation:  vertePositionLocation,
		texturePositionLocation: texturePositionLocation,
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
