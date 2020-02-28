package glclear_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/tools/glclear"
)

func TestNew(t *testing.T) {
	t.Run("should panic for nil command", func(t *testing.T) {
		assert.Panics(t, func() {
			// when
			glclear.New(nil)
		})
	})
}
