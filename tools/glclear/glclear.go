package glclear

import (
	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/image"
)

func New(command *gl.ClearCommand) *Tool {
	return &Tool{
		command: command,
	}
}

type Tool struct {
	command *gl.ClearCommand
}

func (t *Tool) SetColor(color image.Color) {
	t.command.SetColor(color)
}

func (t *Tool) Clear(selection image.Selection) {
	selection.Modify(t.command)
}
