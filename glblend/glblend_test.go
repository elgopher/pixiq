package glblend_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacekolszak/pixiq/glblend"
)

func TestNewSource(t *testing.T) {
	t.Run("should panic when context is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = glblend.NewSource(nil)
		})
	})
}
