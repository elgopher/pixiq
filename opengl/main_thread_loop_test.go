package opengl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainThreadLoop_Execute(t *testing.T) {
	t.Run("should not do anything for nil job", func(t *testing.T) {
		mainThreadLoop.Execute(nil)
	})
	t.Run("should execute job synchronously", func(t *testing.T) {
		var executed bool
		// when
		mainThreadLoop.Execute(func() {
			executed = true
		})
		assert.True(t, executed)
	})
	t.Run("should execute two jobs", func(t *testing.T) {
		var executionCount int
		job := func() {
			executionCount += 1
		}
		// when
		mainThreadLoop.Execute(job)
		mainThreadLoop.Execute(job)
		// then
		assert.Equal(t, 2, executionCount)
	})
}
