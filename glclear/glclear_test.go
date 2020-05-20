package glclear_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/gl"
	"github.com/jacekolszak/pixiq/glclear"
)

func TestNew(t *testing.T) {
	t.Run("should panic for nil context", func(t *testing.T) {
		assert.Panics(t, func() {
			// when
			glclear.New(nil)
		})
	})
	t.Run("should panic for nil command", func(t *testing.T) {
		assert.Panics(t, func() {
			// when
			glclear.New(nilCommandContext{})
		})
	})
}

type nilCommandContext struct {
}

func (a nilCommandContext) NewClearCommand() *gl.ClearCommand {
	return nil
}
