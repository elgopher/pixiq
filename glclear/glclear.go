// Package glclear provides GPU tools for clearing selections
// using specific color
package glclear

import (
	"github.com/elgopher/pixiq/gl"
	"github.com/elgopher/pixiq/image"
)

// glContext is an OpenGL context. Possible implementation is *gl.Context which can be
// obtained by calling OpenGL.Context()
type glContext interface {
	NewClearCommand() *gl.ClearCommand
}

// New returns a new instance of *glclear.Tool.
func New(context glContext) *Tool {
	if context == nil {
		panic("nil context")
	}
	command := context.NewClearCommand()
	if command == nil {
		panic("nil *gl.ClearCommand returned by context")
	}
	return &Tool{
		command: command,
	}
}

// Tool is a clearing tool. It clears the image.Selection with specific color
// which basically means setting the color for each pixel in the selection.
//
// Tool uses GPU through use of *gl.ClearCommand
type Tool struct {
	command *gl.ClearCommand
}

// SetColor sets color which will be used by Clear method
func (t *Tool) SetColor(color image.Color) {
	t.command.SetColor(color)
}

// Clear clears selection with previously set color
func (t *Tool) Clear(selection image.Selection) {
	selection.Modify(t.command)
}
