package glfw_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jacekolszak/pixiq/glfw"
)

func BenchmarkWindow_Draw(b *testing.B) {
	openGL, err := glfw.NewOpenGL(mainThreadLoop)
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

func BenchmarkWindow_PollMouseEvent(b *testing.B) {
	openGL, err := glfw.NewOpenGL(mainThreadLoop)
	require.NoError(b, err)
	defer openGL.Destroy()
	win, err := openGL.OpenWindow(640, 360)
	if err != nil {
		panic(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 3; j++ {
			// TODO Right now it is extremely inefficient
			win.PollMouseEvent()
		}
	}
}
