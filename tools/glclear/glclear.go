package glclear

import (
	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/image"
)

// New returns new instance of *glclear.Tool
func New(command *gl.ClearCommand) *Tool {
	if command == nil {
		panic("nil command")
	}
	return &Tool{
		command: command,
	}
}

// Tool is a clearing tool. It clears the image.Selection with specific color
// which basically means setting the color for each pixel in the selection.
//
// Tool is using GPU through use of *gl.ClearCommand
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
