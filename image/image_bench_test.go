package image_test

import (
	"testing"

	"github.com/jacekolszak/pixiq/image"
)

func BenchmarkSelection_SetColor(b *testing.B) {
	b.StopTimer()
	var (
		color     = image.RGBA(10, 20, 30, 40)
		img, _    = image.New(1920, 1080, acceleratedImageStub{})
		selection = img.WholeImageSelection()
		height    = selection.Height()
		width     = selection.Width()
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.SetColor(x, y, color)
			}
		}
	}
}

func BenchmarkSelection_Color(b *testing.B) {
	b.StopTimer()
	var (
		img, _    = image.New(1920, 1080, acceleratedImageStub{})
		selection = img.WholeImageSelection()
		height    = selection.Height()
		width     = selection.Width()
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				selection.Color(x, y)
			}
		}
	}
}

// Must be 0 allocs/op
func BenchmarkSelection(b *testing.B) {
	b.StopTimer()
	var (
		img, _    = image.New(1, 1, acceleratedImageStub{})
		selection = img.WholeImageSelection()
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		selection.Selection(0, 0).WithSize(1, 1)
	}
}

// Should be 0 allocs/op
func BenchmarkSelection_Modify(b *testing.B) {
	b.StopTimer()
	var (
		img, _    = image.New(1920, 1080, acceleratedImageStub{})
		selection = img.WholeImageSelection()
		program   = struct{}{}
		primitive = struct{}{}
		param     = struct{}{}
	)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		procedure := func(drawer image.Drawer) {
			drawer.SetSelection("selection", selection)
			_ = drawer.Draw(primitive, param)
		}
		_ = selection.Modify(program, procedure)
	}
}
