package shader

import "github.com/jacekolszak/pixiq/image"

type Compiler interface {
	Compile(fragmentShaderSource string) (Program, error)
}

type Program interface {
	Call() ProgramCall
}

type ProgramCall interface {
	SetTexture(uniformName string, selection image.Selection)
	// Release should be called after program is executed. The Program implementation
	// might reuse the released ProgramCall instance next time a new ProgramCall is
	// created by Program.Call
	Release()
}
