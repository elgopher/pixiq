package opengl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/opengl"
)

func BenchmarkWindow_Draw(b *testing.B) {
	openGL, err := opengl.New(mainThreadLoop)
	require.NoError(b, err)
	defer openGL.Destroy()
	win, err := openGL.OpenWindow(640, 360)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		win.Draw()
	}
}
